package tmdb

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type TVDetails struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	OriginalName     string      `json:"original_name"`
	OriginalLanguage string      `json:"original_language"`
	Overview         string      `json:"overview"`
	FirstAirDate     string      `json:"first_air_date"`
	LastAirDate      string      `json:"last_air_date"`
	PosterPath       string      `json:"poster_path"`
	BackdropPath     string      `json:"backdrop_path"`
	VoteAverage      float64     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	Genres           []GenreItem `json:"genres"`
	Popularity       float64     `json:"popularity"`
	Status           string      `json:"status"`
	Type             string      `json:"type"`
	NumberOfSeasons  int         `json:"number_of_seasons"`
	NumberOfEpisodes int         `json:"number_of_episodes"`
	EpisodeRunTime   []int       `json:"episode_run_time"`
	InProduction     bool        `json:"in_production"`
	Networks         []Network   `json:"networks"`
	CreatedBy        []Creator   `json:"created_by"`
	Seasons          []Season    `json:"seasons"`
}

type Season struct {
	ID           int      `json:"id"`
	EpisodeCount int      `json:"episode_count"`
	Name         string   `json:"name"`
	Overview     string   `json:"overview"`
	PosterPath   string   `json:"poster_path"`
	AirDate      TmdbTime `json:"air_date"`
	SeasonNumber int      `json:"season_number"`
}

type Network struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type Creator struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type TVDiscoverResult struct {
	Page         int      `json:"page"`
	Results      []TVItem `json:"results"`
	TotalPages   int      `json:"total_pages"`
	TotalResults int      `json:"total_results"`
}

type TVItem struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	OriginalName     string      `json:"original_name"`
	OriginalLanguage string      `json:"original_language"`
	Overview         string      `json:"overview"`
	FirstAirDate     string      `json:"first_air_date"`
	PosterPath       string      `json:"poster_path"`
	BackdropPath     string      `json:"backdrop_path"`
	VoteAverage      float64     `json:"vote_average"`
	VoteCount        int         `json:"vote_count"`
	GenreIds         []int    `json:"genre_ids"`
	Popularity       float64  `json:"popularity"`
	OriginCountry    []string `json:"origin_country"`
}

type TVCredits struct {
	ID   int          `json:"id"`
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

func (t *TmdbConfig) GetTVDetails(tvID int, language string) (*TVDetails, error) {
	if language == "" {
		language = "en-US"
	}

	if tv, found := t.cache.Get("tv-" + strconv.Itoa(tvID) + language); found {
		return tv.(*TVDetails), nil
	}

	endpoint := fmt.Sprintf("%s/tv/%d", t.config.baseURL, tvID)

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

	var result TVDetails

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t.cache.Set("tv-"+strconv.Itoa(tvID)+language, &result, 0)

	return &result, nil
}

func (t *TmdbConfig) GetTVCredits(tvID int) (*TVCredits, error) {
	endpoint := fmt.Sprintf("%s/tv/%d/credits", t.config.baseURL, tvID)

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

	var result TVCredits

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

type SeasonDetails struct {
	ID           int       `json:"id"`
	AirDate      TmdbTime  `json:"air_date"`
	Episodes     []Episode `json:"episodes"`
	Name         string    `json:"name"`
	Overview     string    `json:"overview"`
	PosterPath   string    `json:"poster_path"`
	SeasonNumber int       `json:"season_number"`
}

type Episode struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Overview       string   `json:"overview"`
	AirDate        TmdbTime `json:"air_date"`
	EpisodeNumber  int      `json:"episode_number"`
	SeasonNumber   int      `json:"season_number"`
	StillPath      string   `json:"still_path"`
	VoteAverage    float64  `json:"vote_average"`
	VoteCount      int      `json:"vote_count"`
	Runtime        int      `json:"runtime"`
}

func (t *TmdbConfig) GetSeasonDetails(tvID int, seasonNumber int, language string) (*SeasonDetails, error) {
	if language == "" {
		language = "en-US"
	}

	if season, found := t.cache.Get("season-" + strconv.Itoa(tvID) + "-" + strconv.Itoa(seasonNumber) + language); found {
		return season.(*SeasonDetails), nil
	}

	endpoint := fmt.Sprintf("%s/tv/%d/season/%d", t.config.baseURL, tvID, seasonNumber)

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

	var result SeasonDetails

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	t.cache.Set("season-"+strconv.Itoa(tvID)+"-"+strconv.Itoa(seasonNumber)+language, &result, 0)

	return &result, nil
}

func (t *TmdbConfig) DiscoverTV(genre string, minimumRating float64, firstAirDateLt string, firstAirDateGt string, translation string, page int) (*TVDiscoverResult, error) {
	endpoint := fmt.Sprintf("%s/discover/tv", t.config.baseURL)

	queryParams := url.Values{}
	queryParams.Add("api_key", t.config.apiKey)

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

	var result TVDiscoverResult

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
