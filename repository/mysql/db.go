package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string
	Password string
	Port     int
	Host     string
	DbName   string
}

type MySQLDB struct {
	Config Config
	db     *sql.DB
}

func New(config Config) *MySQLDB {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		config.Username, config.Password, config.Host, config.Port, config.DbName))
	if err != nil {
		panic(fmt.Errorf("cant open database, %v", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{Config: config, db: db}
}
