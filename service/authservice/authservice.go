package authservice

import (
	"game-app/entity/user"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessSubject         string        `koanf:"access_subject"`
	RefreshSubject        string        `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}

type AuthParser interface {
	VerifyToken(bearerToken string) (*Claims, error)
}

func New(config Config) Service {
	return Service{config: config}
}

func (s Service) CreateAccessToken(user user.User) (string, error) {
	return s.CreateToken(user.ID, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) CreateRefreshToken(user user.User) (string, error) {
	return s.CreateToken(user.ID, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {
	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) CreateToken(userID uint, subject string, expiredDuration time.Duration) (string, error) {
	// create a signer for rsa 256
	// TODO - replace with rsa 256 RS256 - https://github.com/golang-jwt/jwt/blob/main/http_example_test.go

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredDuration)),
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
