package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	dbHelper "github.com/lotproject/go-helpers/db"
	"github.com/lotproject/go-proto/go/user_service"
	"github.com/lotproject/user-service/internal/repository/models"
	"go.uber.org/zap"
)

type userRepository dbRepository

// UserRepositoryInterface is abstraction layer for working with users and its representation in database.
type UserRepositoryInterface interface {
	// Insert adds the user to the collection.
	Insert(ctx context.Context, user *user_service.User) error

	// Update updates the user in the collection.
	Update(context.Context, *user_service.User) error

	// GetById returns the user by identity.
	GetById(ctx context.Context, id string) (*user_service.User, error)

	// GetByLogin returns the active user by login.
	GetByLogin(ctx context.Context, login string) (*user_service.User, error)
}

// NewUserRepository create and return an object for working with the user repository.
// The returned object implements the UserRepositoryInterface interface.
func NewUserRepository(db *sqlx.DB, logger *zap.Logger) UserRepositoryInterface {
	s := &userRepository{
		db:     db,
		logger: logger,
		mapper: models.NewUserMapper(),
	}
	return s
}

func (r *userRepository) Insert(ctx context.Context, user *user_service.User) error {
	if user.Id == "" {
		user.Id = uuid.NewString()
	}

	model, err := r.mapper.MapProtoToModel(user)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, user),
		)
		return err
	}

	query := `
		INSERT INTO user (id, login, password, username, email_code, email_confirmed, recovery_code, is_active, created_at, updated_at)
		VALUES (:id, :login, :password, :username, :email_code, :email_confirmed, :recovery_code, :is_active, :created_at, :updated_at)`
	_, err = r.db.NamedExecContext(ctx, query, model)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(dbHelper.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	return nil
}

func (r *userRepository) Update(ctx context.Context, user *user_service.User) error {
	model, err := r.mapper.MapProtoToModel(user)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, user),
		)
		return err
	}

	query := `
		UPDATE user
		SET login=:login,
			password=:password,
			username=:username,
			email_code=:email_code,
			email_confirmed=:email_confirmed,
			is_active=:is_active,
			created_at=:created_at,
			updated_at=:updated_at
		WHERE id = :id`
	_, err = r.db.NamedExecContext(ctx, query, model)

	if err != nil {
		r.logger.Error(
			dbHelper.ErrorDatabaseMapModelFailed,
			zap.Error(err),
			zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
			zap.Any(dbHelper.ErrorDatabaseFieldDocument, model),
		)
		return err
	}

	return nil
}

func (r *userRepository) GetById(ctx context.Context, id string) (*user_service.User, error) {
	var model = models.User{}

	query := r.getMainSelectQuery()
	query = fmt.Sprintf(query, "WHERE u.id = ?")

	err := r.db.GetContext(ctx, &model, query, id)

	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(
				dbHelper.ErrorDatabaseQueryFailed,
				zap.Error(err),
				zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
				zap.String("id", id),
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

	return obj.(*user_service.User), nil
}

func (r *userRepository) GetByLogin(ctx context.Context, _login string) (*user_service.User, error) {
	var model = models.User{}

	query := r.getMainSelectQuery()
	query = fmt.Sprintf(query, "WHERE u.login = ? AND u.is_active=1")

	err := r.db.GetContext(ctx, &model, query, _login)

	if err != nil {
		if err != sql.ErrNoRows {
			r.logger.Error(
				dbHelper.ErrorDatabaseQueryFailed,
				zap.Error(err),
				zap.Any(dbHelper.ErrorDatabaseFieldQuery, dbHelper.CleanQueryForLog(query)),
				zap.String("login", _login),
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

	return obj.(*user_service.User), nil
}

func (r *userRepository) getMainSelectQuery() string {
	return `
		SELECT u.*
		FROM user AS u
		%s`
}
