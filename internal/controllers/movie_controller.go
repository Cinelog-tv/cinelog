package controllers

import (
	"strconv"

	"github.com/Cinelog-tv/backend/internal/dto"
	"github.com/Cinelog-tv/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type MovieController struct {
	mediaService *services.MediaService
}

func RegisterMovieController(api *echo.Group, mediaService *services.MediaService) *MovieController {
	mc := &MovieController{
		mediaService: mediaService,
	}

	api.GET("/discovers", mc.discoverMoviesController)
	api.GET("/movie/:id", mc.getMovieController)
	api.GET("/movie/:id/cast", mc.getMovieCastController)

	return mc
}

// @Summary      Discover movies
// @Description  Get a list of movies based on filters
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        genre          query     string  false  "Genre ID"
// @Param        minimum_rating query     number  false  "Minimum rating"
// @Param        release_date_lt query    string  false  "Release date less than (YYYY-MM-DD)"
// @Param        release_date_gt query    string  false  "Release date greater than (YYYY-MM-DD)"
// @Param        translation    query     string  false  "Language (e.g., en-US)"
// @Param        page           query     int     false  "Page number"
// @Success      200  {object}  dto.DiscoverResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /discovers [get]
func (mc *MovieController) discoverMoviesController(c echo.Context) error {
	var req dto.DiscoverRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	response, err := mc.mediaService.DiscoverMovies(req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to discover movies"})
	}

	return c.JSON(200, response)
}

// @Summary      Get movie details
// @Description  Get detailed information about a specific movie
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "Movie ID"
// @Param        language query     string  false "Language (e.g., en-US)"
// @Success      200  {object}  dto.MediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id} [get]
func (mc *MovieController) getMovieController(c echo.Context) error {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid movie ID"})
	}

	language := c.QueryParam("language")
	if language == "" {
		language = "en-US"
	}

	response, err := mc.mediaService.GetMovieDetails(movieID, language)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch movie details"})
	}

	return c.JSON(200, response)
}

// @Summary      Get movie cast
// @Description  Get the cast and crew of a specific movie
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "Movie ID"
// @Param        language query     string  false "Language (e.g., en-US)"
// @Success      200  {array}   dto.CastResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /movie/{id}/cast [get]
func (mc *MovieController) getMovieCastController(c echo.Context) error {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid movie ID"})
	}

	language := c.QueryParam("language")
	if language == "" {
		language = "en-US"
	}

	cast, err := mc.mediaService.GetMovieCast(movieID, language)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch movie cast"})
	}

	return c.JSON(200, cast)
}
