package service

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/golang/protobuf/ptypes/empty"
	gameService "github.com/lotproject/game-service/pkg"
	dbHelper "github.com/lotproject/go-helpers/db"
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
	gameService  gameService.GameService
}

// NewService
func NewService(
	repositories *repository.Repositories,
	gameService gameService.GameService,
	cfg *config.Config,
	log *zap.Logger,
) *Service {
	return &Service{
		repositories: repositories,
		gameService:  gameService,
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

func (s *Service) convertUserToProfile(user *pkg.User) *pkg.UserProfile {
	centrifugoToken, _ := createJwtToken(user.Id, s.cfg.RefreshTokenLifetime, jwt.SigningMethodHS256, s.cfg.CentrifugoSecret)

	servers, err := s.gameService.GetUserServers(context.TODO(), &gameService.GetUserServersRequest{UserId: user.Id})
	if err != nil {
		s.log.Error(
			"Unable to get server list for user profile",
			zap.Error(err),
		)
	}

	var serverList []*pkg.GameServer

	if servers != nil && len(servers.List) > 0 {
		for _, server := range servers.List {
			serverList = append(serverList, &pkg.GameServer{
				Id:   server.Id,
				Name: server.Name,
			})
		}
	}

	if servers != nil && len(servers.List) > 0 {
		for _, server := range serverList {
			serverList = append(serverList, &pkg.GameServer{
				Id:   server.Id,
				Name: server.Name,
			})
		}
	}

	var walletList []*pkg.AuthProvider

	wallets, _ := s.repositories.AuthProvider.GetByUserId(context.TODO(), user.Id)
	for _, wallet := range wallets {
		walletList = append(walletList, &pkg.AuthProvider{
			Provider: wallet.Provider,
			Token:    wallet.Token,
		})
	}

	profile := &pkg.UserProfile{
		Id:                user.Id,
		Login:             user.Login,
		Username:          user.Username,
		EmailConfirmed:    user.EmailConfirmed,
		CentrifugoToken:   centrifugoToken,
		CentrifugoChannel: fmt.Sprintf(s.cfg.CentrifugoUserChannel, user.Id),
		GameServers:       serverList,
		Wallets:           walletList,
	}

	return profile
}

func (s *Service) buildGetUserError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(pkg.ServiceName, pkg.ErrorUserNotFound)
	}

	return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
}

func (s *Service) buildGetWalletError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(pkg.ServiceName, pkg.ErrorWalletNotFound)
	}

	return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
}

func (s *Service) buildGetAuthLogError(err error) error {
	if dbHelper.IsNotFound(err) {
		return errors.NotFound(pkg.ServiceName, pkg.ErrorAuthenticationNotFound)
	}

	return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
}

func createJwtToken(userId string, lifeTime int, signingMethod jwt.SigningMethod, secret string) (string, error) {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(lifeTime)).Unix(),
		Subject:   userId,
	}
	token := jwt.NewWithClaims(signingMethod, claims)
	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
