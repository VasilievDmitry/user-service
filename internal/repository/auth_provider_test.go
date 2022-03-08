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

type AuthProviderTestSuite struct {
	suite.Suite
	db              *sqlx.DB
	userRep         UserRepositoryInterface
	authProviderRep AuthProviderRepositoryInterface
	user            *user_service.User
	cfg             *config.Config
}

func Test_AuthProviderTestSuite(t *testing.T) {
	suite.Run(t, new(AuthProviderTestSuite))
}

func (suite *AuthProviderTestSuite) SetupSuite() {
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
	suite.authProviderRep = NewAuthProviderRepository(suite.db, log)
}

func (suite *AuthProviderTestSuite) SetupTest() {
	suite.user = &user_service.User{
		Id: uuid.NewString(),
	}

	if err := suite.userRep.Insert(context.Background(), suite.user); err != nil {
		assert.FailNow(suite.T(), "unable to create user for test", suite.user)
	}
}

func (suite *AuthProviderTestSuite) TearDownTest() {
	_, err := suite.db.Exec("TRUNCATE TABLE auth_provider")
	if err != nil {
		suite.FailNow("Unable to truncate table", err)
	}

	_, err = suite.db.Exec("TRUNCATE TABLE user")
	if err != nil {
		suite.FailNow("Unable to truncate table", err)
	}
}

func (suite *AuthLogTestSuite) Test_NewAuthProviderRepository() {
	assert.Implements(suite.T(), (*AuthProviderRepositoryInterface)(nil), NewAuthProviderRepository(nil, nil))
}

func (suite *AuthProviderTestSuite) Test_CRUD() {
	var (
		ctx      = context.Background()
		provider = &user_service.AuthProvider{
			Token:    "token",
			Provider: "provider",
			User:     suite.user,
		}
	)

	err := suite.authProviderRep.Insert(ctx, provider)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), provider.Id)

	provider2, err := suite.authProviderRep.GetByToken(ctx, provider.Provider, provider.Token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), provider.Id, provider2.Id)
	assert.Equal(suite.T(), provider.Provider, provider2.Provider)
	assert.Equal(suite.T(), provider.Token, provider2.Token)
	assert.NotEmpty(suite.T(), provider2.CreatedAt)
	assert.NotEmpty(suite.T(), provider2.UpdatedAt)

	provider2.Token = "token2"
	err = suite.authProviderRep.Update(ctx, provider2)
	assert.NoError(suite.T(), err)

	provider3, err := suite.authProviderRep.GetByToken(ctx, provider.Provider, provider2.Token)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), provider.Id, provider3.Id)
	assert.Equal(suite.T(), provider.Provider, provider3.Provider)
	assert.Equal(suite.T(), provider2.Token, provider3.Token)
	assert.Equal(suite.T(), provider2.CreatedAt, provider3.CreatedAt)
	assert.GreaterOrEqual(suite.T(), provider3.UpdatedAt.Seconds, provider2.UpdatedAt.Seconds)
}

func (suite *AuthProviderTestSuite) Test_GetByToken_UnknownProvider() {
	var (
		ctx      = context.Background()
		provider = &user_service.AuthProvider{
			Token:    "token",
			Provider: "provider",
			User:     suite.user,
		}
	)

	err := suite.authProviderRep.Insert(ctx, provider)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), provider.Id)

	_, err = suite.authProviderRep.GetByToken(ctx, "unknown", provider.Token)
	assert.Error(suite.T(), err)
}

func (suite *AuthProviderTestSuite) Test_GetByToken_UnknownToken() {
	var (
		ctx      = context.Background()
		provider = &user_service.AuthProvider{
			Token:    "token",
			Provider: "provider",
			User:     suite.user,
		}
	)

	err := suite.authProviderRep.Insert(ctx, provider)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), provider.Id)

	_, err = suite.authProviderRep.GetByToken(ctx, provider.Provider, "unknown")
	assert.Error(suite.T(), err)
}
