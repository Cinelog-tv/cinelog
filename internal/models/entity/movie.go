package entity

import "time"

type MovieEntity struct {
	ID            uint
	TMDBID        int
	IMDBID        string
	ReleaseDate   time.Time
	Runtime       int
	DefaultPoster string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type MovieTranslationEntity struct {
	Movie    MovieEntity
	Language LanguageEntity
	Title    string
	Overview string
	Tagline  string
}

type MoviePosterEntity struct {
	ID       uint
	Movie    MovieEntity
	Language LanguageEntity
	Path     string
	source   string
}
