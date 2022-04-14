package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/user-service/pkg"
	"github.com/micro/go-micro/errors"
	"go.uber.org/zap"
	"strconv"
)

func (s *Service) GetUserById(
	ctx context.Context,
	req *pkg.GetUserByIdRequest,
	res *pkg.ResponseWithUserProfile,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}

func (s *Service) GetUserByLogin(
	ctx context.Context,
	req *pkg.GetUserByLoginRequest,
	res *pkg.ResponseWithUserProfile,
) error {
	user, err := s.repositories.User.GetByLogin(ctx, req.Login)

	if err != nil {
		return s.buildGetUserError(err)
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}

func (s *Service) GetUserByAccessToken(
	ctx context.Context,
	req *pkg.GetUserByAccessTokenRequest,
	res *pkg.ResponseWithUserProfile,
) error {
	authLog, err := s.repositories.AuthLog.GetByAccessToken(ctx, req.AccessToken)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	user, err := s.repositories.User.GetById(ctx, authLog.User.Id)

	if err != nil {
		return s.buildGetUserError(err)
	}

	res.UserProfile = s.convertUserToProfile(user)
	s.log.Info("GetUserByAccessToken", zap.Any("profile", res.UserProfile))
	return nil
}

func (s *Service) SetUsername(
	ctx context.Context,
	req *pkg.SetUsernameRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	user.Username = req.Username

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	return nil
}

func (s *Service) SetLogin(
	ctx context.Context,
	req *pkg.SetLoginRequest,
	res *pkg.SetLoginResponse,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	if user.EmailConfirmed {
		return errors.BadRequest(pkg.ServiceName, pkg.ErrorLoginAlreadyConfirmed)
	}

	user.Login = req.Login
	user.EmailCode = strconv.Itoa(random.RandomInt(100000, 999999))

	if err = s.repositories.User.Update(ctx, user); err != nil {
		if dbHelper.IsDuplicateEntry(err) {
			return errors.Conflict(pkg.ServiceName, pkg.ErrorLoginAlreadyExists)
		}

		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	res.Code = user.EmailCode

	return nil
}

func (s *Service) ConfirmLogin(
	ctx context.Context,
	req *pkg.ConfirmLoginRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	if user.EmailConfirmed {
		return errors.BadRequest(pkg.ServiceName, pkg.ErrorLoginAlreadyConfirmed)
	}

	if user.EmailCode != req.Code {
		return errors.BadRequest(pkg.ServiceName, pkg.ErrorRecoveryCodeInvalid)
	}

	user.EmailCode = ""
	user.EmailConfirmed = true

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
	}

	return nil
}
