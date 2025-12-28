package tmdb

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type GenreResponse struct {
	Genres []GenreItem `json:"genres"`
}

type GenreItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (t *TmdbConfig) GetMovieGenre(translation string) (*GenreResponse, error) {
	if translation == "" {
		translation = "en"
	}

	if genres, found := t.cache.Get("genre-movie-" + translation); found {
		return genres.(*GenreResponse), nil
	}

	endpoint := t.config.baseURL + "/genre/movie/list"

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)
	queryParams.Add("language", translation)

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

	var result GenreResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t.cache.Set("genre-"+translation, &result, 0)

	return &result, nil
}

func (t *TmdbConfig) GetShowGenre(translation string) (*GenreResponse, error) {
	if translation == "" {
		translation = "en"
	}

	if genres, found := t.cache.Get("genre-movie-" + translation); found {
		return genres.(*GenreResponse), nil
	}

	endpoint := t.config.baseURL + "/genre/tv/list"

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)
	queryParams.Add("language", translation)

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

	var result GenreResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t.cache.Set("genre-"+translation, &result, 0)

	return &result, nil
}
