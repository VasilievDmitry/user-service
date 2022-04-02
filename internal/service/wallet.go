package service

import (
	"context"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/user-service/pkg"
	"github.com/micro/go-micro/errors"
)

func (s *Service) CreateUserByWallet(
	ctx context.Context,
	req *pkg.CreateUserByWalletRequest,
	res *pkg.ResponseWithUserProfile,
) error {
	var user = &pkg.User{}

	if !pkg.IsSupportedWalletType(req.Provider) {
		return errors.BadRequest(pkg.ServiceName, pkg.ErrorWalletUnsupportedType)
	}

	authProvider, err := s.repositories.AuthProvider.GetByToken(ctx, req.Provider, req.Token)

	if err != nil {
		if !dbHelper.IsNotFound(err) {
			return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
		}

		user = &pkg.User{
			IsActive: true,
		}

		if err = s.repositories.User.Insert(ctx, user); err != nil {
			return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
		}

		authProvider = &pkg.AuthProvider{
			User:     user,
			Provider: req.Provider,
			Token:    req.Token,
		}

		if err = s.repositories.AuthProvider.Insert(ctx, authProvider); err != nil {
			return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
		}
	} else {
		user = authProvider.User
	}

	res.UserProfile = s.convertUserToProfile(user)

	return nil
}
