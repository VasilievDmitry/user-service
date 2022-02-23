package service

import (
	"context"
	"errors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"strconv"
)

func (s *Service) GetUserById(
	ctx context.Context,
	req *user_service.GetUserByIdRequest,
	res *user_service.ResponseWithUserProfile,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return err
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
		return err
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
		return err
	}

	user, err := s.repositories.User.GetById(ctx, authLog.User.Id)

	if err != nil {
		return err
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}

func (s *Service) SetUsername(
	ctx context.Context,
	req *user_service.SetUsernameRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	user.Username = req.Username

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) SetLogin(
	ctx context.Context,
	req *user_service.SetLoginRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	if user.EmailConfirmed {
		return errors.New(user_service.ErrorLoginAlreadyConfirmed)
	}

	check, err := s.repositories.User.GetByLogin(ctx, req.Login)

	if err != nil {
		return err
	}

	if check != nil {
		return errors.New(user_service.ErrorLoginAlreadyExists)
	}

	user.Login = req.Login
	user.EmailCode = strconv.Itoa(random.RandomInt(100000, 999999))

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	// TODO: Отправить письмо с подтверждением (проверить лимит)

	return nil
}

func (s *Service) ConfirmLogin(
	ctx context.Context,
	req *user_service.ConfirmLoginRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	if user.EmailConfirmed {
		return errors.New(user_service.ErrorLoginAlreadyConfirmed)
	}

	if user.EmailCode != req.Code {
		return errors.New(user_service.ErrorInvalidRecoveryCode)
	}

	user.EmailCode = ""
	user.EmailConfirmed = true

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
