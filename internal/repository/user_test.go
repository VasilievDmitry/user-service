package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	db      *sqlx.DB
	userRep UserRepositoryInterface
	cfg     *config.Config
}

func Test_UserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (suite *UserTestSuite) SetupSuite() {
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

	suite.db, err = sqlx.Open("mysql", suite.cfg.MysqlDsn)
	if err != nil {
		suite.FailNow("Database connection failed", err)
	}

	err = db.Migrate("file://../../migrations", "mysql://"+suite.cfg.MysqlDsn, true, suite.cfg.MigrationsLockTimeout)
	if err != nil {
		suite.FailNow("Database migration failed", err)
	}

	suite.userRep = NewUserRepository(suite.db, log)
}

func (suite *UserTestSuite) SetupTest() {
}

func (suite *UserTestSuite) TearDownTest() {
	_, err := suite.db.Exec("TRUNCATE TABLE user")
	if err != nil {
		suite.FailNow("Unable to truncate table", err)
	}
}

func (suite *UserTestSuite) Test_CRUD() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Login:          "login",
			Password:       "password",
			Username:       "username",
			IsActive:       true,
			EmailConfirmed: true,
			EmailCode:      "email_code",
			RecoveryCode:   "recovery_code",
		}
	)

	err := suite.userRep.Insert(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.Id)

	user2, err := suite.userRep.GetById(ctx, user.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Id, user2.Id)
	assert.Equal(suite.T(), user.Login, user2.Login)
	assert.Equal(suite.T(), user.Password, user2.Password)
	assert.Equal(suite.T(), user.Username, user2.Username)
	assert.Equal(suite.T(), user.IsActive, user2.IsActive)
	assert.Equal(suite.T(), user.EmailConfirmed, user2.EmailConfirmed)
	assert.Equal(suite.T(), user.EmailCode, user2.EmailCode)
	assert.Equal(suite.T(), user.RecoveryCode, user2.RecoveryCode)
	assert.NotEmpty(suite.T(), user2.CreatedAt)
	assert.NotEmpty(suite.T(), user2.UpdatedAt)

	user2.IsActive = false
	err = suite.userRep.Update(ctx, user2)
	assert.NoError(suite.T(), err)

	user3, err := suite.userRep.GetById(ctx, user.Id)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Id, user3.Id)
	assert.Equal(suite.T(), user.Login, user2.Login)
	assert.Equal(suite.T(), user.Password, user2.Password)
	assert.Equal(suite.T(), user.Username, user2.Username)
	assert.Equal(suite.T(), user2.IsActive, user2.IsActive)
	assert.Equal(suite.T(), user.EmailConfirmed, user2.EmailConfirmed)
	assert.Equal(suite.T(), user.EmailCode, user2.EmailCode)
	assert.Equal(suite.T(), user.RecoveryCode, user2.RecoveryCode)
	assert.Equal(suite.T(), user2.CreatedAt, user3.CreatedAt)
	assert.GreaterOrEqual(suite.T(), user3.UpdatedAt.Seconds, user2.UpdatedAt.Seconds)
}

func (suite *UserTestSuite) Test_GetByActiveLogin() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Login:          "login",
			Password:       "password",
			Username:       "username",
			IsActive:       true,
			EmailConfirmed: true,
			EmailCode:      "email_code",
			RecoveryCode:   "recovery_code",
		}
	)

	err := suite.userRep.Insert(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.Id)

	user2, err := suite.userRep.GetByLogin(ctx, user.Login)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.Id, user2.Id)
	assert.Equal(suite.T(), user.Login, user2.Login)
	assert.Equal(suite.T(), user.Password, user2.Password)
	assert.Equal(suite.T(), user.Username, user2.Username)
	assert.Equal(suite.T(), user.IsActive, user2.IsActive)
	assert.Equal(suite.T(), user.EmailConfirmed, user2.EmailConfirmed)
	assert.Equal(suite.T(), user.EmailCode, user2.EmailCode)
	assert.Equal(suite.T(), user.RecoveryCode, user2.RecoveryCode)
	assert.NotEmpty(suite.T(), user2.CreatedAt)
	assert.NotEmpty(suite.T(), user2.UpdatedAt)
}

func (suite *UserTestSuite) Test_GetByDisabledLogin() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Login:          "login",
			Password:       "password",
			Username:       "username",
			IsActive:       false,
			EmailConfirmed: true,
			EmailCode:      "email_code",
			RecoveryCode:   "recovery_code",
		}
	)

	err := suite.userRep.Insert(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.Id)

	_, err = suite.userRep.GetByLogin(ctx, user.Login)
	assert.Error(suite.T(), err)
}

func (suite *UserTestSuite) Test_GetByUnknownLogin() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Login:          "login",
			Password:       "password",
			Username:       "username",
			IsActive:       true,
			EmailConfirmed: true,
			EmailCode:      "email_code",
			RecoveryCode:   "recovery_code",
		}
	)

	err := suite.userRep.Insert(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.Id)

	_, err = suite.userRep.GetByLogin(ctx, "unknown")
	assert.Error(suite.T(), err)
}

func (suite *UserTestSuite) Test_GetByUnknownId() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Login:          "login",
			Password:       "password",
			Username:       "username",
			IsActive:       true,
			EmailConfirmed: true,
			EmailCode:      "email_code",
			RecoveryCode:   "recovery_code",
		}
	)

	err := suite.userRep.Insert(ctx, user)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), user.Id)

	_, err = suite.userRep.GetById(ctx, "unknown")
	assert.Error(suite.T(), err)
}
