package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type AuthLogTestSuite struct {
	suite.Suite
	repository AuthLogRepositoryInterface
	cfg        *config.Config
}

func Test_AuthLogTestSuite(t *testing.T) {
	suite.Run(t, new(AuthLogTestSuite))
}

func (suite *AuthLogTestSuite) SetupSuite() {
	var err error

	suite.cfg, err = config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", err)
	}

	log, _ := zap.NewProduction()

	err = db.CreateScheme(suite.cfg.MysqlDsn)
	if err != nil {
		suite.FailNow("Unable to create database scheme", err)
	}

	driver, err := sqlx.Open("mysql", suite.cfg.MysqlDsn)
	if err != nil {
		suite.FailNow("Database connection failed", err)
	}

	err = db.Migrate("file://../../migrations", "mysql://"+suite.cfg.MysqlDsn, true, suite.cfg.MigrationsLockTimeout)
	if err != nil {
		suite.FailNow("Database migration failed", err)
	}

	suite.repository = NewAuthLogRepository(driver, log)
}

func (suite *AuthLogTestSuite) TestOneTimeToken_CRUD() {
	var (
		ctx = context.Background()
		log = &user_service.AuthLog{
			AccessToken:  "access_token",
			RefreshToken: "refresh_token",
			UserAgent:    "user_agent",
			Ip:           "ip",
			User: &user_service.User{
				Id: uuid.NewString(),
			},
		}
	)

	err := suite.repository.Insert(ctx, log)
	assert.Error(suite.T(), err)
}
