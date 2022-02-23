package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AccessToken jwt.StandardClaims

func (c *AccessToken) Valid() error {
	if time.Now().Unix() > c.ExpiresAt {
		return errors.New("jwt token is expired")
	}

	return nil
}
