package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/lotproject/go-helpers/hash"
	"github.com/lotproject/go-helpers/random"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/micro/go-micro/errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) VerifyPassword(
	ctx context.Context,
	req *user_service.VerifyPasswordRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))

	if err != nil {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorInvalidPassword)
	}

	return nil
}

func (s *Service) SetPassword(
	ctx context.Context,
	req *user_service.SetPasswordRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)

	if err != nil {
		s.log.Error("Unable to generate bcrypt password", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	user.Password = string(password)

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	return nil
}

func (s *Service) CreatePasswordRecoveryCode(
	ctx context.Context,
	req *user_service.CreatePasswordRecoveryCodeRequest,
	res *user_service.CreatePasswordRecoveryCodeResponse,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	user.RecoveryCode, err = hash.GetSha256HashString(random.RandomString(10))

	if err != nil {
		s.log.Error("Unable to create recovery code", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	res.Code = user.RecoveryCode

	return nil
}

func (s *Service) UsePasswordRecoveryCode(
	ctx context.Context,
	req *user_service.UsePasswordRecoveryCodeRequest,
	_ *empty.Empty,
) error {
	user, err := s.repositories.User.GetById(ctx, req.UserId)

	if err != nil {
		return s.buildGetUserError(err)
	}

	if user.RecoveryCode != req.Code {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorRecoveryCodeInvalid)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)

	if err != nil {
		s.log.Error("Unable to generate bcrypt password", zap.Error(err))
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	user.Password = string(password)
	user.RecoveryCode = ""

	if err = s.repositories.User.Update(ctx, user); err != nil {
		return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
	}

	return nil
}
