package config

import (
	MySql "game-app/repository/mysql"
	"game-app/service/authService"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HttpConfig HTTPServer
	Auth       authService.Config
	Mysql      MySql.Config
}
