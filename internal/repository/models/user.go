package models

import (
	"database/sql"
	"time"

	"github.com/golang/protobuf/ptypes"

	userService "github.com/lotproject/user-service/proto/v1"
)

type userMapper struct{}

func NewUserMapper() Mapper {
	return &userMapper{}
}

type User struct {
	Id             sql.NullString `db:"id"`
	Login          sql.NullString `db:"login"`
	Password       sql.NullString `db:"password"`
	Username       sql.NullString `db:"username"`
	EmailCode      sql.NullString `db:"email_code"`
	EmailConfirmed sql.NullBool   `db:"email_confirmed"`
	RecoveryCode   sql.NullString `db:"recovery_code"`
	IsActive       sql.NullBool   `db:"is_active"`
	CreatedAt      sql.NullTime   `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	Balance      	sql.NullFloat64   `db:"balance"`
}

func (m *userMapper) MapProtoToModel(obj interface{}) (interface{}, error) {
	in := obj.(*userService.User)
	out := &User{
		EmailConfirmed: sql.NullBool{Bool: in.EmailConfirmed, Valid: true},
		IsActive:       sql.NullBool{Bool: in.IsActive, Valid: true},
	}

	if in.Id != "" {
		out.Id = sql.NullString{String: in.Id, Valid: true}
	}

	if in.Login != "" {
		out.Login = sql.NullString{String: in.Login, Valid: true}
	}

	if in.Password != "" {
		out.Password = sql.NullString{String: in.Password, Valid: true}
	}

	if in.Username != "" {
		out.Username = sql.NullString{String: in.Username, Valid: true}
	}

	if in.EmailCode != "" {
		out.EmailCode = sql.NullString{String: in.EmailCode, Valid: true}
	}

	if in.RecoveryCode != "" {
		out.RecoveryCode = sql.NullString{String: in.RecoveryCode, Valid: true}
	}

	if in.CreatedAt != nil {
		t, err := ptypes.Timestamp(in.CreatedAt)
		if err != nil {
			return nil, err
		}

		out.CreatedAt = sql.NullTime{Time: t, Valid: true}
	} else {
		out.CreatedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}

	out.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return out, nil
}

func (m *userMapper) MapModelToProto(obj interface{}) (interface{}, error) {
	var err error

	in := obj.(*User)
	out := &userService.User{}

	if in.Id.Valid {
		out.Id = in.Id.String
	}

	if in.Login.Valid {
		out.Login = in.Login.String
	}

	if in.Password.Valid {
		out.Password = in.Password.String
	}

	if in.Username.Valid {
		out.Username = in.Username.String
	}

	if in.EmailCode.Valid {
		out.EmailCode = in.EmailCode.String
	}

	if in.EmailConfirmed.Valid {
		out.EmailConfirmed = in.EmailConfirmed.Bool
	}

	if in.RecoveryCode.Valid {
		out.RecoveryCode = in.RecoveryCode.String
	}

	if in.IsActive.Valid {
		out.IsActive = in.IsActive.Bool
	}
	if in.Balance.Valid {
		out.Balance = in.Balance.Float64
	}
	if in.CreatedAt.Valid {
		out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt.Time)
		if err != nil {
			return nil, err
		}
	}

	if in.UpdatedAt.Valid {
		out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt.Time)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
