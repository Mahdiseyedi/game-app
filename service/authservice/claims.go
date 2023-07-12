package authservice

import (
	"game-app/entity/role"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	RegisteredClaims jwt.RegisteredClaims
	UserID           uint      `json:"user_id"`
	Role             role.Role `json:"role"`
}

func (c Claims) Valid() error {
	return c.RegisteredClaims.Valid()
}
