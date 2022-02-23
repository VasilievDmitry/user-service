package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/config"
	"github.com/lotproject/user-service/internal/repository"
	"github.com/lotproject/user-service/pkg"
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

func (s *Service) convertUserToProfile(user *user_service.User) *user_service.UserProfile {
	return &user_service.UserProfile{
		Id:             user.Id,
		Login:          user.Login,
		Username:       user.Username,
		IsActive:       user.IsActive,
		EmailConfirmed: user.EmailConfirmed,
	}
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
