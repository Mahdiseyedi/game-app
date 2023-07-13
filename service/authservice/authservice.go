package authservice

import (
	"game-app/entity/role"
	"game-app/entity/user"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
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
	return s.CreateToken(user.ID, user.Role, s.config.AccessSubject, s.config.AccessExpirationTime)
}

func (s Service) VerifyToken(bearerToken string) (*Claims, error) {
	const op = "authservice.VerifyToken"
	bearerToken = strings.Replace(bearerToken, "Bearer ", "", 1)
	token, err := jwt.ParseWithClaims(bearerToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindForbidden).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, richerror.New(op).WithErr(err).
			WithKind(richerror.KindInvalid).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong)
	}
}

func (s Service) CreateRefreshToken(user user.User) (string, error) {
	return s.CreateToken(user.ID, user.Role, s.config.RefreshSubject, s.config.RefreshExpirationTime)
}

func (s Service) CreateToken(userID uint, role role.Role,
	subject string, expiredDuration time.Duration) (string, error) {
	const op = "authservice.CreateToken"

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiredDuration)),
		},
		UserID: userID,
		Role:   role,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).
			WithMessage(errmsg.ErrorMsgSomethingWentWrong)
	}

	return tokenString, nil
}
