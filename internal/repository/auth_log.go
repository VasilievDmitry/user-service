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

type authLogRepository dbRepository

// AuthLogRepositoryInterface is abstraction layer for working with authorizations of users and its representation in database.
type AuthLogRepositoryInterface interface {
	// Insert adds the user auth to the collection.
	Insert(ctx context.Context, auth *user_service.AuthLog) error

	// Update updates the user auth to the collection.
	Update(ctx context.Context, auth *user_service.AuthLog) error

	// GetByAccessToken returns the user auth by access token.
	GetByAccessToken(ctx context.Context, token string) (*user_service.AuthLog, error)

	// GetByRefreshToken returns the user auth by refresh token.
	GetByRefreshToken(ctx context.Context, token string) (*user_service.AuthLog, error)
}

// NewAuthLogRepository create and return an object for working with the authorizations of user repository.
// The returned object implements the AuthLogRepositoryInterface interface.
func NewAuthLogRepository(db *sqlx.DB, logger *zap.Logger) AuthLogRepositoryInterface {
	s := &authLogRepository{
		db:     db,
		logger: logger,
		mapper: models.NewAuthLogMapper(),
	}
	return s
}

func (r *authLogRepository) Insert(ctx context.Context, auth *user_service.AuthLog) error {
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
		INSERT INTO auth_log (
			user_id, 
			ip, 
			user_agent, 
			access_token, 
			refresh_token, 
			is_active, 
			created_at, 
			updated_at, 
			expire_at
		) VALUES (
			:user_id, 
			:ip, 
			:user_agent, 
			:access_token, 
			:refresh_token, 
			:is_active, 
			:created_at, 
			:updated_at, 
			:expire_at
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

func (r *authLogRepository) Update(ctx context.Context, auth *user_service.AuthLog) error {
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
			auth_log 
		SET 
			user_id=:user_id, 
			ip=:ip, 
			user_agent=:user_agent, 
			access_token=:access_token, 
			refresh_token=:refresh_token, 
			is_active=:is_active, 
			expire_at=:expire_at, 
			created_at=:created_at 
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

func (r *authLogRepository) GetByAccessToken(ctx context.Context, token string) (*user_service.AuthLog, error) {
	var (
		model = models.AuthLog{}
		query = fmt.Sprintf(r.getMainSelectQuery(), "WHERE al.access_token=? AND al.is_active=1")
	)

	err := r.db.GetContext(ctx, &model, query, token)

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

	return obj.(*user_service.AuthLog), nil
}

func (r *authLogRepository) GetByRefreshToken(ctx context.Context, token string) (*user_service.AuthLog, error) {
	var (
		model = models.AuthLog{}
		query = fmt.Sprintf(r.getMainSelectQuery(), "WHERE al.refresh_token=? AND al.is_active=1")
	)

	err := r.db.GetContext(ctx, &model, query, token)

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

	return obj.(*user_service.AuthLog), nil
}

func (r *authLogRepository) getMainSelectQuery() string {
	return `
		SELECT 
			al.*,
			u.id              AS 'user.id',
			u.login           AS 'user.login',
			u.username        AS 'user.username',
			u.email_confirmed AS 'user.email_confirmed',
			u.is_active       AS 'user.is_active'
		FROM auth_log AS al
		LEFT OUTER JOIN user AS u ON u.id = al.user_id
			%s`
}