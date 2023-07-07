package httpServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s Server) healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "everything is good!",
	})
}
