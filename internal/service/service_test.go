package service

import (
	"database/sql"
	"errors"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	microErrors "github.com/micro/go-micro/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	cfg     *config.Config
	service *Service
}

func Test_ServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupSuite() {
	var err error

	suite.cfg, err = config.NewConfig()
	if err != nil {
		suite.FailNow("Config load failed", err)
	}

	log, _ := zap.NewProduction()
	suite.service = NewService(nil, suite.cfg, log)
}

func (suite *ServiceTestSuite) SetupTest() {
}

func (suite *ServiceTestSuite) TearDownTest() {
}

func (suite *ServiceTestSuite) Test_NewService() {
	srv := NewService(nil, nil, nil)
	assert.IsType(suite.T(), (*Service)(nil), srv)
}

func (suite *ServiceTestSuite) Test_Ping() {
	err := suite.service.Ping(nil, nil, nil)
	assert.NoError(suite.T(), err)
}

func (suite *ServiceTestSuite) Test_buildGetUserError_NotFound() {
	err := microErrors.Parse(suite.service.buildGetUserError(sql.ErrNoRows).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(404), err.Code)
	assert.Equal(suite.T(), user_service.ErrorUserNotFound, err.Detail)
}

func (suite *ServiceTestSuite) Test_buildGetUserError_InternalServerError() {
	err := microErrors.Parse(suite.service.buildGetUserError(errors.New("error")).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(500), err.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, err.Detail)
}

func (suite *ServiceTestSuite) Test_buildGetWalletError_NotFound() {
	err := microErrors.Parse(suite.service.buildGetWalletError(sql.ErrNoRows).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(404), err.Code)
	assert.Equal(suite.T(), user_service.ErrorWalletNotFound, err.Detail)
}

func (suite *ServiceTestSuite) Test_buildGetWalletError_InternalServerError() {
	err := microErrors.Parse(suite.service.buildGetWalletError(errors.New("error")).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(500), err.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, err.Detail)
}

func (suite *ServiceTestSuite) Test_buildGetAuthLogError_NotFound() {
	err := microErrors.Parse(suite.service.buildGetAuthLogError(sql.ErrNoRows).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(404), err.Code)
	assert.Equal(suite.T(), user_service.ErrorAuthenticationNotFound, err.Detail)
}

func (suite *ServiceTestSuite) Test_buildGetAuthLogError_InternalServerError() {
	err := microErrors.Parse(suite.service.buildGetAuthLogError(errors.New("error")).Error())
	assert.NotEmpty(suite.T(), err)
	assert.Equal(suite.T(), user_service.ServiceName, err.Id)
	assert.Equal(suite.T(), int32(500), err.Code)
	assert.Equal(suite.T(), user_service.ErrorInternalError, err.Detail)
}

func (suite *ServiceTestSuite) Test_createJwtToken_InvalidAlg() {
	assert.Panics(suite.T(), func() { _, _ = createJwtToken("user_id", 10, "", "") })
}

func (suite *ServiceTestSuite) Test_createJwtToken() {
	token, err := createJwtToken("user_id", 0, "HS256", "")
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

func (suite *ServiceTestSuite) Test_convertUserToProfile() {
	var (
		user = &user_service.User{
			Id:             "user_id",
			Login:          "login",
			Username:       "username",
			IsActive:       true,
			EmailConfirmed: true,
		}
	)

	profile := suite.service.convertUserToProfile(user)
	assert.Equal(suite.T(), user.Id, profile.Id)
	assert.Equal(suite.T(), user.Login, profile.Login)
	assert.Equal(suite.T(), user.Username, profile.Username)
	assert.Equal(suite.T(), user.IsActive, profile.IsActive)
	assert.Equal(suite.T(), user.EmailConfirmed, profile.EmailConfirmed)
}
