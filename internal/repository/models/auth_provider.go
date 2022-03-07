package models

import (
	"database/sql"
	"errors"
	"github.com/lotproject/go-proto/go/user_service"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		UserId:    in.User.Id,
		Provider:  in.Provider,
		Token:     in.Token,
		CreatedAt: time.Now(),
		UpdatedAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	if in.Id != 0 {
		out.Id = in.Id
	}

	if in.CreatedAt != nil {
		out.CreatedAt = in.CreatedAt.AsTime()
	}

	return out, nil
}

func (m *authProviderMapper) MapModelToProto(obj interface{}) (interface{}, error) {
	var err error

	in := obj.(*AuthProvider)
	out := &user_service.AuthProvider{
		Id:        in.Id,
		Provider:  in.Provider,
		Token:     in.Token,
		CreatedAt: timestamppb.New(in.CreatedAt),
		UpdatedAt: timestamppb.New(in.UpdatedAt.Time),
	}

	if in.User == nil {
		return nil, errors.New("user cannot be empty")
	}

	user, err := NewUserMapper().MapModelToProto(in.User)
	if err != nil {
		return nil, err
	}

	out.User = user.(*user_service.User)

	if in.CreatedAt.IsZero() {
		return nil, errors.New("created time cannot be empty")
	}

	if !in.UpdatedAt.Valid {
		return nil, errors.New("updated time cannot be empty")
	}

	return out, nil
}
