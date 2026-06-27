package tmdb

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	MOVIE_DETAILS_ENDPOINT = "/movie/"
)

const (
	MOVIE_CACHE_KEY = "movie-%d-%s"
)

type TMDBMovieDetails struct {
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
	Genres           []TMDBGenre `json:"genres"`
	Adult            bool        `json:"adult"`
	Popularity       float64     `json:"popularity"`
	Budget           int         `json:"budget"`
	Revenue          int         `json:"revenue"`
	Status           string      `json:"status"`
	Tagline          string      `json:"tagline"`
	Video            bool        `json:"video"`
}

type TMDBMovieTranslations struct {
	ID           int                           `json:"id"`
	Translations []TMDBMovieTranslationDetails `json:"translations"`
}

type TMDBMovieTranslationDetails struct {
	ISO639_1    string                   `json:"iso_639_1"`
	ISO3166_1   string                   `json:"iso_3166_1"`
	Name        string                   `json:"name"`
	EnglishName string                   `json:"english_name"`
	Data        TMDBMovieTranslationData `json:"data"`
}

type TMDBMovieTranslationData struct {
	Title    string `json:"title"`
	Overview string `json:"overview"`
	Homepage string `json:"homepage"`
	Tagline  string `json:"tagline"`
}

type TMDBMovieParameters struct {
	MovieID  int    `json:"movie_id"`
	Language string `json:"language"`
}

func (client *TMDBClient) GetMovieDetails(parameters *TMDBMovieParameters) (*TMDBMovieDetails, error) {

	if parameters == nil {
		return nil, fmt.Errorf("parameters.MovieID is required")
	}

	var movieID int
	if movieID = parameters.MovieID; movieID <= 0 {
		return nil, fmt.Errorf("MovieID must be positive integer")
	}

	queryParams := url.Values{}

	var language string
	if v := strings.TrimSpace(parameters.Language); v != "" {
		language = parameters.Language
	} else {
		language = "en"
	}

	cacheKey := fmt.Sprintf(MOVIE_CACHE_KEY, movieID, language)

	if movie, found := client.config.Cache.Get(cacheKey); found {
		return movie.(*TMDBMovieDetails), nil
	}

	queryParams.Add("language", language)

	req, err := client.newRequest("GET", MOVIE_DETAILS_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var movieDetails TMDBMovieDetails
	if err := client.do(req, &movieDetails); err != nil {
		return nil, err
	}

	client.config.Cache.Set(cacheKey, &movieDetails, 0)

	return &movieDetails, nil
}
