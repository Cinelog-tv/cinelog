package tmdb

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

const (
	SEARCH_MULTI_ENDPOINT = "search/multi"
)

type TMDBSearchMultiResult struct {
	ID               string        `json:"id"`
	Title            string        `json:"title"`
	OriginalTitle    string        `json:"original_title"`
	OriginalLanguage string        `json:"original_language"`
	Adult            bool          `json:"adult"`
	BackdropPath     string        `json:"backdrop_path"`
	PosterPath       string        `json:"poster_path"`
	MediaType        TMDBMediaType `json:"media_type"`
	GenreIDs         []int         `json:"genre_ids"`
	Popularity       float32       `json:"popularity"`
	ReleaseDate      TMDBTime      `json:"release_date"`
	Video            bool          `json:"video"`
	VoteAverage      float32       `json:"vote_average"`
	VoteCount        int           `json:"vote_count"`
}

type TMDBSearchMulti struct {
	Page    int                     `json:"page"`
	Results []TMDBSearchMultiResult `json:"results"`
}

type TMDBSearchParameters struct {
	Query        string
	FirstAirDate *TMDBTime
	IncludeAdult bool
	Language     string
	Page         int
	Year         int
}

func (c *TMDBClient) SearchMulti(parameters *TMDBSearchParameters) (*TMDBSearchMulti, error) {
	if parameters == nil {
		return nil, fmt.Errorf("parameters is required")
	}

	var q string
	if q = strings.TrimSpace(parameters.Query); q == "" {
		return nil, fmt.Errorf("query is required")
	}

	queryParams := url.Values{}

	if parameters.Page > 0 {
		page := strconv.Itoa(parameters.Page)
		queryParams.Add("page", page)
	}

	if l := strings.TrimSpace(parameters.Language); l != "" {
		queryParams.Add("language", l)
	}

	if parameters.IncludeAdult {
		queryParams.Add("include_adult", "true")
	}

	queryParams.Add("query", q)

	req, err := c.newRequest("GET", SEARCH_MULTI_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var tmdbSearchMulti TMDBSearchMulti
	if err = c.do(req, &tmdbSearchMulti); err != nil {
		return nil, err
	}

	return &tmdbSearchMulti, nil

}
