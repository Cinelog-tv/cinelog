package database

import (
	"database/sql"
	"fmt"

	_ "github.com/glebarez/go-sqlite"
)

type DBDriver string

const (
	DriverPostgres DBDriver = "postgres"
	DriverSQLite   DBDriver = "sqlite"
	DriverMySQL    DBDriver = "mysql"
)

type DatabaseConfig struct {
	Driver DBDriver
	DSN    string
}

func NewDatabase(cfg DatabaseConfig) (*sql.DB, error) {
	db, err := openDB(cfg)

	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := migrate(db, cfg.Driver); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func openDB(cfg DatabaseConfig) (*sql.DB, error) {
	switch cfg.Driver {
	case DriverSQLite:
		return sql.Open("sqlite", cfg.DSN)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
