package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessToken jwt.RegisteredClaims

func (c *AccessToken) Valid() error {
	if time.Now().Unix() > c.ExpiresAt.Unix() {
		return errors.New("jwt token is expired")
	}

	return nil
}
