package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/huandu/go-sqlbuilder"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upMovieInit, downMovieInit)
}

func upMovieInit(ctx context.Context, tx *sql.Tx) error {
	fmt.Println("Running migration: 4_movie.go")

	sbMovie := sqlbuilder.SQLite.NewCreateTableBuilder()
	sbMovie.CreateTable("movie").IfNotExists()
	sbMovie.Define("id", "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	sbMovie.Define("tmdb_id", "INTEGER", "UNIQUE")
	sbMovie.Define("imdb_id", "TEXT", "UNIQUE")
	sbMovie.Define("release_date", "DATE")
	sbMovie.Define("runtime", "INTEGER")
	sbMovie.Define("default_poster", "TEXT")
	sbMovie.Define("created_at", "DATETIME", "NOT NULL", "DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))")
	sbMovie.Define("updated_at", "DATETIME", "NOT NULL", "DEFAULT (strftime('%Y-%m-%dT%H:%M:%fZ', 'now'))")

	// movie_translation
	sbMovieTranslation := sqlbuilder.SQLite.NewCreateTableBuilder()
	sbMovieTranslation.CreateTable("movie_translation").IfNotExists()
	sbMovieTranslation.Define("movie_id", "INTEGER", "NOT NULL")
	sbMovieTranslation.Define("language_id", "INTEGER", "NOT NULL")
	sbMovieTranslation.Define("title", "TEXT")
	sbMovieTranslation.Define("overview", "TEXT")
	sbMovieTranslation.Define("tagline", "TEXT")
	sbMovieTranslation.Define("PRIMARY KEY (movie_id, language_id)")
	sbMovieTranslation.Define("FOREIGN KEY (movie_id)", "REFERENCES movie(id)", "ON DELETE CASCADE")
	sbMovieTranslation.Define("FOREIGN KEY (language_id)", "REFERENCES language(id)", "ON DELETE CASCADE")

	// movie_poster
	sbMoviePoster := sqlbuilder.SQLite.NewCreateTableBuilder()
	sbMoviePoster.CreateTable("movie_poster").IfNotExists()
	sbMoviePoster.Define("id", "INTEGER", "PRIMARY KEY", "AUTOINCREMENT")
	sbMoviePoster.Define("movie_id", "INTEGER", "NOT NULL")
	sbMoviePoster.Define("language_id", "INTEGER")
	sbMoviePoster.Define("path", "TEXT", "NOT NULL")
	sbMoviePoster.Define("source", "TEXT", "NOT NULL", "DEFAULT 'tmdb'")
	sbMoviePoster.Define("FOREIGN KEY (movie_id)", "REFERENCES movie(id)", "ON DELETE CASCADE")
	sbMoviePoster.Define("FOREIGN KEY (language_id)", "REFERENCES language(id)", "ON DELETE CASCADE")

	builders := []sqlbuilder.Builder{sbMovie, sbMovieTranslation, sbMoviePoster}

	for _, b := range builders {
		query, args := b.Build()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	}

	return nil
}

func downMovieInit(ctx context.Context, tx *sql.Tx) error {
	tables := []string{"movie_poster", "movie_translation", "movie"}

	for _, table := range tables {
		if _, err := tx.ExecContext(ctx, "DROP TABLE IF EXISTS "+table); err != nil {
			return err
		}
	}

	return nil
}
