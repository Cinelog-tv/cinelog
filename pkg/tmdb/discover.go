package tmdb

import (
	"net/url"
	"strconv"
)

type TMDBMediaType string

const (
	DISCOVER_TV_ENDPOINT    = "/discover/tv"
	DISCOVER_MOVIE_ENDPOINT = "/discover/movie"
)

const (
	TMDB_Movie TMDBMediaType = "movie"
	TMDB_TV    TMDBMediaType = "tv"
)

type TmdbDiscoverResult struct {
	Page         int                `json:"page"`
	TotalPages   int                `json:"total_pages"`
	TotalResults int                `json:"total_results"`
	Results      []TMDBDiscoverItem `json:"results"`
}

type TMDBDiscoverItem struct {
	ID               int           `json:"id"`
	Title            string        `json:"title"`
	Name             string        `json:"name"`
	OriginalTitle    string        `json:"original_title"`
	OriginalLanguage string        `json:"original_language"`
	MediaType        TMDBMediaType `json:"media_type"`
	ReleaseDate      TMDBTime      `json:"release_date"`
	FirstAirDate     TMDBTime      `json:"first_air_date"`
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

func (i TMDBDiscoverItem) GetYear() int {
	switch i.MediaType {
	case TMDB_Movie:
		if i.ReleaseDate.IsZero() {
			return 0
		}
		return i.ReleaseDate.Year()
	case TMDB_TV:
		if i.FirstAirDate.IsZero() {
			return 0
		}
		return i.FirstAirDate.Year()
	default:
		return 0
	}
}

func (i TMDBDiscoverItem) GetTitle() string {
	switch i.MediaType {
	case TMDB_Movie:
		return i.Title
	case TMDB_TV:
		return i.Name
	}
	return ""
}

type TMDBTVDiscoverResult struct {
	Page         int          `json:"page"`
	Results      []TMDBTVItem `json:"results"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
}

func (client *TMDBClient) GetDiscoverTV(genre string, minimumRating float64, firstAirDateLt string, firstAirDateGt string, translation string, page int) (*TMDBTVDiscoverResult, error) {
	queryParams := url.Values{}

	if genre != "" {
		queryParams.Add("with_genres", genre)
	}

	if minimumRating > 0 {
		queryParams.Add("vote_average.gte", strconv.FormatFloat(minimumRating, 'f', 1, 64))
	}

	if firstAirDateLt != "" {
		queryParams.Add("first_air_date.lte", firstAirDateLt)
	}

	if firstAirDateGt != "" {
		queryParams.Add("first_air_date.gte", firstAirDateGt)
	}

	if translation != "" {
		queryParams.Add("language", translation)
	} else {
		queryParams.Add("language", "en-US")
	}

	if page > 0 {
		queryParams.Add("page", strconv.Itoa(page))
	} else {
		queryParams.Add("page", "1")
	}

	queryParams.Add("sort_by", "popularity.desc")

	req, err := client.newRequest("GET", DISCOVER_TV_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var tmdbTvDiscoverResult TMDBTVDiscoverResult
	if err := client.do(req, &tmdbTvDiscoverResult); err != nil {
		return nil, err
	}

	return &tmdbTvDiscoverResult, nil

}

func (client *TMDBClient) DiscoverMovie(page int, releaseLt, releaseGt, genre, language string) (*TmdbDiscoverResult, error) {

	queryParams := url.Values{}

	if page > 0 {
		p := strconv.Itoa(page)

		queryParams.Add("page", p)
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

	if language != "" {
		queryParams.Add("language", language)
	}

	req, err := client.newRequest("GET", DISCOVER_TV_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var tmdbMovieDiscoverResult TmdbDiscoverResult
	if err := client.do(req, &tmdbMovieDiscoverResult); err != nil {
		return nil, err
	}

	return &tmdbMovieDiscoverResult, nil
}
