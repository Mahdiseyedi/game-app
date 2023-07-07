package authService

import (
	"game-app/entity/user"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Service struct {
	signKey               string
	accessExpirationTime  time.Duration
	refreshExpirationTime time.Duration
	accessSubject         string
	refreshSubject        string
}

type AuthParser interface {
	VerifyToken(bearerToken string) (*Claims, error)
}

func New(signkey, accessSubject, refreshSubject string,
	accessExpirationTime, refreshExpirationTime time.Duration) Service {
	return Service{
		signKey:               signkey,
		accessExpirationTime:  accessExpirationTime,
		refreshExpirationTime: refreshExpirationTime,
		accessSubject:         accessSubject,
		refreshSubject:        refreshSubject,
	}
}

func (s Service) CreateAccessToken(user user.User) (string, error) {
	return s.CreateToken(user.ID, s.accessSubject, s.accessExpirationTime)
}

func (s Service) CreateRefreshToken(user user.User) (string, error) {
	return s.CreateToken(user.ID, s.refreshSubject, s.refreshExpirationTime)
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {
	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.signKey), nil
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
	tokenString, err := accessToken.SignedString([]byte(s.signKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
