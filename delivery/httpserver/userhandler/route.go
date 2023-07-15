package userhandler

import (
	middlewares "game-app/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/users")

	userGroup.POST("/register", h.userRegister)
	userGroup.POST("/login", h.userLogin)
	userGroup.GET("/profile", h.userProfile,
		middlewares.Auth(h.authSvc, h.authConfig),
		middlewares.UpsertPresence(h.presenceSvc))
}
