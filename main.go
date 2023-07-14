package main

import (
	"fmt"
	"game-app/Validator/matchingvalidator"
	"game-app/Validator/uservalidator"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/migrator"
	"game-app/repository/mysql"
	"game-app/repository/mysql/mysqlaccesscontrol"
	"game-app/repository/mysql/mysqluser"
	"game-app/repository/redis/redismatching"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/userservice"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: +%v\n", cfg)

	mgr := migrator.New(cfg.Mysql)
	//TODO - add command for migrations
	mgr.Up()

	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc,
		matchingSvc, matchingV := SetupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator,
		backofficeSvc, authorizationSvc, matchingSvc, matchingV)

	server.Serve()
}

func SetupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
) {
	authSvc := authservice.New(cfg.Auth)
	MysqlRepo := mysql.New(cfg.Mysql)
	userMysql := mysqluser.New(MysqlRepo)
	userSvc := userservice.New(userMysql, authSvc)
	backofficeUserSvc := backofficeuserservice.New()

	aclMysql := mysqlaccesscontrol.New(MysqlRepo)
	authorizationSvc := authorizationservice.New(aclMysql)

	uV := uservalidator.New(userMysql)

	matchingV := matchingvalidator.New()
	redisAdapter := redis.New(cfg.Redis)
	matchingRepo := redismatching.New(redisAdapter)
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo)

	return authSvc, userSvc, uV, backofficeUserSvc, authorizationSvc, matchingSvc, matchingV
}
