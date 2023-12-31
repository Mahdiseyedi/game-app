package mysql

import (
	"database/sql"
	"fmt"
	"game-app/pkg/errmsg"
	"game-app/pkg/richerror"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Config struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Port     int    `koanf:"port"`
	Host     string `koanf:"host"`
	DBName   string `koanf:"db_name"`
}

type MySQLDB struct {
	Config Config
	db     *sql.DB
}

func (m *MySQLDB) Conn() *sql.DB {
	return m.db
}

func New(config Config) *MySQLDB {
	const op = "mysql.New"
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s",
		config.Username, config.Password, config.Host, config.Port, config.DBName))
	if err != nil {
		panic(richerror.New(op).WithErr(err).
			WithKind(richerror.KindUnexpected).WithMessage(errmsg.ErrorMsgCantOpenDatabase))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySQLDB{Config: config, db: db}
}
