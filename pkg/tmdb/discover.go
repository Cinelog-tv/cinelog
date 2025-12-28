package tmdb

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type TmdbMediaType string

const (
	TmdbMovie TmdbMediaType = "movie"
	TmdbTV    TmdbMediaType = "tv"
)

type TmdbDiscoverResult struct {
	Page         int                `json:"page"`
	TotalPages   int                `json:"total_pages"`
	TotalResults int                `json:"total_results"`
	Results      []TmdbDiscoverItem `json:"results"`
}

type TmdbDiscoverItem struct {
	ID               int           `json:"id"`
	Title            string        `json:"title"`
	Name             string        `json:"name"`
	OriginalTitle    string        `json:"original_title"`
	OriginalLanguage string        `json:"original_language"`
	MediaType        TmdbMediaType `json:"media_type"`
	ReleaseDate      TmdbTime      `json:"release_date"`
	FirstAirDate     TmdbTime      `json:"first_air_date"`
	Overview         string        `json:"overview"`
	PosterPath       string        `json:"poster_path"`
	BackdropPath     string        `json:"backdrop_path"`
	VoteAverage      float64       `json:"vote_average"`
	VoteCount        int           `json:"vote_count"`
	GenreIDs         []int         `json:"genre_ids"`
	Adult            bool          `json:"adult"`
	Popularity       float64       `json:"popularity"`
	Video            bool          `json:"video"`
}

func (i TmdbDiscoverItem) GetYear() int {
	switch i.MediaType {
	case TmdbMovie:
		if i.ReleaseDate.IsZero() {
			return 0
		}
		return i.ReleaseDate.Year()
	case TmdbTV:
		if i.FirstAirDate.IsZero() {
			return 0
		}
		return i.FirstAirDate.Year()
	default:
		return 0
	}
}

func (i TmdbDiscoverItem) GetTitle() string {
	switch i.MediaType {
	case TmdbMovie:
		return i.Title
	case TmdbTV:
		return i.Name
	}
	return ""
}

func (t *TmdbConfig) DiscoverMovie(page int, releaseLt, releaseGt, genre, language string) (*TmdbDiscoverResult, error) {
	endpoint := t.config.baseURL + "/discover/movie"

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)

	if page > 0 {
		queryParams.Add("page", string(rune(page)))
	} else {
		queryParams.Add("page", "1")
	}

	if releaseLt != "" && releaseLt != "0000-00-00" {
		queryParams.Add("release_date.lte", releaseLt)
	}

	if releaseGt != "" && releaseGt != "0000-00-00" {
		queryParams.Add("release_date.gte", releaseGt)
	}

	if genre != "" {
		queryParams.Add("with_genres", genre)
	}

	fullURL := endpoint + "?" + queryParams.Encode()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TmdbDiscoverResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for i := range result.Results {
		result.Results[i].MediaType = TmdbMovie
	}

	return &result, nil
}
