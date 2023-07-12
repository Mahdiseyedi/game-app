package backofficeuserhandler

import (
	middlewares "game-app/delivery/httpserver/middleware"
	"game-app/entity/permission"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoutes(e *echo.Echo) {
	userGroup := e.Group("/backoffice/users")

	userGroup.GET("/", h.listUsers, middlewares.Auth(h.authSvc, h.authConfig),
		middlewares.AccessCheck(h.authorizationSvc, permission.UserListPermission))
}
