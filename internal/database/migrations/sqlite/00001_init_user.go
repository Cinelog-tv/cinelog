package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {

	fmt.Println("Running migration: 1_init_user.go")

	ctbRole := sqlbuilder.NewCreateTableBuilder()

	ctbRole.CreateTable("role").IfNotExists()
	ctbRole.Define("id", "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	ctbRole.Define("name", "TEXT", "UNIQUE", "NOT NULL")

	query, _ := ctbRole.Build()

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create role table: %w", err)
	}

	ctbUser := sqlbuilder.NewCreateTableBuilder()
	ctbUser.CreateTable("user").IfNotExists()
	ctbUser.Define("id", "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	ctbUser.Define("email", "TEXT", "UNIQUE", "NOT NULL")
	ctbUser.Define("username", "TEXT", "UNIQUE", "NOT NULL")
	ctbUser.Define("password", "TEXT", "NOT NULL")
	ctbUser.Define("role_id", "INTEGER", "NOT NULL", "REFERENCES role(id)")

	query, _ = ctbUser.Build()

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}

	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS user"); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS role"); err != nil {
		return err
	}

	return nil

}
