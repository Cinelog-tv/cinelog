package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDefaultRole, downDefaultRole)
}

func upDefaultRole(ctx context.Context, tx *sql.Tx) error {
	fmt.Println("Running migration: 2_default_role.go")

	sb := sqlbuilder.SQLite.NewInsertBuilder()
	sb.InsertInto("role")
	sb.Cols("name")
	sb.Values("admin")
	sb.Values("user")

	query, args := sb.Build()

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("failed to insert default roles: %w", err)
	}

	return nil
}

func downDefaultRole(ctx context.Context, tx *sql.Tx) error {
	query, _ := sqlbuilder.DeleteFrom("role").Build()
	if _, err := tx.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to delete default roles: %w", err)
	}

	return nil
}
