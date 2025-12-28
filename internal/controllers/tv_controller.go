package controllers

import (
	"strconv"

	"github.com/Cinelog-tv/backend/internal/dto"
	"github.com/Cinelog-tv/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type TVController struct {
	mediaService *services.MediaService
}

func RegisterTVController(api *echo.Group, mediaService *services.MediaService) *TVController {
	tc := &TVController{
		mediaService: mediaService,
	}

	api.GET("/discovers/tv", tc.discoverTVController)
	api.GET("/tv/:id", tc.getTVController)
	api.GET("/tv/:id/cast", tc.getTVCastController)
	api.GET("/tv/:id/season/:season", tc.getSeasonController)

	return tc
}

// @Summary      Discover TV shows
// @Description  Get a list of TV shows based on filters
// @Tags         TV Shows
// @Accept       json
// @Produce      json
// @Param        genre          query     string  false  "Genre ID"
// @Param        minimum_rating query     number  false  "Minimum rating"
// @Param        release_date_lt query    string  false  "First air date less than (YYYY-MM-DD)"
// @Param        release_date_gt query    string  false  "First air date greater than (YYYY-MM-DD)"
// @Param        translation    query     string  false  "Language (e.g., en-US)"
// @Param        page           query     int     false  "Page number"
// @Success      200  {object}  dto.DiscoverResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /discovers/tv [get]
func (tc *TVController) discoverTVController(c echo.Context) error {
	var req dto.DiscoverRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	response, err := tc.mediaService.DiscoverTV(req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to discover TV shows"})
	}

	return c.JSON(200, response)
}

// @Summary      Get TV show details
// @Description  Get detailed information about a specific TV show
// @Tags         TV Shows
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "TV Show ID"
// @Param        language query     string  false "Language (e.g., en-US)"
// @Success      200  {object}  dto.MediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tv/{id} [get]
func (tc *TVController) getTVController(c echo.Context) error {
	tvID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid TV ID"})
	}

	language := c.QueryParam("language")
	if language == "" {
		language = "en-US"
	}

	response, err := tc.mediaService.GetTVDetails(tvID, language)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch TV details"})
	}

	return c.JSON(200, response)
}

// @Summary      Get TV show cast
// @Description  Get the cast and crew of a specific TV show
// @Tags         TV Shows
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "TV Show ID"
// @Param        language query     string  false "Language (e.g., en-US)"
// @Success      200  {array}   dto.CastResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tv/{id}/cast [get]
func (tc *TVController) getTVCastController(c echo.Context) error {
	tvID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid TV ID"})
	}

	language := c.QueryParam("language")
	if language == "" {
		language = "en-US"
	}

	cast, err := tc.mediaService.GetTVCast(tvID, language)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch TV cast"})
	}

	return c.JSON(200, cast)
}

// @Summary      Get season details
// @Description  Get detailed information about a specific season of a TV show
// @Tags         TV Shows
// @Accept       json
// @Produce      json
// @Param        id       path      int     true  "TV Show ID"
// @Param        season   path      int     true  "Season number"
// @Param        language query     string  false "Language (e.g., en-US)"
// @Success      200  {object}  object
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /tv/{id}/season/{season} [get]
func (tc *TVController) getSeasonController(c echo.Context) error {
	tvID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid TV ID"})
	}

	seasonNum, err := strconv.Atoi(c.Param("season"))
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid season number"})
	}

	language := c.QueryParam("language")
	if language == "" {
		language = "en-US"
	}

	seasonDetails, err := tc.mediaService.GetSeasonDetails(tvID, seasonNum, language)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch season details"})
	}

	for i := range seasonDetails.Episodes {
		if seasonDetails.Episodes[i].StillPath != "" {
			seasonDetails.Episodes[i].StillPath = "https://image.tmdb.org/t/p/w500/" + seasonDetails.Episodes[i].StillPath
		}
	}

	if seasonDetails.PosterPath != "" {
		seasonDetails.PosterPath = "https://image.tmdb.org/t/p/w500/" + seasonDetails.PosterPath
	}

	return c.JSON(200, seasonDetails)
}
