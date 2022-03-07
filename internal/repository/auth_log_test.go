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
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
)

type AuthLogTestSuite struct {
	suite.Suite
	db         *sqlx.DB
	userRep    UserRepositoryInterface
	authLogRep AuthLogRepositoryInterface
	user       *user_service.User
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

	suite.db, err = sqlx.Open("mysql", suite.cfg.MysqlDsn)
	if err != nil {
		suite.FailNow("Database connection failed", err)
	}

	err = db.Migrate("file://../../migrations", "mysql://"+suite.cfg.MysqlDsn, true, suite.cfg.MigrationsLockTimeout)
	if err != nil {
		suite.FailNow("Database migration failed", err)
	}

	suite.userRep = NewUserRepository(suite.db, log)
	suite.authLogRep = NewAuthLogRepository(suite.db, log)
}

func (suite *AuthLogTestSuite) SetupTest() {
	suite.user = &user_service.User{
		Id: uuid.NewString(),
	}

	if err := suite.userRep.Insert(context.TODO(), suite.user); err != nil {
		assert.FailNow(suite.T(), "unable to create user for test", suite.user)
	}
}

func (suite *AuthLogTestSuite) TearDownTest() {
	_, err := suite.db.Exec("TRUNCATE TABLE auth_log")
	if err != nil {
		suite.FailNow("Unable to truncate table", err)
	}

	_, err = suite.db.Exec("TRUNCATE TABLE user")
	if err != nil {
		suite.FailNow("Unable to truncate table", err)
	}
}

func (suite *AuthLogTestSuite) Test_NewAuthLogRepository() {
	assert.Implements(suite.T(), (*AuthLogRepositoryInterface)(nil), NewAuthLogRepository(nil, nil))
}

func (suite *AuthLogTestSuite) Test_CRUD() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	log2, err := suite.authLogRep.GetByAccessToken(ctx, log.AccessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), log.Id, log2.Id)
	assert.NotEmpty(suite.T(), log2.CreatedAt)
	assert.NotEmpty(suite.T(), log2.UpdatedAt)

	log2.Ip = "ip2"
	err = suite.authLogRep.Update(ctx, log2)
	assert.NoError(suite.T(), err)

	log3, err := suite.authLogRep.GetByAccessToken(ctx, log.AccessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), log.Id, log3.Id)
	assert.Equal(suite.T(), log2.Ip, log3.Ip)
	assert.Equal(suite.T(), log2.CreatedAt, log3.CreatedAt)
	assert.GreaterOrEqual(suite.T(), log3.UpdatedAt.Seconds, log2.UpdatedAt.Seconds)
}

func (suite *AuthLogTestSuite) Test_Insert_MappingError() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	log.ExpireAt = nil

	err := suite.authLogRep.Insert(ctx, log)
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) Test_Update_MappingError() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	log.ExpireAt = nil

	err := suite.authLogRep.Update(ctx, log)
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) Test_GetByActiveAccessToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	log2, err := suite.authLogRep.GetByAccessToken(ctx, log.AccessToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), log.Id, log2.Id)
	assert.Equal(suite.T(), log.AccessToken, log2.AccessToken)
	assert.Equal(suite.T(), log.RefreshToken, log2.RefreshToken)
	assert.Equal(suite.T(), log.UserAgent, log2.UserAgent)
	assert.Equal(suite.T(), log.Ip, log2.Ip)
	assert.Equal(suite.T(), log.IsActive, log2.IsActive)
	assert.Equal(suite.T(), log.ExpireAt.Seconds, log2.ExpireAt.Seconds)
	assert.NotEmpty(suite.T(), log2.CreatedAt.Seconds)
	assert.NotEmpty(suite.T(), log2.UpdatedAt.Seconds)
	assert.NotEmpty(suite.T(), log2.User)
	assert.Equal(suite.T(), log.User.Id, log2.User.Id)
}

func (suite *AuthLogTestSuite) Test_GetByDisabledAccessToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	log.IsActive = false

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	_, err = suite.authLogRep.GetByAccessToken(ctx, log.AccessToken)
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) Test_GetByUnknownAccessToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	_, err = suite.authLogRep.GetByAccessToken(ctx, "unknown")
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) Test_GetByActiveRefreshToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	log2, err := suite.authLogRep.GetByRefreshToken(ctx, log.RefreshToken)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), log.Id, log2.Id)
	assert.Equal(suite.T(), log.AccessToken, log2.AccessToken)
	assert.Equal(suite.T(), log.RefreshToken, log2.RefreshToken)
	assert.Equal(suite.T(), log.UserAgent, log2.UserAgent)
	assert.Equal(suite.T(), log.Ip, log2.Ip)
	assert.Equal(suite.T(), log.IsActive, log2.IsActive)
	assert.Equal(suite.T(), log.ExpireAt.Seconds, log2.ExpireAt.Seconds)
	assert.NotEmpty(suite.T(), log2.CreatedAt.Seconds)
	assert.NotEmpty(suite.T(), log2.UpdatedAt.Seconds)
	assert.NotEmpty(suite.T(), log2.User)
	assert.Equal(suite.T(), log.User.Id, log2.User.Id)
}

func (suite *AuthLogTestSuite) Test_GetByDisabledRefreshToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	log.IsActive = false

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	_, err = suite.authLogRep.GetByRefreshToken(ctx, log.RefreshToken)
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) Test_GetByUnknownRefreshToken() {
	var (
		ctx = context.Background()
		log = suite.getDefaultLog()
	)

	err := suite.authLogRep.Insert(ctx, log)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), log.Id)

	_, err = suite.authLogRep.GetByRefreshToken(ctx, "unknown")
	assert.Error(suite.T(), err)
}

func (suite *AuthLogTestSuite) getDefaultLog() *user_service.AuthLog {
	return &user_service.AuthLog{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		UserAgent:    "user_agent",
		Ip:           "ip",
		IsActive:     true,
		ExpireAt:     timestamppb.Now(),
		User:         suite.user,
	}
}
