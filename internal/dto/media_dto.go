package dto

import "github.com/Cinelog-tv/backend/pkg/tmdb"

type MediaType string

const (
	Movie   MediaType = "movie"
	TV      MediaType = "tv"
	Episode MediaType = "episode"
)

type DiscoverRequest struct {
	Genre         string  `json:"genre"`
	MinimumRating float64 `json:"minimum_rating"`
	ReleaseDateLt string  `json:"release_date_lt"`
	ReleaseDateGt string  `json:"release_date_gt"`
	Translation   string  `json:"translation"`
	Page          int     `json:"page"`
}

type MediaResponse struct {
	ID               int              `json:"id"`
	Title            string           `json:"title"`
	Year             int              `json:"year"`
	Type             MediaType        `json:"type"`
	Genre            tmdb.GenreItem   `json:"genre"`
	Rating           float64          `json:"rating"`
	Poster           string           `json:"poster"`
	Backdrop         string           `json:"backdrop,omitempty"`
	Overview         string           `json:"overview,omitempty"`
	ReleaseDate      string           `json:"release_date,omitempty"`
	Runtime          int              `json:"runtime,omitempty"`
	Genres           []tmdb.GenreItem `json:"genres,omitempty"`
	Seasons          []SeasonResponse `json:"seasons,omitempty"`
	NumberOfSeasons  int              `json:"number_of_seasons,omitempty"`
	NumberOfEpisodes int              `json:"number_of_episodes,omitempty"`
	Cast             []CastResponse   `json:"cast,omitempty"`
}

type SeasonResponse struct {
	ID           int    `json:"id"`
	SeasonNumber int    `json:"season_number"`
	EpisodeCount int    `json:"episode_count"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	AirDate      string `json:"air_date"`
}

type EpisodeResponse struct {
	ID            int     `json:"id"`
	EpisodeNumber int     `json:"episode_number"`
	SeasonNumber  int     `json:"season_number"`
	Name          string  `json:"name"`
	Overview      string  `json:"overview"`
	AirDate       string  `json:"air_date"`
	Runtime       int     `json:"runtime"`
	StillPath     string  `json:"still_path"`
	VoteAverage   float64 `json:"vote_average"`
}

type CastResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
	Order       int    `json:"order"`
}

type DiscoverResponse struct {
	Page         int             `json:"page"`
	TotalPages   int             `json:"total_pages"`
	TotalResults int             `json:"total_results"`
	Results      []MediaResponse `json:"results"`
}
