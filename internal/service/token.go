package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/hash"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/micro/go-micro/errors"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (s *Service) CreateAuthToken(
	ctx context.Context,
	req *user_service.CreateAuthTokenRequest,
	res *user_service.ResponseWithAuthToken,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	accessToken, err := createJwtToken(
		user.Id,
		s.cfg.AccessTokenLifetime,
		s.cfg.AccessTokenSigningMethod,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {
		s.log.Error("Unable to create JWT token", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	refreshToken, err := hash.GetSha256HashString(
		fmt.Sprintf("%s%s", user.Id, random.RandomString(10)),
	)

	if err != nil {
		s.log.Error("Unable to generate refresh token hash", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	refreshExp := time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenLifetime))

	authLog := &user_service.AuthLog{
		User:         user,
		IsActive:     true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     timestamppb.New(refreshExp),
		Ip:           req.Ip,
		UserAgent:    req.UserAgent,
	}

	if err = s.repositories.AuthLog.Insert(ctx, authLog); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	res.AuthToken = &user_service.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return nil
}

func (s *Service) RefreshAccessToken(
	ctx context.Context,
	req *user_service.RefreshAccessTokenRequest,
	res *user_service.ResponseWithAuthToken,
) error {
	authLog, err := s.repositories.AuthLog.GetByRefreshToken(ctx, req.RefreshToken)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	user, err := s.repositories.User.GetById(ctx, authLog.User.Id)

	if err != nil {
		return s.buildGetUserError(err)
	}

	accessToken, err := createJwtToken(
		user.Id,
		s.cfg.AccessTokenLifetime,
		s.cfg.AccessTokenSigningMethod,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {

		s.log.Error("Unable to create JWT token", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	res.AuthToken = &user_service.AuthToken{
		AccessToken:  accessToken,
		RefreshToken: authLog.RefreshToken,
	}

	return nil
}

func (s *Service) DeactivateAuthToken(
	ctx context.Context,
	req *user_service.DeactivateAuthTokenRequest,
	_ *empty.Empty,
) error {
	authLog, err := s.repositories.AuthLog.GetByAccessToken(ctx, req.AccessToken)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	if req.UserId != authLog.User.Id {
		return errors.Forbidden(user_service.ServiceName, user_service.ErrorTokenOwnerInvalid)
	}

	authLog.IsActive = false

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	return nil
}
