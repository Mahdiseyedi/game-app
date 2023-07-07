package httpServer

import (
	"game-app/config"
	"game-app/service/authService"
	"game-app/service/userService"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	config  config.Config
	authSvc authService.Service
	userSvc userService.Service
}

func New(config config.Config,
	authSvc authService.Service,
	userSvc userService.Service) Server {

	return Server{
		config:  config,
		authSvc: authSvc,
		userSvc: userSvc,
	}
}

func (s Server) UserRegister(c echo.Context) error {
	var rReq userService.RegisterRequest

	if err := c.Bind(&rReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	resp, err := s.userSvc.Register(rReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, resp)
}

func (s Server) Login(c echo.Context) error {
	var req userService.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	resp, err := s.userSvc.Login(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

func (s Server) UserProfile(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")

	claims, err := s.authSvc.VerifyToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	resp, err := s.userSvc.GetUserProfile(userService.UserProfileRequest{UserID: claims.UserID})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
