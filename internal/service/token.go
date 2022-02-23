package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/hash"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"time"
)

func (s *Service) CreateAuthToken(
	ctx context.Context,
	req *user_service.CreateAuthTokenRequest,
	res *user_service.ResponseWithAuthToken,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	accessToken, err := createJwtToken(
		user.Id,
		s.cfg.AccessTokenLifetime,
		s.cfg.AccessTokenSigningMethod,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {
		return err
	}

	refreshToken, err := hash.GetSha256HashString(
		fmt.Sprintf("%s%s", user.Id, random.RandomString(10)),
	)

	refreshExp := time.Now().Add(time.Hour * 24 * time.Duration(s.cfg.RefreshTokenLifetime))
	expTs, err := ptypes.TimestampProto(refreshExp)

	authLog := &user_service.AuthLog{
		User:         user,
		IsActive:     true,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpireAt:     expTs,
		Ip:           req.Ip,
		UserAgent:    req.UserAgent,
	}

	if err = s.repositories.AuthLog.Insert(ctx, authLog); err != nil {
		return err
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
		return err
	}

	user, err := s.repositories.User.GetById(ctx, authLog.User.Id)

	if err != nil {
		return err
	}

	accessToken, err := createJwtToken(
		user.Id,
		s.cfg.AccessTokenLifetime,
		s.cfg.AccessTokenSigningMethod,
		s.cfg.AccessTokenSecret,
	)

	if err != nil {
		return err
	}

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return err
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
		return err
	}

	if req.User.Id != authLog.User.Id {
		return errors.New(user_service.ErrorDifferentTokenOwner)
	}

	authLog.IsActive = false

	if err = s.repositories.AuthLog.Update(ctx, authLog); err != nil {
		return err
	}

	return nil
}
