package main

import (
	"game-app/config"
	"game-app/delivery/httpServer"
	"game-app/repository/mysql"
	"game-app/service/authService"
	"game-app/service/userService"
	"time"
)

const (
	JwtSignKey           = "jwt_secret"
	AccessTokenSubject   = "at"
	RefreshTokenSubject  = "rt"
	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 7
)

func main() {
	cfg := config.Config{
		HttpConfig: config.HTTPServer{Port: 8088},
		Auth: authService.Config{
			SignKey:               JwtSignKey,
			AccessExpirationTime:  AccessTokenDuration,
			RefreshExpirationTime: RefreshTokenDuration,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
		},
		Mysql: mysql.Config{
			Username: "gameapp",
			Password: "gameappt0lk2o20",
			Port:     3308,
			Host:     "localhost",
			DbName:   "gameapp_db",
		},
	}

	authSvc, userSvc := SetupServices(cfg)
	server := httpServer.New(cfg, authSvc, userSvc)

	server.Serve()
}

func SetupServices(cfg config.Config) (authService.Service, userService.Service) {
	authSvc := authService.New(cfg.Auth)
	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userService.New(MysqlRepo, authSvc)

	return authSvc, userSvc
}
