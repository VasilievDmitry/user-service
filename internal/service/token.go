package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/hash"
	"github.com/lotproject/go-helpers/random"
	"go-micro.dev/v4/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/lotproject/user-service/pkg"
	userService "github.com/lotproject/user-service/proto/v1"
)

func (s *Service) CreateAuthToken(
	ctx context.Context,
	req *userService.CreateAuthTokenRequest,
	res *userService.ResponseWithAuthToken,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	accessToken, err := createJwtToken(
		user.Id,
		s.cfg.AccessTokenLifetime,
		jwt.SigningMethodHS256,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {
		s.log.Error("Unable to create JWT token", zap.Error(err))
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	refreshToken, err := hash.GetSha256HashString(
		fmt.Sprintf("%s%s", user.Id, random.RandomString(10)),
	)

	if err != nil {
		s.log.Error("Unable to generate refresh token hash", zap.Error(err))
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	refreshExp := time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenLifetime))

	authLog := &userService.AuthLog{
		User:         user,
		IsActive:     true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     timestamppb.New(refreshExp),
		Ip:           req.Ip,
		UserAgent:    req.UserAgent,
	}

	if err = s.repositories.AuthLog.Insert(ctx, authLog); err != nil {
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	res.AuthToken = &userService.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return nil
}

func (s *Service) RefreshAccessToken(
	ctx context.Context,
	req *userService.RefreshAccessTokenRequest,
	res *userService.ResponseWithAuthToken,
) error {
	authLog, err := s.repositories.AuthLog.GetByRefreshToken(ctx, req.RefreshToken)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	accessToken, err := createJwtToken(
		authLog.User.Id,
		s.cfg.AccessTokenLifetime,
		jwt.SigningMethodHS256,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {
		s.log.Error("Unable to create JWT token", zap.Error(err))
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	authLog.AccessToken = accessToken

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	res.AuthToken = &userService.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: authLog.RefreshToken,
	}

	return nil
}

func (s *Service) DeactivateAuthToken(
	ctx context.Context,
	req *userService.DeactivateAuthTokenRequest,
	_ *empty.Empty,
) error {
	authLog, err := s.repositories.AuthLog.GetByAccessToken(ctx, req.AccessToken)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	if req.UserId != authLog.User.Id {
		return errors.Forbidden(pkg.ServiceName, pkg.ErrorTokenOwnerInvalid)
	}

	authLog.IsActive = false

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	return nil
}
