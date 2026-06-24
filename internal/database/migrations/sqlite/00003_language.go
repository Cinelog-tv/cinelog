package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upLanguage, downLanguage)
}

func upLanguage(ctx context.Context, tx *sql.Tx) error {
	fmt.Println("Running migration: 3_language.go")

	sbLanguage := sqlbuilder.SQLite.NewCreateTableBuilder()
	sbLanguage.CreateTable("language").IfNotExists()
	sbLanguage.Define("id", "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	sbLanguage.Define("code", "TEXT", "NOT NULL", "UNIQUE")

	query, args := sbLanguage.Build()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}

	return nil
}

func downLanguage(ctx context.Context, tx *sql.Tx) error {

	if _, err := tx.ExecContext(ctx, "DELETE TABLE IF EXISTS language"); err != nil {
		return err
	}

	return nil
}
