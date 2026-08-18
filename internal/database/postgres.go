package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBConfig interface {
	DBUser() string
	DBPass() string
	DBHost() string
	DBPort() string
	DBName() string
}

func GetPostgresConnectionString(cfg DBConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser(), cfg.DBPass(), cfg.DBHost(), cfg.DBPort(), cfg.DBName(),
	)
}

func ConnectPostgres(cfg DBConfig) (*sql.DB, error) {
	connStr := GetPostgresConnectionString(cfg)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)

	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Minute * 15)

	return db, nil
}
