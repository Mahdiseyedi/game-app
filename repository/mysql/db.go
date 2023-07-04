package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db *sql.DB
}

func New() *DB {

	db, err := sql.Open("mysql", "gameapp:gameappt0lk2o20@(localhost:3308)/gameapp_db")
	if err != nil {
		panic(fmt.Errorf("cant open database, %v", err))
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &DB{db}
}
