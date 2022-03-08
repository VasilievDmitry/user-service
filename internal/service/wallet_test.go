package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/internal/repository/mocks"
	microErrors "github.com/micro/go-micro/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type WalletTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *Service
}

func Test_WalletTestSuite(t *testing.T) {
	suite.Run(t, new(WalletTestSuite))
}

func (suite *WalletTestSuite) SetupSuite() {
	var err error

	suite.cfg, err = config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", err)
	}

	log, _ := zap.NewProduction()
	suite.service = NewService(repository.InitRepositories(nil, nil), suite.cfg, log)
}

func (suite *WalletTestSuite) SetupTest() {
}

func (suite *WalletTestSuite) TearDownTest() {
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_UnsupportedWalletType() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: "unknown",
		}
		res = &user_service.ResponseWithUserProfile{}
	)

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorWalletUnsupportedType, mErr.Detail)
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_GetByTokenDbError() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: user_service.WalletTypePhantom,
			Token:    "token",
		}
		res = &user_service.ResponseWithUserProfile{}
	)

	authProviderRep := &mocks.AuthProviderRepositoryInterface{}
	authProviderRep.On("GetByToken", ctx, req.Provider, req.Token).Return(nil, errors.New("db_error"))
	suite.service.repositories.AuthProvider = authProviderRep

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_InsertUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: user_service.WalletTypePhantom,
			Token:    "token",
		}
		res  = &user_service.ResponseWithUserProfile{}
		user = &user_service.User{
			IsActive: true,
		}
	)

	authProviderRep := &mocks.AuthProviderRepositoryInterface{}
	authProviderRep.On("GetByToken", ctx, req.Provider, req.Token).Return(nil, sql.ErrNoRows)
	suite.service.repositories.AuthProvider = authProviderRep

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("Insert", ctx, user).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_InsertAuthProviderDbError() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: user_service.WalletTypePhantom,
			Token:    "token",
		}
		res  = &user_service.ResponseWithUserProfile{}
		user = &user_service.User{
			Id:       "user_id",
			IsActive: true,
		}
	)

	authProviderRep := &mocks.AuthProviderRepositoryInterface{}
	authProviderRep.On("GetByToken", ctx, req.Provider, req.Token).Return(nil, sql.ErrNoRows)
	authProviderRep.On("Insert", ctx, mock.MatchedBy(func(input *user_service.AuthProvider) bool {
		return input.User != nil && input.User.Id == user.Id && input.Provider == req.Provider && input.Token == req.Token
	})).Return(errors.New("db_error"))
	suite.service.repositories.AuthProvider = authProviderRep

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("Insert", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		if input.Id == "" && input.IsActive == true {
			input.Id = user.Id
			return true
		}
		return false
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_NewUser() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: user_service.WalletTypePhantom,
			Token:    "token",
		}
		res  = &user_service.ResponseWithUserProfile{}
		user = &user_service.User{
			Id:       "user_id",
			IsActive: true,
		}
		profile = suite.service.convertUserToProfile(user)
	)

	authProviderRep := &mocks.AuthProviderRepositoryInterface{}
	authProviderRep.On("GetByToken", ctx, req.Provider, req.Token).Return(nil, sql.ErrNoRows)
	authProviderRep.On("Insert", ctx, mock.MatchedBy(func(input *user_service.AuthProvider) bool {
		return input.User != nil && input.User.Id == user.Id && input.Provider == req.Provider && input.Token == req.Token
	})).Return(nil)
	suite.service.repositories.AuthProvider = authProviderRep

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("Insert", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		if input.Id == "" && input.IsActive == true {
			input.Id = user.Id
			return true
		}
		return false
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res.UserProfile, profile)
}

func (suite *WalletTestSuite) Test_CreateUserByWallet_ExistsUser() {
	var (
		ctx = context.Background()
		req = &user_service.CreateUserByWalletRequest{
			Provider: user_service.WalletTypePhantom,
			Token:    "token",
		}
		res  = &user_service.ResponseWithUserProfile{}
		user = &user_service.User{
			Id:       "user_id",
			IsActive: true,
		}
		authProvider = &user_service.AuthProvider{
			User:     user,
			Provider: req.Provider,
			Token:    req.Token,
		}
		profile = suite.service.convertUserToProfile(authProvider.User)
	)

	authProviderRep := &mocks.AuthProviderRepositoryInterface{}
	authProviderRep.On("GetByToken", ctx, req.Provider, req.Token).Return(authProvider, nil)
	suite.service.repositories.AuthProvider = authProviderRep

	err := suite.service.CreateUserByWallet(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), res.UserProfile, profile)
}