package tmdb

import (
	"net/url"
	"strings"
)

const (
	GENRE_MOVIE_ENDPOINT = "/genre/movie/list"
	GENRE_TV_ENDPOINT    = "/genre/tv/list"
)

type TMDBGenreResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TMDBGenreParameters struct {
	Language string `json:"language"`
}

func (client *TMDBClient) GetMovieGenres(parameters *TMDBGenreParameters) (*TMDBGenreResponse, error) {
	return client.getGenre(parameters, GENRE_MOVIE_ENDPOINT)
}

func (client *TMDBClient) GetTVGenres(parameters *TMDBGenreParameters) (*TMDBGenreResponse, error) {
	return client.getGenre(parameters, GENRE_TV_ENDPOINT)
}

func (client *TMDBClient) getGenre(parameters *TMDBGenreParameters, endpoint string) (*TMDBGenreResponse, error) {
	queryParams := url.Values{}
	if v := strings.TrimSpace(parameters.Language); v != "" {
		queryParams.Add("language", v)
	}

	req, err := client.newRequest("GET", endpoint, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var tmdbGenreResponse TMDBGenreResponse
	if err := client.do(req, &tmdbGenreResponse); err != nil {
		return nil, err
	}

	return &tmdbGenreResponse, nil
}
