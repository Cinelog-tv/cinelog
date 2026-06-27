package tmdb

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	TV_DETAILS_ENDPOINT = "/tv/"
)

const (
	TV_CACHE_NAME = "tv-%d-%s"
)

type TMDBTVDetails struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	OriginalName     string        `json:"original_name"`
	OriginalLanguage string        `json:"original_language"`
	Overview         string        `json:"overview"`
	FirstAirDate     string        `json:"first_air_date"`
	LastAirDate      string        `json:"last_air_date"`
	PosterPath       string        `json:"poster_path"`
	BackdropPath     string        `json:"backdrop_path"`
	VoteAverage      float64       `json:"vote_average"`
	VoteCount        int           `json:"vote_count"`
	Genres           []TMDBGenre   `json:"genres"`
	Popularity       float64       `json:"popularity"`
	Status           string        `json:"status"`
	Type             string        `json:"type"`
	NumberOfSeasons  int           `json:"number_of_seasons"`
	NumberOfEpisodes int           `json:"number_of_episodes"`
	EpisodeRunTime   []int         `json:"episode_run_time"`
	InProduction     bool          `json:"in_production"`
	Networks         []TMDBNetwork `json:"networks"`
	CreatedBy        []TMDBCreator `json:"created_by"`
	Seasons          []TMDBSeason  `json:"seasons"`
}

type TMDBSeason struct {
	ID           int      `json:"id"`
	EpisodeCount int      `json:"episode_count"`
	Name         string   `json:"name"`
	Overview     string   `json:"overview"`
	PosterPath   string   `json:"poster_path"`
	AirDate      TMDBTime `json:"air_date"`
	SeasonNumber int      `json:"season_number"`
}

type TMDBTVItem struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
	OriginalLanguage string   `json:"original_language"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	PosterPath       string   `json:"poster_path"`
	BackdropPath     string   `json:"backdrop_path"`
	VoteAverage      float64  `json:"vote_average"`
	VoteCount        int      `json:"vote_count"`
	GenreIds         []int    `json:"genre_ids"`
	Popularity       float64  `json:"popularity"`
	OriginCountry    []string `json:"origin_country"`
}

type TMDBTVCredits struct {
	ID   int              `json:"id"`
	Cast []TMDBCastMember `json:"cast"`
	Crew []TMDBCrewMember `json:"crew"`
}

type TMDBSeasonDetails struct {
	ID           int           `json:"id"`
	AirDate      TMDBTime      `json:"air_date"`
	Episodes     []TMDBEpisode `json:"episodes"`
	Name         string        `json:"name"`
	Overview     string        `json:"overview"`
	PosterPath   string        `json:"poster_path"`
	SeasonNumber int           `json:"season_number"`
}

type TMDBEpisode struct {
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Overview      string   `json:"overview"`
	AirDate       TMDBTime `json:"air_date"`
	EpisodeNumber int      `json:"episode_number"`
	SeasonNumber  int      `json:"season_number"`
	StillPath     string   `json:"still_path"`
	VoteAverage   float64  `json:"vote_average"`
	VoteCount     int      `json:"vote_count"`
	Runtime       int      `json:"runtime"`
}

type TMDBTVParameters struct {
	TVID     int    `json:"tv_id"`
	Language string `json:"language"`
}

func (client *TMDBClient) GetTVDetail(parameters *TMDBTVParameters) (*TMDBTVDetails, error) {

	if parameters == nil {
		return nil, fmt.Errorf("parameters.MovieID is required")
	}

	var tvID int
	if tvID = parameters.TVID; tvID <= 0 {
		return nil, fmt.Errorf("MovieID must be positive integer")
	}

	queryParams := url.Values{}

	var language string
	if v := strings.TrimSpace(parameters.Language); v != "" {
		language = parameters.Language
	} else {
		language = "en"
	}

	cacheKey := fmt.Sprintf(TV_CACHE_NAME, tvID, language)

	if tv, found := client.config.Cache.Get(cacheKey); found {
		return tv.(*TMDBTVDetails), nil
	}

	queryParams.Add("language", language)

	req, err := client.newRequest("GET", TV_DETAILS_ENDPOINT, &queryParams, nil)

	if err != nil {
		return nil, err
	}

	var tvDetails TMDBTVDetails
	if err := client.do(req, &tvDetails); err != nil {
		return nil, err
	}

	client.config.Cache.Set(cacheKey, &tvDetails, 0)

	return &tvDetails, nil
}
