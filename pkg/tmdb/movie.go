package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type MovieDetails struct {
	ID               int         `json:"id"`
	Title            string      `json:"title"`
	OriginalTitle    string      `json:"original_title"`
	OriginalLanguage string      `json:"original_language"`
	Overview         string      `json:"overview"`
	ReleaseDate      string      `json:"release_date"`
	PosterPath       string      `json:"poster_path"`
	BackdropPath     string      `json:"backdrop_path"`
	VoteAverage      float64     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	Runtime          int         `json:"runtime"`
	Genres           []GenreItem `json:"genres"`
	Adult            bool        `json:"adult"`
	Popularity       float64     `json:"popularity"`
	Budget           int         `json:"budget"`
	Revenue          int         `json:"revenue"`
	Status           string      `json:"status"`
	Tagline          string      `json:"tagline"`
	Video            bool        `json:"video"`
}

type MovieTranslations struct {
	ID           int                       `json:"id"`
	Translations []MovieTranslationDetails `json:"translations"`
}

type MovieTranslationDetails struct {
	ISO639_1    string               `json:"iso_639_1"`
	ISO3166_1   string               `json:"iso_3166_1"`
	Name        string               `json:"name"`
	EnglishName string               `json:"english_name"`
	Data        MovieTranslationData `json:"data"`
}

type MovieTranslationData struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Homepage string `json:"homepage"`
	Tagline  string `json:"tagline"`
}

func (t *TmdbConfig) GetMovieDetails(movieID int, language string) (*MovieDetails, error) {
	if language == "" {
		language = "en-US"
	}

	if movie, found := t.cache.Get("movie-" + strconv.Itoa(movieID) + language); found {
		return movie.(*MovieDetails), nil
	}

	endpoint := fmt.Sprintf("%s/movie/%d", t.config.baseURL, movieID)

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)
	queryParams.Add("language", language)

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

	var result MovieDetails

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t.cache.Set("movie-"+strconv.Itoa(movieID)+language, &result, 0)

	return &result, nil
}

func (t *TmdbConfig) GetMovieTranslations(movieID int) (*MovieTranslations, error) {
	endpoint := fmt.Sprintf("%s/movie/%d/translations", t.config.baseURL, movieID)

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)

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

	var result MovieTranslations

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type MovieCredits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

type CastMember struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Character          string `json:"character"`
	Order              int    `json:"order"`
	ProfilePath        string `json:"profile_path"`
	Gender             int    `json:"gender"`
	KnownForDepartment string `json:"known_for_department"`
}

type CrewMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
	Gender      int    `json:"gender"`
}

func (t *TmdbConfig) GetMovieCredits(movieID int) (*MovieCredits, error) {
	endpoint := fmt.Sprintf("%s/movie/%d/credits", t.config.baseURL, movieID)

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)

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

	var result MovieCredits

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
