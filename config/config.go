package config

import (
	"github.com/kelseyhightower/envconfig"
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
	RefreshTokenLifetime     int    `envconfig:"REFRESH_TOKEN_LIFETIME" required:"false" default:"30"`
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
