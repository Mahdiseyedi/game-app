package config

import (
	"game-app/adapter/redis"
	MySql "game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/matchingservice"
	"time"
)

type Application struct {
	GracefulShutdownTimeout time.Duration `koanf:"graceful_shutdown_timeout"`
}

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	Application     Application            `koanf:"application"`
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           MySql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
}
