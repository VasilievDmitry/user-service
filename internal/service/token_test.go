package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/internal/repository/mocks"
	"github.com/lotproject/user-service/pkg"
	microErrors "github.com/micro/go-micro/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type TokenTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *Service
}

func Test_TokenTestSuite(t *testing.T) {
	suite.Run(t, new(TokenTestSuite))
}

func (suite *TokenTestSuite) SetupSuite() {
	var err error

	suite.cfg, err = config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", err)
	}

	log, _ := zap.NewProduction()
	suite.service = NewService(repository.InitRepositories(nil, nil), suite.cfg, log)
}

func (suite *TokenTestSuite) SetupTest() {
	suite.service.cfg.BcryptCost = bcrypt.MinCost
}

func (suite *TokenTestSuite) TearDownTest() {
}

func (suite *TokenTestSuite) Test_CreateAuthToken_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &pkg.CreateAuthTokenRequest{
			UserId: "user_id",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.CreateAuthToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorUserNotFound, mErr.Detail)
}

func (suite *TokenTestSuite) Test_CreateAuthToken_InsertAuthLogDbError() {
	var (
		ctx  = context.Background()
		user = &pkg.User{
			Id: "user_id",
		}
		req = &pkg.CreateAuthTokenRequest{
			UserId:    user.Id,
			UserAgent: "user_agent",
			Ip:        "ip",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("Insert", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == user.Id && input.IsActive == true && input.AccessToken != "" &&
			input.RefreshToken != "" && input.ExpireAt != nil && input.Ip == req.Ip && input.UserAgent == req.UserAgent
	})).Return(errors.New("db_error"))
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.CreateAuthToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorInternalError, mErr.Detail)
}

func (suite *TokenTestSuite) Test_CreateAuthToken_Success() {
	var (
		ctx  = context.Background()
		user = &pkg.User{
			Id: "user_id",
		}
		req = &pkg.CreateAuthTokenRequest{
			UserId:    user.Id,
			UserAgent: "user_agent",
			Ip:        "ip",
		}
		res = &pkg.ResponseWithAuthToken{}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("Insert", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == user.Id && input.IsActive == true && input.AccessToken != "" &&
			input.RefreshToken != "" && input.ExpireAt != nil && input.Ip == req.Ip && input.UserAgent == req.UserAgent
	})).Return(nil)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.CreateAuthToken(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.AuthToken)
	assert.NotEmpty(suite.T(), res.AuthToken.AccessToken)
	assert.NotEmpty(suite.T(), res.AuthToken.RefreshToken)
}

func (suite *TokenTestSuite) Test_RefreshAccessToken_GetByRefreshTokenDbError() {
	var (
		ctx = context.Background()
		req = &pkg.RefreshAccessTokenRequest{
			RefreshToken: "refresh_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByRefreshToken", ctx, req.RefreshToken).Return(nil, sql.ErrNoRows)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.RefreshAccessToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorAuthenticationNotFound, mErr.Detail)
}

func (suite *TokenTestSuite) Test_RefreshAccessToken_UpdateAuthLogDbError() {
	var (
		ctx            = context.Background()
		oldAccessToken = "access_token"
		authLog        = &pkg.AuthLog{
			User: &pkg.User{
				Id: "user_id",
			},
			RefreshToken: "refresh_token",
			AccessToken:  oldAccessToken,
		}
		req = &pkg.RefreshAccessTokenRequest{
			RefreshToken: authLog.RefreshToken,
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByRefreshToken", ctx, req.RefreshToken).Return(authLog, nil)
	authLogRep.On("Update", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == authLog.User.Id &&
			input.AccessToken != oldAccessToken && input.RefreshToken != ""
	})).Return(errors.New("db_error"))
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.RefreshAccessToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorInternalError, mErr.Detail)
}

func (suite *TokenTestSuite) Test_RefreshAccessToken_Success() {
	var (
		ctx            = context.Background()
		oldAccessToken = "access_token"
		authLog        = &pkg.AuthLog{
			User: &pkg.User{
				Id: "user_id",
			},
			RefreshToken: "refresh_token",
			AccessToken:  oldAccessToken,
		}
		req = &pkg.RefreshAccessTokenRequest{
			RefreshToken: authLog.RefreshToken,
		}
		res = &pkg.ResponseWithAuthToken{}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByRefreshToken", ctx, req.RefreshToken).Return(authLog, nil)
	authLogRep.On("Update", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == authLog.User.Id &&
			input.AccessToken != oldAccessToken && input.RefreshToken != ""
	})).Return(nil)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.RefreshAccessToken(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.AuthToken)
	assert.NotEmpty(suite.T(), res.AuthToken.AccessToken)
	assert.NotEmpty(suite.T(), res.AuthToken.RefreshToken)
}

func (suite *TokenTestSuite) Test_DeactivateAuthToken_GetByAccessTokenDbError() {
	var (
		ctx = context.Background()
		req = &pkg.DeactivateAuthTokenRequest{
			UserId:      "user_id",
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(nil, sql.ErrNoRows)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.DeactivateAuthToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorAuthenticationNotFound, mErr.Detail)
}

func (suite *TokenTestSuite) Test_DeactivateAuthToken_TokenOwnerInvalid() {
	var (
		ctx     = context.Background()
		authLog = &pkg.AuthLog{
			User: &pkg.User{
				Id: "user_id",
			},
		}
		req = &pkg.DeactivateAuthTokenRequest{
			UserId:      "user_id2",
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(authLog, nil)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.DeactivateAuthToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(403), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorTokenOwnerInvalid, mErr.Detail)
}

func (suite *TokenTestSuite) Test_DeactivateAuthToken_UpdateAuthLogDbError() {
	var (
		ctx     = context.Background()
		authLog = &pkg.AuthLog{
			User: &pkg.User{
				Id: "user_id",
			},
		}
		req = &pkg.DeactivateAuthTokenRequest{
			UserId:      authLog.User.Id,
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(authLog, nil)
	authLogRep.On("Update", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == authLog.User.Id && input.IsActive == false
	})).Return(errors.New("db_error"))
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.DeactivateAuthToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), pkg.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), pkg.ErrorInternalError, mErr.Detail)
}

func (suite *TokenTestSuite) Test_DeactivateAuthToken_Success() {
	var (
		ctx     = context.Background()
		authLog = &pkg.AuthLog{
			User: &pkg.User{
				Id: "user_id",
			},
		}
		req = &pkg.DeactivateAuthTokenRequest{
			UserId:      authLog.User.Id,
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(authLog, nil)
	authLogRep.On("Update", ctx, mock.MatchedBy(func(input *pkg.AuthLog) bool {
		return input.User != nil && input.User.Id == authLog.User.Id && input.IsActive == false
	})).Return(nil)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.DeactivateAuthToken(ctx, req, nil)
	assert.NoError(suite.T(), err)
}
