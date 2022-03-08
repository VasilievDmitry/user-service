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

type PasswordTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *Service
}

func Test_PasswordTestSuite(t *testing.T) {
	suite.Run(t, new(PasswordTestSuite))
}

func (suite *PasswordTestSuite) SetupSuite() {
	var err error

	suite.cfg, err = config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", err)
	}

	log, _ := zap.NewProduction()
	suite.service = NewService(repository.InitRepositories(nil, nil), suite.cfg, log)
}

func (suite *PasswordTestSuite) SetupTest() {
}

func (suite *PasswordTestSuite) TearDownTest() {
}

func (suite *PasswordTestSuite) Test_VerifyPassword_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.VerifyPasswordRequest{
			UserId:   "user_id",
			Password: "password",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.VerifyPassword(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_VerifyPassword_CompareError() {
	var (
		ctx = context.Background()
		req = &user_service.VerifyPasswordRequest{
			UserId:   "user_id",
			Password: "password",
		}
		user = &user_service.User{
			Id:       "user_id",
			Password: "password",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.VerifyPassword(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInvalidPassword, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_VerifyPassword_CompareSuccess() {
	var (
		ctx = context.Background()
		req = &user_service.VerifyPasswordRequest{
			UserId:   "user_id",
			Password: "password",
		}
		password, _ = bcrypt.GenerateFromPassword([]byte(req.Password), suite.cfg.BcryptCost)
		user        = &user_service.User{
			Id:       "user_id",
			Password: string(password),
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.VerifyPassword(ctx, req, nil)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordTestSuite) Test_SetPassword_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.SetPasswordRequest{
			UserId:   "user_id",
			Password: "password",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.SetPassword(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_SetPassword_UpdateUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.SetPasswordRequest{
			UserId: "user_id",
		}
		user = &user_service.User{
			Id: req.UserId,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.Password != ""
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.SetPassword(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_SetPassword_Success() {
	var (
		ctx = context.Background()
		req = &user_service.SetPasswordRequest{
			UserId: "user_id",
		}
		user = &user_service.User{
			Id: req.UserId,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.Password != ""
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.SetPassword(ctx, req, nil)
	assert.NoError(suite.T(), err)
}

func (suite *PasswordTestSuite) Test_CreatePasswordRecoveryCode_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.CreatePasswordRecoveryCodeRequest{
			UserId: "user_id",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.CreatePasswordRecoveryCode(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_CreatePasswordRecoveryCode_UpdateUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.CreatePasswordRecoveryCodeRequest{
			UserId: "user_id",
		}
		user = &user_service.User{
			Id: req.UserId,
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.RecoveryCode != ""
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.CreatePasswordRecoveryCode(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_CreatePasswordRecoveryCode_Success() {
	var (
		ctx = context.Background()
		req = &user_service.CreatePasswordRecoveryCodeRequest{
			UserId: "user_id",
		}
		user = &user_service.User{
			Id: req.UserId,
		}
		res = &user_service.CreatePasswordRecoveryCodeResponse{}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.RecoveryCode != ""
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.CreatePasswordRecoveryCode(ctx, req, res)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), res)
	assert.NotEmpty(suite.T(), res.Code)
}

func (suite *PasswordTestSuite) Test_UsePasswordRecoveryCode_GetUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.UsePasswordRecoveryCodeRequest{
			UserId: "user_id",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(nil, sql.ErrNoRows)
	suite.service.repositories.User = userRep

	err := suite.service.UsePasswordRecoveryCode(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(404), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_UsePasswordRecoveryCode_InvalidRecoveryCode() {
	var (
		ctx = context.Background()
		req = &user_service.UsePasswordRecoveryCodeRequest{
			UserId: "user_id",
			Code:   "code",
		}
		user = &user_service.User{
			Id:           req.UserId,
			RecoveryCode: "recovery_code",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	suite.service.repositories.User = userRep

	err := suite.service.UsePasswordRecoveryCode(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(400), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorRecoveryCodeInvalid, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_UsePasswordRecoveryCode_UpdateUserDbError() {
	var (
		ctx = context.Background()
		req = &user_service.UsePasswordRecoveryCodeRequest{
			UserId: "user_id",
			Code:   "code",
		}
		user = &user_service.User{
			Id:           req.UserId,
			RecoveryCode: "code",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.RecoveryCode == "" && input.Password != ""
	})).Return(errors.New("db_error"))
	suite.service.repositories.User = userRep

	err := suite.service.UsePasswordRecoveryCode(ctx, req, nil)
	assert.Error(suite.T(), err)

	mErr := microErrors.Parse(err.Error())
	assert.NotEmpty(suite.T(), mErr)
	assert.Equal(suite.T(), user_service.ServiceName, mErr.Id)
	assert.Equal(suite.T(), int32(500), mErr.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, mErr.Detail)
}

func (suite *PasswordTestSuite) Test_UsePasswordRecoveryCode_Success() {
	var (
		ctx = context.Background()
		req = &user_service.UsePasswordRecoveryCodeRequest{
			UserId: "user_id",
			Code:   "code",
		}
		user = &user_service.User{
			Id:           req.UserId,
			RecoveryCode: "code",
		}
	)

	userRep := &mocks.UserRepositoryInterface{}
	userRep.On("GetById", ctx, req.UserId).Return(user, nil)
	userRep.On("Update", ctx, mock.MatchedBy(func(input *user_service.User) bool {
		return input.Id == req.UserId && input.RecoveryCode == "" && input.Password != ""
	})).Return(nil)
	suite.service.repositories.User = userRep

	err := suite.service.UsePasswordRecoveryCode(ctx, req, nil)
	assert.NoError(suite.T(), err)
}
