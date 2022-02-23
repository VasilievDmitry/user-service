package repository

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/lotproject/user-service/internal/repository/models"
	"go.uber.org/zap"
)

type dbRepository struct {
	db     *sqlx.DB
	redis  redis.Cmdable
	logger *zap.Logger
	mapper models.Mapper
}

type Repositories struct {
	User         UserRepositoryInterface
	AuthProvider AuthProviderRepositoryInterface
	AuthLog      AuthLogRepositoryInterface
}

func InitRepositories(db *sqlx.DB, redis redis.Cmdable, log *zap.Logger) *Repositories {
	return &Repositories{
		User:         NewUserRepository(db, log),
		AuthProvider: NewAuthProviderRepository(db, log),
		AuthLog:      NewAuthLogRepository(db, log),
	}
}
