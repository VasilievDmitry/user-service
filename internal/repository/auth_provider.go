package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	dbHelper "github.com/lotproject/go-helpers/db"
	"go.uber.org/zap"

	"github.com/lotproject/user-service/internal/repository/models"
	userService "github.com/lotproject/user-service/proto/v1"
)

type authProviderRepository dbRepository

// AuthProviderRepositoryInterface is abstraction layer for working with providers of users and its representation in database.
type AuthProviderRepositoryInterface interface {
	// Insert adds the user provider to the collection.
	Insert(ctx context.Context, provider *userService.AuthProvider) error

	// Update updates the user provider to the collection.
	Update(ctx context.Context, provider *userService.AuthProvider) error

	// GetByToken returns the user provider by token.
	GetByToken(ctx context.Context, provider, token string) (*userService.AuthProvider, error)

	// GetByUserId returns the user provider by user identifier.
	GetByUserId(ctx context.Context, provider string) ([]*userService.AuthProvider, error)
}

// NewAuthProviderRepository create and return an object for working with the providers of user repository.
// The returned object implements the NewAuthProviderRepository interface.
func NewAuthProviderRepository(db *sqlx.DB, logger *zap.Logger) AuthProviderRepositoryInterface {
	s := &authProviderRepository{
		db:     db,
		logger: logger,
		mapper: models.NewAuthProviderMapper(),
	}
	return s
}

func (r *authProviderRepository) Insert(ctx context.Context, provider *userService.AuthProvider) error {
	model, err := r.mapper.MapProtoToModel(provider)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, provider),
		)
		return err
	}

	var query = `
		INSERT INTO auth_provider (
			user_id, 
			provider, 
			token, 
			created_at, 
			updated_at
		) VALUES (
			:user_id, 
			:provider, 
			:token, 
			:created_at, 
			:updated_at
		)`
	res, err := r.db.NamedExecContext(ctx, query, model)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(dbHelper.ErrorDatabaseFieldOperation, dbHelper.ErrorDatabaseFieldOperationInsert),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(dbHelper.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	id, err := res.LastInsertId()

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseGetLatestId,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(dbHelper.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	provider.Id = id

	return nil
}

func (r *authProviderRepository) Update(ctx context.Context, provider *userService.AuthProvider) error {
	model, err := r.mapper.MapProtoToModel(provider)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, provider),
		)
		return err
	}

	var query = `
		UPDATE 
			auth_provider 
		SET 
			user_id=:user_id, 
			provider=:provider, 
			token=:token, 
			created_at=:created_at, 
			updated_at=:updated_at 
		WHERE id=:id`
	_, err = r.db.NamedExecContext(ctx, query, model)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(dbHelper.ErrorDatabaseFieldOperation, dbHelper.ErrorDatabaseFieldOperationInsert),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(dbHelper.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	return nil
}

func (r *authProviderRepository) GetByToken(ctx context.Context, provider, token string) (*userService.AuthProvider, error) {
	var (
		model = models.AuthProvider{}
		query = fmt.Sprintf(r.getMainSelectQuery(), "WHERE provider=? AND token=?")
	)

	err := r.db.GetContext(ctx, &model, query, provider, token)

	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(
				dbHelper.ErrorDatabaseQueryFailed,
				zap.Error(err),
				zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			)
		}
		return nil, err
	}

	obj, err := r.mapper.MapModelToProto(&model)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, model),
		)
		return nil, err
	}

	return obj.(*userService.AuthProvider), nil
}

func (r *authProviderRepository) GetByUserId(ctx context.Context, userId string) ([]*userService.AuthProvider, error) {
	var (
		list  []*models.AuthProvider
		query = fmt.Sprintf(r.getMainSelectQuery(), "WHERE user_id=?")
	)

	err := r.db.SelectContext(ctx, &list, query, userId)

	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(
				dbHelper.ErrorDatabaseQueryFailed,
				zap.Error(err),
				zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			)
		}
		return nil, err
	}

	objs := make([]*userService.AuthProvider, len(list))

	for i, obj := range list {
		v, err := r.mapper.MapModelToProto(obj)
		if err != nil {
			r.logger.Error(
				dbHelper.ErrorDatabaseMapModelFailed,
				zap.Error(err),
				zap.Any(dbHelper.ErrorDatabaseFieldQuery, obj),
			)
			return nil, err
		}
		objs[i] = v.(*userService.AuthProvider)
	}

	return objs, nil
}

func (r *authProviderRepository) getMainSelectQuery() string {
	return `
		SELECT 
			ap.*,
			u.id              AS 'user.id',
			u.login           AS 'user.login',
			u.username        AS 'user.username',
			u.email_confirmed AS 'user.email_confirmed',
			u.is_active       AS 'user.is_active'
		FROM auth_provider AS ap
		LEFT OUTER JOIN user AS u ON u.id = ap.user_id
			%s`
}
