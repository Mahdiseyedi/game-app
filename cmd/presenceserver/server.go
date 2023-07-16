package main

import (
	"game-app/adapter/redis"
	"game-app/config"
	"game-app/delivery/grpcserver/presenceserver"
	"game-app/repository/redis/redispresence"
	"game-app/service/presenceservice"
)

func main() {
	cfg := config.Load("config.yml")

	redisAdapter := redis.New(cfg.Redis)

	presenceRepo := redispresence.New(redisAdapter)
	presenceSvc := presenceservice.New(cfg.PresenceService, presenceRepo)

	server := presenceserver.New(presenceSvc)
	server.Start()
}
