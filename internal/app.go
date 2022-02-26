package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/InVisionApp/go-health/v2"
	"github.com/InVisionApp/go-health/v2/handlers"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/log"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/internal/service"
	"github.com/micro/go-micro"
	"github.com/natefinch/lumberjack"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
	"time"
)

// Application is application entry point.
type Application struct {
	cfg          *config.Config
	log          *zap.Logger
	database     *sqlx.DB
	service      *service.Service
	micro        micro.Service
	healthServer *http.Server
	healthRouter *http.ServeMux
}

// NewApplication create new Application.
func NewApplication() (app *Application) {
	app = &Application{
		healthRouter: http.NewServeMux(),
	}

	app.initConfig()
	app.initLogger()
	app.initDatabase()
	app.initMicroServices()
	app.initHealth()
	app.initMetrics()

	app.service = service.NewService(
		repository.InitRepositories(app.database, app.log),
		app.cfg,
		app.log,
	)

	return
}

func (app *Application) initHealth() {
	h := health.New()
	err := h.AddChecks([]*health.Config{
		{
			Name:     "health-check",
			Checker:  app,
			Interval: time.Duration(1) * time.Second,
			Fatal:    true,
		},
	})

	if err != nil {
		app.log.Fatal("Health check register failed", zap.Error(err))
	}

	if err = h.Start(); err != nil {
		app.log.Fatal("Health check start failed", zap.Error(err))
	}

	app.log.Info("Health check listener started", zap.Int("port", app.cfg.MetricsPort))

	app.healthRouter.HandleFunc("/health", handlers.NewJSONHandlerFunc(h, nil))
}

func (app *Application) initMetrics() {
	app.healthRouter.Handle("/metrics", promhttp.Handler())
}

func (app *Application) Status() (interface{}, error) {
	err := app.database.Ping()

	if err != nil {
		return "fail", errors.New("db connection lost: " + err.Error())
	}

	return "ok", nil
}

func (app *Application) initLogger() {
	cfg := zap.NewProductionEncoderConfig()

	if app.cfg.DevelopMode {
		cfg = zap.NewDevelopmentEncoderConfig()
	}

	var writer io.Writer

	if app.cfg.LogToFileEnabled {
		writer = &lumberjack.Logger{
			Filename:   app.cfg.LogFilePath,
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   false,
		}
	} else {
		writer = os.Stderr
	}

	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(writer),
		log.GetZapLevel(app.cfg.LogLevel),
	)
	logger := zap.New(core, zap.AddCaller())

	app.log = logger.Named(user_service.ServiceName)
	zap.ReplaceGlobals(app.log)

	app.log.Info("Logger init...")
}

func (app *Application) initConfig() {
	var err error

	app.cfg, err = config.NewConfig()

	if err != nil {
		zap.S().Panic("Config init failed", zap.Error(err))
	}

	zap.S().Info("Configuration parsed successfully...")
}

func (app *Application) initDatabase() {
	var err error

	app.database, err = sqlx.Open("mysql", app.cfg.MysqlDsn)
	if err != nil {
		app.log.Fatal("Database connection failed", zap.Error(err))
	}

	app.log.Info("Database initialization successfully...")
}

func (app *Application) initMicroServices() {
	options := []micro.Option{
		micro.Name(user_service.ServiceName),
		micro.AfterStop(func() error {
			app.log.Info("Micro service stopped")
			app.Stop()
			return nil
		}),
	}

	app.micro = micro.NewService(options...)
	app.micro.Init()

	app.log.Info("Micro service initialization successfully...")
}

// Run starts application
func (app *Application) Run() {
	app.log.Info("Starting the login application")

	dsn := fmt.Sprintf("mysql://%s", app.cfg.MysqlDsn)
	err := db.Migrate("file://./migrations", dsn, true, app.cfg.MigrationsLockTimeout)

	if err != nil {
		app.log.Fatal("DB migrations failed", zap.Error(err))
	}

	app.healthServer = &http.Server{
		Addr:              fmt.Sprintf(":%d", app.cfg.MetricsPort),
		Handler:           app.healthRouter,
		ReadTimeout:       time.Duration(app.cfg.MetricsReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(app.cfg.MetricsReadHeaderTimeout) * time.Second,
	}

	if err := user_service.RegisterUserServiceHandler(app.micro.Server(), app.service); err != nil {
		app.log.Fatal("Micro service starting failed", zap.Error(err))
	}

	if err := app.micro.Run(); err != nil {
		panic("Can`t run service")
	}
}

func (app *Application) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app.log.Info("Shutdown the game service application")

	err := app.database.Close()
	if err != nil {
		app.log.Error("DB connection close failed", zap.Error(err))
	} else {
		app.log.Info("DB connection close success")
	}

	if app.healthServer != nil {
		if err := app.healthServer.Shutdown(ctx); err != nil {
			app.log.Error("Health server shutdown failed", zap.Error(err))
		}
		app.log.Info("Health server stopped")
	}

	if err := app.log.Sync(); err != nil {
		app.log.Error("Logger sync failed", zap.Error(err))
	} else {
		app.log.Info("Logger synced")
	}
}
