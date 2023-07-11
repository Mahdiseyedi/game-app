package main

import (
	"fmt"
	"game-app/Validator/uservalidator"
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/migrator"
	"game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/userservice"
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
	cfgz := config.Load("config.yml")
	fmt.Printf("cfgz: +%v\n", cfgz)

	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8088},
		Auth: authservice.Config{
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

	mgr := migrator.New(cfg.Mysql)
	//TODO - add command for migrations
	mgr.Up()

	authSvc, userSvc, userValidator := SetupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	server.Serve()
}

func SetupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)
	MysqlRepo := mysql.New(cfg.Mysql)
	userSvc := userservice.New(MysqlRepo, authSvc)
	uV := uservalidator.New(MysqlRepo)

	return authSvc, userSvc, uV
}
