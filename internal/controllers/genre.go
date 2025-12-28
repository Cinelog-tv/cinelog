package controllers

import (
	"strconv"

	"github.com/Cinelog-tv/backend/pkg/tmdb"
	"github.com/labstack/echo/v4"
	"github.com/patrickmn/go-cache"
)

type GenreController struct {
	tmdb *tmdb.TmdbConfig
	c    *cache.Cache
}

type GenreRequest struct {
	Language string `json:"language"`
}

func RegisterGenreController(api *echo.Group, tmdb *tmdb.TmdbConfig, c *cache.Cache) {
	genreController := &GenreController{
		tmdb: tmdb,
		c:    c,
	}

	api.GET("/movies/genre", genreController.getMovieGenre)
	api.GET("/shows/genre", genreController.getTVGenre)

}

func (g *GenreController) getMovieGenre(c echo.Context) error {
	var genreRequest GenreRequest
	language := c.QueryParam("language")
	if language == "" {
		genreRequest = GenreRequest{
			Language: "en-US",
		}
	} else {
		genreRequest = GenreRequest{
			Language: language,
		}
	}

	genres, err := g.tmdb.GetMovieGenre(genreRequest.Language)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch tv genre"})
	}

	for i := range genres.Genres {
		g.c.Set("genre-"+genreRequest.Language+strconv.Itoa(genres.Genres[i].ID), &genres.Genres[i], cache.NoExpiration)
	}

	return c.JSON(200, genres.Genres)
}

func (g *GenreController) getTVGenre(c echo.Context) error {

	var genreRequest GenreRequest
	language := c.QueryParam("language")
	if language == "" {
		genreRequest = GenreRequest{
			Language: "en-US",
		}
	} else {
		genreRequest = GenreRequest{
			Language: language,
		}
	}

	genres, err := g.tmdb.GetShowGenre(genreRequest.Language)

	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch tv genre"})
	}

	for i := range genres.Genres {
		g.c.Set("genre-"+genreRequest.Language+strconv.Itoa(genres.Genres[i].ID), &genres.Genres[i], cache.NoExpiration)
	}

	return c.JSON(200, genres.Genres)
}
