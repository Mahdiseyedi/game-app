package config

import (
	"game-app/adapter/redis"
	MySql "game-app/repository/mysql"
	"game-app/service/authservice"
	"game-app/service/matchingservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer      HTTPServer             `koanf:"http_server"`
	Auth            authservice.Config     `koanf:"auth"`
	Mysql           MySql.Config           `koanf:"mysql"`
	MatchingService matchingservice.Config `koanf:"matching_service"`
	Redis           redis.Config           `koanf:"redis"`
}
