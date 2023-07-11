package config

import (
	MySql "game-app/repository/mysql"
	"game-app/service/authservice"
)

type HTTPServer struct {
	Port int `koanf:"port"`
}

type Config struct {
	HTTPServer HTTPServer         `koanf:"http_server"`
	Auth       authservice.Config `koanf:"auth"`
	Mysql      MySql.Config       `koanf:"mysql"`
}
