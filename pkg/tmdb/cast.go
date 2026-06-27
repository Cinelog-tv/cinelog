package tmdb

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	PERSON_DETAILS_ENDPOINT = "/person/"
)

type TMDBCastMember struct {
	ID                 int    `json:"id"`
	Name               string `json:"name"`
	Character          string `json:"character"`
	Order              int    `json:"order"`
	ProfilePath        string `json:"profile_path"`
	Gender             int    `json:"gender"`
	KnownForDepartment string `json:"known_for_department"`
}

type TMDBCrewMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
	Gender      int    `json:"gender"`
}

type TMDBCreator struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type TMDBPersonDetails struct {
	ID                 int      `json:"id"`
	IMDBID             string   `json:"imdb_id"`
	Adult              bool     `json:"adult"`
	AlsoKnownAs        []string `json:"also_known_as"`
	Biography          string   `biography:"biography"`
	Birthday           TMDBTime `json:"birthday"`
	DeathDay           TMDBTime `json:"deathday"`
	Gender             int      `json:"gender"` // TODO: Replacer with struct
	Homepage           string   `json:"homepage"`
	KnownForDepartment string   `json:"known_for_department"`
	Name               string   `json:"name"`
	PlaceOfBirth       string   `json:"place_of_birth"`
	Popularity         int      `json:"popularity"`
	ProfilePath        string   `json:"profile_path"`
}

type TMDBPersonDetailsParameters struct {
	PersonID int    `json:"person_id"`
	Language string `json:"language"`
}

func (client *TMDBClient) GetPersonDetail(parameters *TMDBTVParameters) (*TMDBTVDetails, error) {

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
