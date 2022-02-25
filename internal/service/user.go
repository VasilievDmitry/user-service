package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/micro/go-micro/errors"
	"strconv"
)

func (s *Service) GetUserById(
	ctx context.Context,
	req *user_service.GetUserByIdRequest,
	res *user_service.ResponseWithUserProfile,
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
	req *user_service.GetUserByLoginRequest,
	res *user_service.ResponseWithUserProfile,
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
	req *user_service.GetUserByAccessTokenRequest,
	res *user_service.ResponseWithUserProfile,
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

func (s *Service) SetUsername(
	ctx context.Context,
	req *user_service.SetUsernameRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	user.Username = req.Username

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	return nil
}

func (s *Service) SetLogin(
	ctx context.Context,
	req *user_service.SetLoginRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	if user.EmailConfirmed {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorLoginAlreadyConfirmed)
	}

	user.Login = req.Login
	user.EmailCode = strconv.Itoa(random.RandomInt(100000, 999999))

	if err = s.repositories.User.Update(ctx, user); err != nil {
		if dbHelper.IsDuplicateEntry(err) {
			return errors.BadRequest(user_service.ServiceName, user_service.ErrorLoginAlreadyExists)
		}

		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	// TODO: Отправить письмо с подтверждением (проверить лимит)

	return nil
}

func (s *Service) ConfirmLogin(
	ctx context.Context,
	req *user_service.ConfirmLoginRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	if user.EmailConfirmed {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorLoginAlreadyConfirmed)
	}

	if user.EmailCode != req.Code {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorRecoveryCodeInvalid)
	}

	user.EmailCode = ""
	user.EmailConfirmed = true

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	return nil
}
