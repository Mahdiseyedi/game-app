package middleware

import (
	cfg "game-app/config"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"game-app/service/authservice"
	mw "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func Auth(service authservice.Service, config authservice.Config) echo.MiddlewareFunc {
	return mw.WithConfig(mw.Config{
		ContextKey:    cfg.AuthMiddlewareContextKey,
		SigningKey:    []byte(config.SignKey),
		SigningMethod: "HS256",
		ParseTokenFunc: func(c echo.Context, auth string) (interface{}, error) {
			const op = "middleware.Auth.ParseTokenFunc"

			claims, err := service.VerifyToken(auth)
			if err != nil {
				return nil, richerror.New(op).WithErr(err).
					WithKind(richerror.KindInvalid).
					WithMessage(errmsg.ErrorMsgSomethingWentWrong)
			}

			return claims, nil
		},
	})
}
