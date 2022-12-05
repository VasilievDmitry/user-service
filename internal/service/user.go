package service

import (
	"context"
	"strconv"

	"github.com/golang/protobuf/ptypes/empty"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/random"
	"go-micro.dev/v4/errors"

	"github.com/lotproject/user-service/pkg"
	userService "github.com/lotproject/user-service/proto/v1"
)

func (s *Service) GetUserById(
	ctx context.Context,
	req *userService.GetUserByIdRequest,
	res *userService.ResponseWithUserProfile,
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
	req *userService.GetUserByLoginRequest,
	res *userService.ResponseWithUserProfile,
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
	req *userService.GetUserByAccessTokenRequest,
	res *userService.ResponseWithUserProfile,
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

	return nil
}

func (s *Service) GetUserByWallet(
	ctx context.Context,
	req *userService.GetUserByWalletRequest,
	res *userService.ResponseWithUserProfile,
) error {
	authProvider, err := s.repositories.AuthProvider.GetByToken(ctx, req.Provider, req.Token)

	if err != nil {
		return s.buildGetAuthLogError(err)
	}

	user, err := s.repositories.User.GetById(ctx, authProvider.User.Id)

	if err != nil {
		return s.buildGetUserError(err)
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}

func (s *Service) SetUsername(
	ctx context.Context,
	req *userService.SetUsernameRequest,
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
	req *userService.SetLoginRequest,
	res *userService.SetLoginResponse,
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
	req *userService.ConfirmLoginRequest,
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
