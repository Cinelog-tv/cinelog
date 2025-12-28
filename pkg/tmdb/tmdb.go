package tmdb

import (
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

type tmdbConfig struct {
	apiKey  string
	baseURL string
}

type TmdbConfig struct {
	config     *tmdbConfig
	httpClient *http.Client
	cache      *cache.Cache
}

func NewTmdbClient(apiKey string) *TmdbConfig {
	return &TmdbConfig{
		config: &tmdbConfig{
			apiKey:  apiKey,
			baseURL: "https://api.themoviedb.org/3",
		},
		httpClient: &http.Client{},
		cache:      cache.New(time.Hour*1, time.Minute*30),
	}
}
