package internal

import (
	"errors"
	"fmt"
	"io"
	"os"

	clientGrpc "github.com/go-micro/plugins/v4/client/grpc"
	"github.com/go-micro/plugins/v4/registry/etcd"
	serverGrpc "github.com/go-micro/plugins/v4/server/grpc"
	transportGrpc "github.com/go-micro/plugins/v4/transport/grpc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/log"
	"github.com/natefinch/lumberjack"
	"go-micro.dev/v4"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	gameService "github.com/lotproject/user-service/proto/game-service"

	"github.com/lotproject/user-service/pkg"
	userService "github.com/lotproject/user-service/proto/v1"

	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/internal/service"
)

// Application is application entry point.
type Application struct {
	cfg               *config.Config
	log               *zap.Logger
	database          *sqlx.DB
	service           *service.Service
	micro             micro.Service
	gameServiceClient gameService.GameService
}

// NewApplication create new Application.
func NewApplication() (app *Application) {
	app = &Application{}

	app.initConfig()
	app.initLogger()
	app.initDatabase()
	app.initMicroServices()

	app.service = service.NewService(
		repository.InitRepositories(app.database, app.log),
		app.gameServiceClient,
		app.cfg,
		app.log,
	)

	return
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

	app.log = logger.Named(pkg.ServiceName)
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
	r := etcd.NewRegistry(registry.Addrs(app.cfg.MicroRegistryAddress))
	t := transportGrpc.NewTransport()
	c := clientGrpc.NewClient()
	s := serverGrpc.NewServer(
		server.Name(pkg.ServiceName),
		server.Registry(r),
		server.Transport(t),
	)

	options := []micro.Option{
		micro.Registry(r),
		micro.Transport(t),
		micro.Server(s),
		micro.Client(c),
		micro.AfterStop(func() error {
			app.log.Info("Micro service stopped")
			app.Stop()
			return nil
		}),
	}

	app.micro = micro.NewService(options...)
	app.micro.Init()

	app.gameServiceClient = gameService.NewGameService("lot.game.v1", app.micro.Client())

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
	if err := userService.RegisterUserServiceHandler(app.micro.Server(), app.service); err != nil {
		app.log.Fatal("Micro service starting failed", zap.Error(err))
	}

	if err := app.micro.Run(); err != nil {
		panic("Can`t run service")
	}
}

func (app *Application) Stop() {
	app.log.Info("Shutdown the game service application")

	err := app.database.Close()
	if err != nil {
		app.log.Error("DB connection close failed", zap.Error(err))
	} else {
		app.log.Info("DB connection close success")
	}

	if err := app.log.Sync(); err != nil {
		app.log.Error("Logger sync failed", zap.Error(err))
	} else {
		app.log.Info("Logger synced")
	}
}
