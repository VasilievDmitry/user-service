package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-helpers/log"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/internal/repository/models"
	"go.uber.org/zap"
)

type authProviderRepository dbRepository

// AuthProviderRepositoryInterface is abstraction layer for working with providers of users and its representation in database.
type AuthProviderRepositoryInterface interface {
	// Insert adds the user provider to the collection.
	Insert(ctx context.Context, auth *user_service.AuthProvider) error

	// Update updates the user provider to the collection.
	Update(ctx context.Context, auth *user_service.AuthProvider) error

	// GetByToken returns the user provider by token.
	GetByToken(ctx context.Context, provider, token string) (*user_service.AuthProvider, error)
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

func (r *authProviderRepository) Insert(ctx context.Context, auth *user_service.AuthProvider) error {
	model, err := r.mapper.MapProtoToModel(auth)

	if err != nil {
		r.logger.Error(
			log.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(log.ErrorDatabaseFieldQuery, auth),
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

	if err != nil {
		r.logger.Error(
			log.ErrorDatabaseCreateStmt,
			zap.Error(err),
			zap.Any(log.ErrorDatabaseFieldQuery, model),
		)
		return err
	}

	_, err = r.db.NamedExecContext(ctx, query, model)

	if err != nil {
		r.logger.Error(
			log.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(log.ErrorDatabaseFieldOperation, log.ErrorDatabaseFieldOperationInsert),
			zap.Any(log.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(log.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	return nil
}

func (r *authProviderRepository) Update(ctx context.Context, auth *user_service.AuthProvider) error {
	model, err := r.mapper.MapProtoToModel(auth)

	if err != nil {
		r.logger.Error(
			log.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(log.ErrorDatabaseFieldQuery, auth),
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
			log.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(log.ErrorDatabaseFieldOperation, log.ErrorDatabaseFieldOperationInsert),
			zap.Any(log.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(log.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	return nil
}

func (r *authProviderRepository) GetByToken(ctx context.Context, provider, token string) (*user_service.AuthProvider, error) {
	var (
		model = models.AuthProvider{}
		query = fmt.Sprintf(r.getMainSelectQuery(), "WHERE provider=? AND token=?")
	)

	err := r.db.GetContext(ctx, &model, query, provider, token)

	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(
				log.ErrorDatabaseQueryFailed,
				zap.Error(err),
				zap.Any(log.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			)
		}
		return nil, err
	}

	obj, err := r.mapper.MapModelToProto(&model)

	if err != nil {
		r.logger.Error(
			log.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(log.ErrorDatabaseFieldQuery, model),
		)
		return nil, err
	}

	return obj.(*user_service.AuthProvider), nil
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
