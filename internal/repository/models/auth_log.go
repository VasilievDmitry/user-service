package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lotproject/go-proto/go/user_service"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		CreatedAt:    time.Now(),
		UpdatedAt:    sql.NullTime{Time: time.Now(), Valid: true},
	}

	if in.ExpireAt == nil {
		return nil, errors.New("expire time cannot be empty")
	}

	if in.Id != 0 {
		out.Id = in.Id
	}

	if in.CreatedAt != nil {
		out.CreatedAt = in.CreatedAt.AsTime()
	}
	out.ExpireAt = in.ExpireAt.AsTime()
	fmt.Println("MapProtoToModel", in.ExpireAt.Seconds)
	fmt.Println("MapProtoToModel", out.ExpireAt.Format(time.RFC3339))
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
		ExpireAt:     timestamppb.New(in.ExpireAt),
		CreatedAt:    timestamppb.New(in.CreatedAt),
		UpdatedAt:    timestamppb.New(in.UpdatedAt.Time),
	}
	fmt.Println("MapModelToProto", in.ExpireAt.Format(time.RFC3339))
	fmt.Println("MapModelToProto", out.ExpireAt.Seconds)
	if in.User == nil {
		return nil, errors.New("user cannot be empty")
	}

	user, err := NewUserMapper().MapModelToProto(in.User)
	if err != nil {
		return nil, err
	}

	out.User = user.(*user_service.User)

	if in.ExpireAt.IsZero() {
		return nil, errors.New("expire time cannot be empty")
	}

	if in.CreatedAt.IsZero() {
		return nil, errors.New("created time cannot be empty")
	}

	if !in.UpdatedAt.Valid {
		return nil, errors.New("updated time cannot be empty")
	}

	return out, nil
}
