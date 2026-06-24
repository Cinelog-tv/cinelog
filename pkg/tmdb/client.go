package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Cinelog-tv/cinelog/pkg/cache"
)

type TMDBClient struct {
	apiKey string
	config *TMDBClientConfig
}

var base_url = "https://api.themoviedb.org/3"

type TMDBClientConfig struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Cache      cache.Cache
}

func NewTMDBClient(apiKey string, config *TMDBClientConfig) (*TMDBClient, error) {

	if config == nil {
		config = &TMDBClientConfig{
			HTTPClient: &http.Client{},
			Cache:      cache.NewMemoryCache(time.Hour*1, time.Minute*30),
		}
	}

	if config.BaseURL == nil {
		parsed, err := url.Parse(base_url)

		if err != nil {
			return nil, err
		}

		config.BaseURL = parsed
	}

	if config.HTTPClient == nil {
		config.HTTPClient = &http.Client{}
	}

	if config.Cache == nil {
		config.Cache = cache.NewMemoryCache(time.Hour*1, time.Minute*30)
	}

	return &TMDBClient{
		apiKey: apiKey,
		config: config,
	}, nil
}

func (client *TMDBClient) newRequest(method, path string, params *url.Values, body io.Reader) (*http.Request, error) {
	endpoint := path
	if params == nil {
		params = &url.Values{}
	}

	params.Add("api_key", client.apiKey)

	endpoint += "?" + params.Encode()

	parsedUrl, err := client.config.BaseURL.Parse(path)

	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %v", err)
	}

	req, err := http.NewRequest(method, parsedUrl.String(), body)

	if err != nil {
		return nil, fmt.Errorf("cannot build request for %s: %v", parsedUrl.String(), err)
	}

	req.Header.Add("User-Agent", "app/cinelog")

	return req, nil
}

func (client *TMDBClient) do(req *http.Request, v interface{}) error {
	resp, err := client.config.HTTPClient.Do(req)

	if err != nil {
		return fmt.Errorf("error during api call: %v", err)
	}

	defer resp.Body.Close()

	if err := parseStatusCode(resp.StatusCode); err != nil {
		return fmt.Errorf("api respond with bad status code %d: %v", resp.StatusCode, err)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
			return fmt.Errorf("cannot decode response body: %v", err)
		}
	}
	return nil
}

func parseStatusCode(code int) error {
	switch code {
	case http.StatusBadRequest:
		return fmt.Errorf("400 bad request")
	case http.StatusNotFound:
		return fmt.Errorf("404 not found")
	case http.StatusInternalServerError:
		return fmt.Errorf("500 internal server error")
	case http.StatusUnauthorized:
		return fmt.Errorf("401 not authorized")
	case http.StatusForbidden:
		return fmt.Errorf("403 forbiden")
	}
	return nil
}
