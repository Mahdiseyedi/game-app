package matchinghandler

import (
	middlewares "game-app/delivery/httpserver/middleware"
	"github.com/labstack/echo/v4"
)

func (h Handler) SetRoute(e *echo.Echo) {
	userGroup := e.Group("/matching")

	userGroup.POST("/add-to-waiting-list", h.addToWaitingList,
		middlewares.Auth(h.authSvc, h.authConfig))
}
