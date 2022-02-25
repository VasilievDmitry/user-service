package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes/empty"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/pkg"
	"github.com/micro/go-micro/errors"
	"go.uber.org/zap"
	"time"
)

type Service struct {
	cfg          *config.Config
	log          *zap.Logger
	repositories *repository.Repositories
}

func NewService(
	repositories *repository.Repositories,
	cfg *config.Config,
	log *zap.Logger,
) *Service {
	return &Service{
		repositories: repositories,
		cfg:          cfg,
		log:          log,
	}
}

func (s *Service) Ping(
	_ context.Context,
	_ *empty.Empty,
	_ *empty.Empty,
) error {
	return nil
}

func (s *Service) convertUserToProfile(user *user_service.User) *user_service.UserProfile {
	return &user_service.UserProfile{
		Id:             user.Id,
		Login:          user.Login,
		Username:       user.Username,
		IsActive:       user.IsActive,
		EmailConfirmed: user.EmailConfirmed,
	}
}

func (s *Service) buildGetUserError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(user_service.ServiceName, user_service.ErrorUserNotFound)
	}

	return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
}

func (s *Service) buildGetWalletError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(user_service.ServiceName, user_service.ErrorWalletNotFound)
	}

	return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
}

func (s *Service) buildGetAuthLogError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(user_service.ServiceName, user_service.ErrorAuthenticationNotFound)
	}

	return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
}

func createJwtToken(userId string, lifeTime int, signingMethod, secret string) (string, error) {
	accessToken := &pkg.AccessToken{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(lifeTime)).Unix(),
		Subject:   userId,
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(signingMethod), accessToken)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
