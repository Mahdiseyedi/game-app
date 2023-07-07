package authService

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	RegisteredClaims jwt.RegisteredClaims
	UserID           uint `json:"user_id"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
