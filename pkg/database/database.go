package database

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"go-kafka-example/config"
)

func ConnectDB(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	switch cfg.Adapter {
	case "sqlite3":
		return connectSqlite3(cfg)
	case "mysql":
		return connectMySQL(cfg)
	case "postgres":
		return connectPostgres(cfg)
	default:
		return nil, errors.New("Database: " + cfg.Adapter + " not exist")
	}
}

func connectSqlite3(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	return sqlx.Open("sqlite3", cfg.Host)

}

func connectMySQL(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	setting := fmt.Sprintf("%s:%s@tcp(%s):%d/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	return sqlx.Open("mysql", setting)
}

func connectPostgres(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	setting := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Db)
	return sqlx.Open("postgres", setting)
}
