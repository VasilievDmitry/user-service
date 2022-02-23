package models

import (
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/lotproject/go-proto/go/user_service"
	"time"
)

type authLogMapper struct{}

func NewAuthLogMapper() Mapper {
	return &authLogMapper{}
}

type AuthLog struct {
	Id           int64        `db:"id"`
	UserId       string       `db:"user_id"`
	User         *User        `db:"user"`
	Ip           string       `db:"ip"`
	UserAgent    string       `db:"user_agent"`
	AccessToken  string       `db:"access_token"`
	RefreshToken string       `db:"refresh_token"`
	IsActive     bool         `db:"is_active"`
	ExpireAt     time.Time    `db:"expire_at"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}

func (m *authLogMapper) MapProtoToModel(obj interface{}) (interface{}, error) {
	in := obj.(*user_service.AuthLog)
	out := &AuthLog{
		UserId:       in.User.Id,
		Ip:           in.Ip,
		UserAgent:    in.UserAgent,
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
		IsActive:     in.IsActive,
	}

	if in.Id != 0 {
		out.Id = in.Id
	}

	if in.CreatedAt != nil {
		t, err := ptypes.Timestamp(in.CreatedAt)
		if err != nil {
			return nil, err
		}

		out.CreatedAt = t
	} else {
		out.CreatedAt = time.Now()
	}

	t, err := ptypes.Timestamp(in.ExpireAt)
	if err != nil {
		return nil, err
	}

	out.ExpireAt = t

	out.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return out, nil
}

func (m *authLogMapper) MapModelToProto(obj interface{}) (interface{}, error) {
	var err error

	in := obj.(*AuthLog)
	out := &user_service.AuthLog{
		Id:           in.Id,
		Ip:           in.Ip,
		UserAgent:    in.UserAgent,
		AccessToken:  in.AccessToken,
		RefreshToken: in.RefreshToken,
		IsActive:     in.IsActive,
	}

	if in.User == nil {
		return nil, errors.New("user cannot be empty")
	}

	user, err := NewUserMapper().MapModelToProto(in.User)
	if err != nil {
		return nil, err
	}

	out.User = user.(*user_service.User)

	out.CreatedAt, err = ptypes.TimestampProto(in.CreatedAt)

	if err != nil {
		return nil, err
	}

	out.ExpireAt, err = ptypes.TimestampProto(in.ExpireAt)

	if err != nil {
		return nil, err
	}

	if in.UpdatedAt.Valid {
		out.UpdatedAt, err = ptypes.TimestampProto(in.UpdatedAt.Time)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
