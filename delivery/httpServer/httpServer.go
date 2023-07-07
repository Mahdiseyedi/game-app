package httpServer

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	userGroup := e.Group("/users")

	userGroup.POST("/register", s.UserRegister)
	userGroup.POST("/login", s.Login)

	e.GET("/health-check", s.healthCheckHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", s.config.HttpConfig.Port)))
}
