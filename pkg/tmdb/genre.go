package tmdb

type TMDBGenreResponse struct {
	Genres []TMDBGenre `json:"genres"`
}

type TMDBGenre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
