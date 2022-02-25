package service

import (
	"context"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/micro/go-micro/errors"
)

func (s *Service) CreateUserByWallet(
	ctx context.Context,
	req *user_service.CreateUserByWalletRequest,
	res *user_service.ResponseWithUserProfile,
) error {
	var user = &user_service.User{}

	if !user_service.IsSupportedWalletType(req.Provider) {
		return errors.BadRequest(user_service.ServiceName, user_service.ErrorWalletUnsupportedType)
	}

	authProvider, err := s.repositories.AuthProvider.GetByToken(ctx, req.Provider, req.Token)

	if err != nil {
		if !dbHelper.IsNotFound(err) {
			return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
		}

		newUser := &user_service.User{
			IsActive: true,
		}

		if err = s.repositories.User.Insert(ctx, newUser); err != nil {
			return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
		}

		authProvider = &user_service.AuthProvider{
			User:     user,
			Provider: req.Provider,
			Token:    req.Token,
		}

		if err = s.repositories.AuthProvider.Insert(ctx, authProvider); err != nil {
			return errors.InternalServerError(user_service.ServiceName, user_service.ErrorInternalError)
		}

		if user, err = s.repositories.User.GetById(ctx, newUser.Id); err != nil {
			return s.buildGetUserError(err)
		}
	} else {
		user = authProvider.User
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}
