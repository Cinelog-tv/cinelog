package tmdb

const (
	LANGUAGE_ENDPOINT = "/configuration/languages"
)

type TMDBLanguage struct {
	ISO639_1    string `json:"iso_639_1"`
	EnglishName string `json:"english_name"`
	Name        string `json:"name"`
}

type TMDBLanguageResponse struct {
	Languages []TMDBLanguage `json:"languages"`
}

func (client *TMDBClient) GetLanguages() (*TMDBLanguageResponse, error) {

	req, err := client.newRequest("GET", LANGUAGE_ENDPOINT, nil, nil)
	if err != nil {
		return nil, err
	}

	var languages TMDBLanguageResponse
	if err := client.do(req, &languages); err != nil {
		return nil, err
	}

	return &languages, nil
}
