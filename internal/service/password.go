package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/hash"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) VerifyPassword(
	ctx context.Context,
	req *user_service.VerifyPasswordRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return errors.New(user_service.ErrorInvalidPassword)
	}

	return nil
}

func (s *Service) SetPassword(
	ctx context.Context,
	req *user_service.SetPasswordRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)

	if err != nil {
		return err
	}

	user.Password = string(password)

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *Service) CreatePasswordRecoveryCode(
	ctx context.Context,
	req *user_service.UserProfile,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.Id)

	if err != nil {
		return err
	}

	user.RecoveryCode, err = hash.GetSha256HashString(random.RandomString(10))

	if err != nil {
		return err
	}

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	// TODO: Отправить письмо с кодом (проверить лимит)
	fmt.Println(user.RecoveryCode)

	return nil
}

func (s *Service) UsePasswordRecoveryCode(
	ctx context.Context,
	req *user_service.UsePasswordRecoveryCodeRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.User.Id)

	if err != nil {
		return err
	}

	if user.RecoveryCode != req.Code {
		return errors.New(user_service.ErrorInvalidRecoveryCode)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)

	if err != nil {
		return err
	}

	user.Password = string(password)
	user.RecoveryCode = ""

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return err
	}

	return nil
}
