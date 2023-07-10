package config

import (
	MySql "game-app/repository/mysql"
	"game-app/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	Mysql      MySql.Config
}
