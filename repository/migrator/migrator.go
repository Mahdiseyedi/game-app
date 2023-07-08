package migrator

import (
	"database/sql"
	"fmt"
	"game-app/repository/mysql"
	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect    string
	dbConfig   mysql.Config
	migrations *migrate.FileMigrationSource
}

func New(dbConfig mysql.Config) Migrator {
	migrator := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}
	return Migrator{
		dialect:    "mysql",
		dbConfig:   dbConfig,
		migrations: migrator,
	}
}

//TODO - Add limit to up and down migrations
//TODO - set migrations table name

func (m Migrator) Up() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DbName))
	if err != nil {
		panic(err)
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Applied %d migrations!", n)
}

func (m Migrator) Down() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.Username, m.dbConfig.Password, m.dbConfig.Host, m.dbConfig.Port, m.dbConfig.DbName))
	if err != nil {
		panic(err)
	}

	n, err := migrate.Exec(db, m.dialect, m.migrations, migrate.Down)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Rollback %d migrations!", n)
}

func (m Migrator) Status() {
	//TODO - status
}
