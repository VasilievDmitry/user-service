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
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *Service
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
	suite.service = NewService(repository.InitRepositories(nil, nil), suite.cfg, log)
}

func (suite *UserTestSuite) SetupTest() {
	suite.service.cfg.BcryptCost = bcrypt.MinCost
}

func (suite *UserTestSuite) TearDownTest() {
}

func (suite *UserTestSuite) Test_GetUserById_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.GetUserByIdRequest{
			UserId: "user_id",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserById(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_GetUserById_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			IsActive: true,
		}
		profile = suite.service.convertUserToProfile(user)
		req     = &user_service.GetUserByIdRequest{
			UserId: user.Id,
		}
		res = &user_service.ResponseWithUserProfile{}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserById(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.UserProfile)
	assert.Equal(suite.T(), res.UserProfile, profile)
}

func (suite *UserTestSuite) Test_GetUserByLogin_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.GetUserByLoginRequest{
			Login: "login",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetByLogin", ctx, req.Login).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserByLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_GetUserByLogin_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			Login:    "login",
			IsActive: true,
		}
		profile = suite.service.convertUserToProfile(user)
		req     = &user_service.GetUserByLoginRequest{
			Login: user.Login,
		}
		res = &user_service.ResponseWithUserProfile{}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetByLogin", ctx, req.Login).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserByLogin(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.UserProfile)
	assert.Equal(suite.T(), res.UserProfile, profile)
}

func (suite *UserTestSuite) Test_GetUserByAccessToken_GetByAccessTokenDbError() {
	var (
		ctx = context.Background()
		req = &user_service.GetUserByAccessTokenRequest{
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(nil, sql.ErrNoRows)
	suite.service.repositories.AuthLog = authLogRep

	err := suite.service.GetUserByAccessToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorAuthenticationNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_GetUserByAccessToken_GetUserDbError() {
	var (
		ctx     = context.Background()
		authLog = &user_service.AuthLog{
			User: &user_service.User{
				Id: "user_id",
			},
		}
		req = &user_service.GetUserByAccessTokenRequest{
			AccessToken: "access_token",
		}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(authLog, nil)
	suite.service.repositories.AuthLog = authLogRep

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, authLog.User.Id).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserByAccessToken(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_GetUserByAccessToken_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			Login:    "login",
			IsActive: true,
		}
		profile = suite.service.convertUserToProfile(user)
		authLog = &user_service.AuthLog{
			User: user,
		}
		req = &user_service.GetUserByAccessTokenRequest{
			AccessToken: "access_token",
		}
		res = &user_service.ResponseWithUserProfile{}
	)

	authLogRep := &mocks.AuthLogRepositoryInterface{}
	authLogRep.On("GetByAccessToken", ctx, req.AccessToken).Return(authLog, nil)
	suite.service.repositories.AuthLog = authLogRep

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, authLog.User.Id).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.GetUserByAccessToken(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.UserProfile)
	assert.Equal(suite.T(), res.UserProfile, profile)
}

func (suite *UserTestSuite) Test_SetUsername_GetUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			Username: "username",
		}
		req = &user_service.SetUsernameRequest{
			UserId:   user.Id,
			Username: user.Username,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.SetUsername(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetUsername_UpdateUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			Username: "username",
		}
		req = &user_service.SetUsernameRequest{
			UserId:   user.Id,
			Username: "username2",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.Username == req.Username
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.SetUsername(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetUsername_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:       "user_id",
			Username: "username",
		}
		req = &user_service.SetUsernameRequest{
			UserId:   user.Id,
			Username: "username2",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.Username == req.Username
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.SetUsername(ctx, req, nil)
	assert.NoError(suite.T(), err)
}

func (suite *UserTestSuite) Test_SetLogin_GetUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			EmailConfirmed: true,
		}
		req = &user_service.SetLoginRequest{
			UserId: user.Id,
			Login:  user.Login,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.SetLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetLogin_AlreadyConfirmed() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			EmailConfirmed: true,
		}
		req = &user_service.SetLoginRequest{
			UserId: user.Id,
			Login:  "login",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.SetLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorLoginAlreadyConfirmed, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetLogin_UpdateUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			EmailConfirmed: false,
		}
		req = &user_service.SetLoginRequest{
			UserId: user.Id,
			Login:  "login_new",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.Login == req.Login && input.EmailCode != ""
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.SetLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetLogin_UpdateDuplicateEntryError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			EmailConfirmed: false,
		}
		req = &user_service.SetLoginRequest{
			UserId: user.Id,
			Login:  "login_new",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.Login == req.Login && input.EmailCode != ""
	})).Return(errors.New("Duplicate entry"))
	suite.service.repositories.User = userRep

	err := suite.service.SetLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(409), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorLoginAlreadyExists, mErr.Detail)
}

func (suite *UserTestSuite) Test_SetLogin_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			EmailConfirmed: false,
		}
		req = &user_service.SetLoginRequest{
			UserId: user.Id,
			Login:  "login_new",
		}
		res = &user_service.SetLoginResponse{}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.Login == req.Login && input.EmailCode != ""
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.SetLogin(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.Code)
}

func (suite *UserTestSuite) Test_ConfirmLogin_GetUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			EmailCode:      "code",
			EmailConfirmed: true,
		}
		req = &user_service.ConfirmLoginRequest{
			UserId: user.Id,
			Code:   user.EmailCode,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.ConfirmLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *UserTestSuite) Test_ConfirmLogin_AlreadyConfirmed() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			EmailCode:      "code",
			EmailConfirmed: true,
		}
		req = &user_service.ConfirmLoginRequest{
			UserId: user.Id,
			Code:   user.EmailCode,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.ConfirmLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorLoginAlreadyConfirmed, mErr.Detail)
}

func (suite *UserTestSuite) Test_ConfirmLogin_InvalidCode() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			EmailCode:      "code",
			EmailConfirmed: false,
		}
		req = &user_service.ConfirmLoginRequest{
			UserId: user.Id,
			Code:   "code2",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.ConfirmLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorRecoveryCodeInvalid, mErr.Detail)
}

func (suite *UserTestSuite) Test_ConfirmLogin_UpdateUserDbError() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			EmailCode:      "code",
			EmailConfirmed: false,
		}
		req = &user_service.ConfirmLoginRequest{
			UserId: user.Id,
			Code:   user.EmailCode,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.EmailCode == "" && input.EmailConfirmed == true
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.ConfirmLogin(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *UserTestSuite) Test_ConfirmLogin_Success() {
	var (
		ctx  = context.Background()
		user = &user_service.User{
			Id:             "user_id",
			EmailCode:      "code",
			EmailConfirmed: false,
		}
		req = &user_service.ConfirmLoginRequest{
			UserId: user.Id,
			Code:   user.EmailCode,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == user.Id && input.EmailCode == "" && input.EmailConfirmed == true
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.ConfirmLogin(ctx, req, nil)
	assert.NoError(suite.T(), err)
}
