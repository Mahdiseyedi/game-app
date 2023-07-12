package backofficeuserhandler

import (
	"game-app/service/authorizeservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
)

type Handler struct {
	authConfig        authservice.Config
	authSvc           authservice.Service
	authorizationSvc  authorizeservice.Service
	backofficeUserSvc backofficeuserservice.Service
}

func New(authConfig authservice.Config, authSvc authservice.Service,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizeservice.Service) Handler {
	return Handler{
		authConfig:        authConfig,
		authSvc:           authSvc,
		authorizationSvc:  authorizationSvc,
		backofficeUserSvc: backofficeUserSvc,
	}
}
