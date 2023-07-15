package httpserver

import (
	"fmt"
	"game-app/Validator/matchingvalidator"
	"game-app/Validator/uservalidator"
	"game-app/config"
	"game-app/delivery/httpserver/backofficeuserhandler"
	"game-app/delivery/httpserver/matchinghandler"
	"game-app/delivery/httpserver/userhandler"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/userservice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config            config.Config
	userHandler       userhandler.Handler
	backofficeHandler backofficeuserhandler.Handler
	matchingHandler   matchinghandler.Handler
	Router            *echo.Echo
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service,
	userValidator uservalidator.Validator,
	backofficeUserSvc backofficeuserservice.Service, authorizationSvc authorizationservice.Service,
	matchingSvc matchingservice.Service,
	matchingValidator matchingvalidator.Validator) Server {

	return Server{
		Router:            echo.New(),
		config:            config,
		userHandler:       userhandler.New(config.Auth, authSvc, userSvc, userValidator),
		backofficeHandler: backofficeuserhandler.New(config.Auth, authSvc, backofficeUserSvc, authorizationSvc),
		matchingHandler:   matchinghandler.New(config.Auth, authSvc, matchingSvc, matchingValidator),
	}
}

func (s Server) Serve() {
	//middlewares
	s.Router.Use(middleware.Logger())
	s.Router.Use(middleware.Recover())

	//Routes
	s.Router.GET("/health-check", s.healthCheck)

	s.userHandler.SetRoutes(s.Router)
	s.backofficeHandler.SetRoutes(s.Router)
	s.matchingHandler.SetRoute(s.Router)

	//Start server
	address := fmt.Sprintf(":%d", s.config.HTTPServer.Port)
	fmt.Printf("start echo server on %s\n", address)

	if err := s.Router.Start(address); err != nil {
		fmt.Println("route start error: ", err)
	}
}
