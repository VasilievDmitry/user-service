package config

import (
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap/zapcore"
)

const (
	//Debug has verbose message
	Debug = "debug"
	//Info is default log level
	Info = "info"
	//Warn is for logging messages about possible issues
	Warn = "warn"
	//Error is for logging errors
	Error = "error"
	//Fatal is for logging fatal messages. The system shutdown after logging the message.
	Fatal = "fatal"
)

// Config define application config object
type Config struct {
	DevelopMode bool `envconfig:"DEVELOP_MODE" required:"false" default:"false"`

	MetricsPort              int `envconfig:"METRICS_PORT" required:"false" default:"8086"`
	MetricsReadTimeout       int `envconfig:"METRICS_READ_TIMEOUT" default:"60"`
	MetricsReadHeaderTimeout int `envconfig:"METRICS_READ_HEADER_TIMEOUT" default:"60"`

	LogFilePath      string `envconfig:"LOG_FILE_PATH" required:"false" default:"./logs/log.txt"`
	LogLevel         string `envconfig:"LOG_LEVEL" required:"false" default:"error"`
	LogToFileEnabled bool   `envconfig:"LOG_TO_FILE_ENABLED" required:"false" default:"false"`

	MysqlDsn              string `envconfig:"MYSQL_DSN" required:"true"`
	MigrationsLockTimeout int64  `envconfig:"MIGRATIONS_LOCK_TIMEOUT" default:"120"`

	BcryptCost               int    `envconfig:"BCRYPT_COST" required:"false" default:"10"`
	RefreshTokenLifetime     int    `envconfig:"REFRESH_TOKEN_LIFETIME" required:"false" default:"365"`
	AccessTokenLifetime      int    `envconfig:"ACCESS_TOKEN_LIFETIME" required:"false" default:"3"`
	AccessTokenSecret        string `envconfig:"ACCESS_TOKEN_SECRET" required:"true"`
	AccessTokenSigningMethod string `envconfig:"ACCESS_TOKEN_SIGNING_METHOD" required:"false" default:"HS256"`
}

// NewConfig returns actual config instance
func NewConfig() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process("", cfg)

	return cfg, err
}

func GetZapLevel(level string) zapcore.Level {
	switch level {
	case Info:
		return zapcore.InfoLevel
	case Warn:
		return zapcore.WarnLevel
	case Debug:
		return zapcore.DebugLevel
	case Error:
		return zapcore.ErrorLevel
	case Fatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
