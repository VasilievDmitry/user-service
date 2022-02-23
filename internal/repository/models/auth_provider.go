package models

import (
	"database/sql"
	"errors"
	"github.com/golang/protobuf/ptypes"
	"github.com/lotproject/go-proto/go/user_service"
	"time"
)

type authProviderMapper struct{}

func NewAuthProviderMapper() Mapper {
	return &authProviderMapper{}
}

type AuthProvider struct {
	Id        int64        `db:"id"`
	UserId    string       `db:"user_id"`
	User      *User        `db:"user"`
	Provider  string       `db:"provider"`
	Token     string       `db:"token"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

func (m *authProviderMapper) MapProtoToModel(obj interface{}) (interface{}, error) {
	in := obj.(*user_service.AuthProvider)
	out := &AuthProvider{
		UserId:   in.User.Id,
		Provider: in.Provider,
		Token:    in.Token,
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

	out.UpdatedAt = sql.NullTime{Time: time.Now(), Valid: true}

	return out, nil
}

func (m *authProviderMapper) MapModelToProto(obj interface{}) (interface{}, error) {
	var err error

	in := obj.(*AuthProvider)
	out := &user_service.AuthProvider{
		Id:       in.Id,
		Provider: in.Provider,
		Token:    in.Token,
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
