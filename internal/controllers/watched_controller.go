package controllers

import (
	"strconv"

	"github.com/Cinelog-tv/backend/internal/dto"
	"github.com/Cinelog-tv/backend/internal/services"
	"github.com/labstack/echo/v4"
)

type WatchedController struct {
	watchedService *services.WatchedService
}

func RegisterWatchedController(api *echo.Group, watchedService *services.WatchedService) *WatchedController {
	wc := &WatchedController{
		watchedService: watchedService,
	}

	api.POST("/media/watched", wc.markMediaAsWatchedController)
	api.POST("/media/increment-watch", wc.incrementWatchCountController)

	api.POST("/media/physical", wc.addPhysicalMediaController)
	api.GET("/media/physical/:id", wc.getPhysicalMediaByIdController)
	api.PUT("/media/physical/:id", wc.updatePhysicalMediaController)
	api.DELETE("/media/physical/:id", wc.deletePhysicalMediaController)

	api.GET("/history", wc.getHistoryController)
	api.GET("/collection", wc.getCollectionController)
	api.GET("/analytics", wc.getAnalyticsController)

	return wc
}

// @Summary      Mark media as watched
// @Description  Mark a movie, TV show, or episode as watched or unwatched
// @Tags         Watched Media
// @Accept       json
// @Produce      json
// @Param        request body dto.WatchedMediaRequest true "Watched media request"
// @Success      200  {object}  dto.WatchedMediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/watched [post]
func (wc *WatchedController) markMediaAsWatchedController(c echo.Context) error {
	var req dto.WatchedMediaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if req.MediaType == "" {
		return c.JSON(400, map[string]string{"error": "media_type is required"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	response, err := wc.watchedService.MarkMediaAsWatched(userID, req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to mark media as watched"})
	}

	return c.JSON(200, response)
}

// @Summary      Increment watch count
// @Description  Increment the watch count for a media item
// @Tags         Watched Media
// @Accept       json
// @Produce      json
// @Param        request body dto.WatchedMediaRequest true "Watched media request"
// @Success      200  {object}  dto.WatchedMediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/increment-watch [post]
func (wc *WatchedController) incrementWatchCountController(c echo.Context) error {
	var req dto.WatchedMediaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if req.MediaType == "" {
		return c.JSON(400, map[string]string{"error": "media_type is required"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	response, err := wc.watchedService.IncrementWatchCount(userID, req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to increment watch count"})
	}

	return c.JSON(200, response)
}

// @Summary      Add physical media
// @Description  Add a physical media item to the collection (Blu-ray, DVD, etc.)
// @Tags         Physical Collection
// @Accept       json
// @Produce      json
// @Param        request body dto.PhysicalMediaRequest true "Physical media request"
// @Success      201  {object}  dto.PhysicalMediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/physical [post]
func (wc *WatchedController) addPhysicalMediaController(c echo.Context) error {
	var req dto.PhysicalMediaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	if req.MediaType == "" {
		return c.JSON(400, map[string]string{"error": "media_type is required"})
	}
	if req.Format == "" {
		return c.JSON(400, map[string]string{"error": "format is required"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	response, err := wc.watchedService.AddPhysicalMedia(userID, req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to add physical media"})
	}

	return c.JSON(201, response)
}

// @Summary      Get physical media by ID
// @Description  Get details of a specific physical media item
// @Tags         Physical Collection
// @Accept       json
// @Produce      json
// @Param        id path int true "Physical media ID"
// @Success      200  {object}  dto.PhysicalMediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /media/physical/{id} [get]
func (wc *WatchedController) getPhysicalMediaByIdController(c echo.Context) error {
	mediaID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid media ID"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	response, err := wc.watchedService.GetPhysicalMediaByID(userID, uint(mediaID))
	if err != nil {
		return c.JSON(404, map[string]string{"error": "Physical media not found"})
	}

	return c.JSON(200, response)
}

// @Summary      Update physical media
// @Description  Update details of a physical media item
// @Tags         Physical Collection
// @Accept       json
// @Produce      json
// @Param        id path int true "Physical media ID"
// @Param        request body dto.PhysicalMediaRequest true "Physical media request"
// @Success      200  {object}  dto.PhysicalMediaResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/physical/{id} [put]
func (wc *WatchedController) updatePhysicalMediaController(c echo.Context) error {
	mediaID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid media ID"})
	}

	var req dto.PhysicalMediaRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	response, err := wc.watchedService.UpdatePhysicalMedia(userID, uint(mediaID), req)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to update physical media"})
	}

	return c.JSON(200, response)
}

// @Summary      Delete physical media
// @Description  Remove a physical media item from the collection
// @Tags         Physical Collection
// @Accept       json
// @Produce      json
// @Param        id path int true "Physical media ID"
// @Success      204
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /media/physical/{id} [delete]
func (wc *WatchedController) deletePhysicalMediaController(c echo.Context) error {
	mediaID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid media ID"})
	}

	// TODO: Get userID from context (JWT)
	userID := uint(1)

	if err := wc.watchedService.DeletePhysicalMedia(userID, uint(mediaID)); err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to delete physical media"})
	}

	return c.JSON(204, nil)
}

// @Summary      Get watch history
// @Description  Get the user's watch history
// @Tags         Watched Media
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.WatchedMediaResponse
// @Failure      500  {object}  map[string]string
// @Router       /history [get]
func (wc *WatchedController) getHistoryController(c echo.Context) error {
	// TODO: Get userID from context (JWT)
	userID := uint(1)

	history, err := wc.watchedService.GetHistory(userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch history"})
	}

	return c.JSON(200, history)
}

// @Summary      Get physical collection
// @Description  Get the user's physical media collection
// @Tags         Physical Collection
// @Accept       json
// @Produce      json
// @Success      200  {array}   dto.PhysicalMediaResponse
// @Failure      500  {object}  map[string]string
// @Router       /collection [get]
func (wc *WatchedController) getCollectionController(c echo.Context) error {
	// TODO: Get userID from context (JWT)
	userID := uint(1)

	collection, err := wc.watchedService.GetCollection(userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch collection"})
	}

	return c.JSON(200, collection)
}

// @Summary      Get analytics
// @Description  Get user's watch statistics and analytics
// @Tags         Analytics
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.AnalyticsResponse
// @Failure      500  {object}  map[string]string
// @Router       /analytics [get]
func (wc *WatchedController) getAnalyticsController(c echo.Context) error {
	// TODO: Get userID from context (JWT)
	userID := uint(1)

	analytics, err := wc.watchedService.GetAnalytics(userID)
	if err != nil {
		return c.JSON(500, map[string]string{"error": "Failed to fetch analytics"})
	}

	return c.JSON(200, analytics)
}
