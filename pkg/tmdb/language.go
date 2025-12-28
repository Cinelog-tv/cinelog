package tmdb

import (
	"encoding/json"
	"net/http"
)

type TmdbLanguageResponse struct {
	Iso639_1    string `json:"iso_639_1"`
	EnglishName string `json:"english_name"`
	Name        string `json:"name"`
}

type TmdbLanguageListResponse struct {
	Languages []TmdbLanguageResponse `json:"languages"`
}

func (t *TmdbConfig) GetLanguages() (*TmdbLanguageListResponse, error) {
	url := t.config.baseURL + "/configuration/languages?api_key=" + t.config.apiKey

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var languages TmdbLanguageListResponse

	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, err
	}

	return &languages, nil
}
