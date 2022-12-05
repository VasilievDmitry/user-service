package service

import (
	"context"

	dbHelper "github.com/lotproject/go-helpers/db"
	"go-micro.dev/v4/errors"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lotproject/user-service/pkg"
	userService "github.com/lotproject/user-service/proto/v1"
)

func (s *Service) GetSupportedWallets(
	ctx context.Context,
	_ *emptypb.Empty,
	res *userService.GetSupportedWalletsResponse,
) error {
	var wallets []string

	for wallet, _ := range pkg.WalletList {
		wallets = append(wallets, wallet)
	}

	res.Wallets = wallets

	return nil
}

func (s *Service) CreateUserByWallet(
	ctx context.Context,
	req *userService.CreateUserByWalletRequest,
	res *userService.ResponseWithUserProfile,
) error {
	var user = &userService.User{}

	if !pkg.IsSupportedWalletType(req.Provider) {
		return errors.BadRequest(pkg.ServiceName, pkg.ErrorWalletUnsupportedType)
	}

	authProvider, err := s.repositories.AuthProvider.GetByToken(ctx, req.Provider, req.Token)

	if err != nil {
		if !dbHelper.IsNotFound(err) {
			return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
		}

		user = &userService.User{
			IsActive: true,
		}

		if err = s.repositories.User.Insert(ctx, user); err != nil {
			return errors.InternalServerError(pkg.ServiceName, pkg.ErrorInternalError)
		}

		authProvider = &userService.AuthProvider{
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
