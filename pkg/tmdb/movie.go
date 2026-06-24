package tmdb

import (
	"net/url"
	"strconv"
)

const (
	MOVIE_DETAILS_ENDPOINT = "/movie/"
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

func (client *TMDBClient) GetMovieDetails(movieID int, language string) (*TMDBMovieDetails, error) {
	if language == "" {
		language = "en-US"
	}

	if movie, found := client.config.Cache.Get("movie-" + strconv.Itoa(movieID) + language); found {
		return movie.(*TMDBMovieDetails), nil
	}

	queryParams := url.Values{}
	queryParams.Add("language", language)

	req, err := client.newRequest("GET", MOVIE_DETAILS_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var movieDetails TMDBMovieDetails
	if err := client.do(req, &movieDetails); err != nil {
		return nil, err
	}

	client.config.Cache.Set("movie-"+strconv.Itoa(movieID)+language, &movieDetails, 0)

	return &movieDetails, nil
}
