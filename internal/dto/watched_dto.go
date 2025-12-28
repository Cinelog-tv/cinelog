package dto

import "github.com/Cinelog-tv/backend/internal/models"

type WatchedMediaRequest struct {
	TmdbID     int    `json:"tmdb_id"`
	MediaType  string `json:"media_type"` // "movie", "tv", "episode"
	SeasonNum  int    `json:"season_number"`
	EpisodeNum int    `json:"episode_number"`
	Watched    bool   `json:"watched"`
}

func (r *WatchedMediaRequest) GetMediaType() models.MediaType {
	return models.MediaType(r.MediaType)
}

type WatchedMediaResponse struct {
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	TmdbID     int    `json:"tmdb_id"`
	MediaType  string `json:"media_type"`
	SeasonNum  int    `json:"season_number,omitempty"`
	EpisodeNum int    `json:"episode_number,omitempty"`
	Watched    bool   `json:"watched"`
	WatchCount int    `json:"watch_count"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

type PhysicalMediaRequest struct {
	TmdbID        int     `json:"tmdb_id"`
	MediaType     string  `json:"media_type"` // "movie", "tv"
	SeasonNum     int     `json:"season_number"`
	Format        string  `json:"format"`
	Edition       string  `json:"edition"`
	Price         float64 `json:"price"`
	Store         string  `json:"store"`
	MarkAsWatched bool    `json:"mark_as_watched"` // Mark as watched when adding
}

func (r *PhysicalMediaRequest) GetMediaType() models.MediaType {
	return models.MediaType(r.MediaType)
}

type PhysicalMediaResponse struct {
	ID         uint    `json:"id"`
	UserID     uint    `json:"user_id"`
	TmdbID     int     `json:"tmdb_id"`
	MediaType  string  `json:"media_type"`
	SeasonNum  int     `json:"season_number,omitempty"`
	Format     string  `json:"format"`
	Edition    string  `json:"edition"`
	Price      float64 `json:"price"`
	Store      string  `json:"store"`
	PurchaseAt string  `json:"purchase_at,omitempty"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

type AnalyticsResponse struct {
	TotalWatched     int                    `json:"total_watched"`
	TotalMovies      int                    `json:"total_movies"`
	TotalTVShows     int                    `json:"total_tv_shows"`
	TotalEpisodes    int                    `json:"total_episodes"`
	TotalPhysical    int                    `json:"total_physical"`
	PhysicalByFormat map[string]int         `json:"physical_by_format"`
	MostWatchedMedia []WatchedMediaResponse `json:"most_watched_media"`
}
