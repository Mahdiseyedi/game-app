package main

import (
	"context"
	"fmt"
	"game-app/Validator/matchingvalidator"
	"game-app/Validator/uservalidator"
	presenceClient "game-app/adapter/presence"
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/delivery/httpserver"
	"game-app/repository/migrator"
	"game-app/repository/mysql"
	"game-app/repository/mysql/mysqlaccesscontrol"
	"game-app/repository/mysql/mysqluser"
	"game-app/repository/redis/redismatching"
	"game-app/repository/redis/redispresence"
	"game-app/scheduler"
	"game-app/service/authorizationservice"
	"game-app/service/authservice"
	"game-app/service/backofficeuserservice"
	"game-app/service/matchingservice"
	"game-app/service/presenceservice"
	"game-app/service/userservice"
	_ "google.golang.org/grpc"
	"os"
	"os/signal"
	"sync"
	"time"
)

const (
	JwtSignKey = "jwt_secret"
)

func main() {
	cfg := config.Load("config.yml")
	fmt.Printf("cfg: %+v\n", cfg)

	mgr := migrator.New(cfg.Mysql)
	//TODO - add command for migrations
	mgr.Up()

	//TODO - add struct and add these returned items as struct field
	authSvc, userSvc, userValidator, backofficeSvc, authorizationSvc,
		matchingSvc, matchingV, presencSvc := SetupServices(cfg)

	server := httpserver.New(cfg, authSvc, userSvc, userValidator,
		backofficeSvc, authorizationSvc, matchingSvc, matchingV, presencSvc)

	go func() {
		server.Serve()
	}()

	done := make(chan bool)
	var wg sync.WaitGroup
	go func() {
		sch := scheduler.New(cfg.Scheduler, matchingSvc)

		wg.Add(1)
		sch.Start(done, &wg)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx := context.Background()
	ctxWithTimeout, cancel := context.WithTimeout(ctx,
		cfg.Application.GracefulShutdownTimeout)
	defer cancel()

	if err := server.Router.Shutdown(ctxWithTimeout); err != nil {
		fmt.Println("http server shutdown error: ", err)
	}

	fmt.Println("received interrupt signal, shutting down gracefully...")
	done <- true
	time.Sleep(cfg.Application.GracefulShutdownTimeout)
	<-ctxWithTimeout.Done()

	wg.Wait()
}

func SetupServices(cfg config.Config) (
	authservice.Service, userservice.Service, uservalidator.Validator,
	backofficeuserservice.Service, authorizationservice.Service,
	matchingservice.Service, matchingvalidator.Validator,
	presenceservice.Service,
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

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	presenceAdapter := presenceClient.New(":8086")
	matchingSvc := matchingservice.New(cfg.MatchingService, matchingRepo, presenceAdapter, redisAdapter)

	return authSvc, userSvc, uV, backofficeUserSvc,
		authorizationSvc, matchingSvc, matchingV, presenceSvc
}
