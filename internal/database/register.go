package database

import (
	"context"
	"database/sql"

	_ "github.com/Cinelog-tv/cinelog/internal/database/migrations/sqlite"
	"github.com/pressly/goose/v3"
)

func migrate(db *sql.DB, driver DBDriver) error {
	goose.SetDialect(string(driver))
	return goose.UpContext(context.Background(), db, ".")
}
