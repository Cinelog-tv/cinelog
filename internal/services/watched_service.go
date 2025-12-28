package services

import (
	"time"

	"github.com/Cinelog-tv/backend/internal/dto"
	"github.com/Cinelog-tv/backend/internal/models"
	"gorm.io/gorm"
)

type WatchedService struct {
	db *gorm.DB
}

func NewWatchedService(db *gorm.DB) *WatchedService {
	return &WatchedService{db: db}
}

func (s *WatchedService) MarkMediaAsWatched(userID uint, req dto.WatchedMediaRequest) (*dto.WatchedMediaResponse, error) {
	var watchedMedia models.WatchedMedia
	whereClause := s.db.Where("user_id = ? AND tmdb_id = ? AND media_type = ?", userID, req.TmdbID, req.GetMediaType())

	if req.GetMediaType() == models.MediaTypeEpisode {
		whereClause = whereClause.Where("season_number = ? AND episode_number = ?", req.SeasonNum, req.EpisodeNum)
	}

	result := whereClause.First(&watchedMedia)

	if result.Error == nil {
		watchedMedia.Watched = req.Watched
		if req.Watched && watchedMedia.WatchCount == 0 {
			watchedMedia.WatchCount = 1
		}
		if err := s.db.Save(&watchedMedia).Error; err != nil {
			return nil, err
		}
	} else {
		watchCount := 0
		if req.Watched {
			watchCount = 1
		}
		watchedMedia = models.WatchedMedia{
			UserID:     userID,
			TmdbID:     req.TmdbID,
			MediaType:  req.GetMediaType(),
			SeasonNum:  req.SeasonNum,
			EpisodeNum: req.EpisodeNum,
			Watched:    req.Watched,
			WatchCount: watchCount,
		}
		if err := s.db.Create(&watchedMedia).Error; err != nil {
			return nil, err
		}
	}

	return s.toWatchedMediaResponse(watchedMedia), nil
}

func (s *WatchedService) IncrementWatchCount(userID uint, req dto.WatchedMediaRequest) (*dto.WatchedMediaResponse, error) {
	var watchedMedia models.WatchedMedia
	whereClause := s.db.Where("user_id = ? AND tmdb_id = ? AND media_type = ?", userID, req.TmdbID, req.GetMediaType())

	if req.GetMediaType() == models.MediaTypeEpisode {
		whereClause = whereClause.Where("season_number = ? AND episode_number = ?", req.SeasonNum, req.EpisodeNum)
	}

	result := whereClause.First(&watchedMedia)

	if result.Error == nil {
		watchedMedia.WatchCount++
		watchedMedia.Watched = true
		if err := s.db.Save(&watchedMedia).Error; err != nil {
			return nil, err
		}
	} else {
		watchedMedia = models.WatchedMedia{
			UserID:     userID,
			TmdbID:     req.TmdbID,
			MediaType:  req.GetMediaType(),
			SeasonNum:  req.SeasonNum,
			EpisodeNum: req.EpisodeNum,
			Watched:    true,
			WatchCount: 1,
		}
		if err := s.db.Create(&watchedMedia).Error; err != nil {
			return nil, err
		}
	}

	return s.toWatchedMediaResponse(watchedMedia), nil
}

func (s *WatchedService) GetWatchedMedia(userID uint, tmdbID int, mediaType models.MediaType) (*dto.WatchedMediaResponse, error) {
	var watchedMedia models.WatchedMedia
	if err := s.db.Where("user_id = ? AND tmdb_id = ? AND media_type = ?", userID, tmdbID, mediaType).First(&watchedMedia).Error; err != nil {
		return nil, err
	}

	return s.toWatchedMediaResponse(watchedMedia), nil
}

func (s *WatchedService) GetHistory(userID uint) ([]dto.WatchedMediaResponse, error) {
	var watchedMedia []models.WatchedMedia
	if err := s.db.Where("user_id = ? AND watched = ?", userID, true).Order("updated_at DESC").Find(&watchedMedia).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.WatchedMediaResponse, len(watchedMedia))
	for i, media := range watchedMedia {
		responses[i] = *s.toWatchedMediaResponse(media)
	}

	return responses, nil
}

func (s *WatchedService) AddPhysicalMedia(userID uint, req dto.PhysicalMediaRequest) (*dto.PhysicalMediaResponse, error) {
	physicalMedia := models.PhysicalMedia{
		UserID:    userID,
		TmdbID:    req.TmdbID,
		MediaType: req.GetMediaType(),
		SeasonNum: req.SeasonNum,
		Format:    req.Format,
		Edition:   req.Edition,
		Price:     req.Price,
		Store:     req.Store,
	}

	if err := s.db.Create(&physicalMedia).Error; err != nil {
		return nil, err
	}

	// Si mark_as_watched, marquer comme vu
	if req.MarkAsWatched {
		watchReq := dto.WatchedMediaRequest{
			TmdbID:    req.TmdbID,
			MediaType: req.MediaType,
			SeasonNum: req.SeasonNum,
			Watched:   true,
		}
		s.MarkMediaAsWatched(userID, watchReq)
	}

	return s.toPhysicalMediaResponse(physicalMedia), nil
}

func (s *WatchedService) GetPhysicalMediaByID(userID uint, mediaID uint) (*dto.PhysicalMediaResponse, error) {
	var physicalMedia models.PhysicalMedia
	if err := s.db.Where("id = ? AND user_id = ?", mediaID, userID).First(&physicalMedia).Error; err != nil {
		return nil, err
	}

	return s.toPhysicalMediaResponse(physicalMedia), nil
}

func (s *WatchedService) UpdatePhysicalMedia(userID uint, mediaID uint, req dto.PhysicalMediaRequest) (*dto.PhysicalMediaResponse, error) {
	var physicalMedia models.PhysicalMedia
	if err := s.db.Where("id = ? AND user_id = ?", mediaID, userID).First(&physicalMedia).Error; err != nil {
		return nil, err
	}

	physicalMedia.Format = req.Format
	physicalMedia.Edition = req.Edition
	physicalMedia.Price = req.Price
	physicalMedia.Store = req.Store

	if err := s.db.Save(&physicalMedia).Error; err != nil {
		return nil, err
	}

	// Si mark_as_watched, marquer comme vu
	if req.MarkAsWatched {
		watchReq := dto.WatchedMediaRequest{
			TmdbID:    physicalMedia.TmdbID,
			MediaType: string(physicalMedia.MediaType),
			SeasonNum: physicalMedia.SeasonNum,
			Watched:   true,
		}
		s.MarkMediaAsWatched(userID, watchReq)
	}

	return s.toPhysicalMediaResponse(physicalMedia), nil
}

func (s *WatchedService) DeletePhysicalMedia(userID uint, mediaID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", mediaID, userID).Delete(&models.PhysicalMedia{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *WatchedService) GetCollection(userID uint) ([]dto.PhysicalMediaResponse, error) {
	var physicalMedia []models.PhysicalMedia
	if err := s.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&physicalMedia).Error; err != nil {
		return nil, err
	}

	responses := make([]dto.PhysicalMediaResponse, len(physicalMedia))
	for i, media := range physicalMedia {
		responses[i] = *s.toPhysicalMediaResponse(media)
	}

	return responses, nil
}

func (s *WatchedService) GetAnalytics(userID uint) (*dto.AnalyticsResponse, error) {
	var totalWatched int64
	s.db.Model(&models.WatchedMedia{}).Where("user_id = ? AND watched = ?", userID, true).Count(&totalWatched)

	var totalMovies int64
	s.db.Model(&models.WatchedMedia{}).Where("user_id = ? AND watched = ? AND media_type = ?", userID, true, models.MediaTypeMovie).Count(&totalMovies)

	var totalTVShows int64
	s.db.Model(&models.WatchedMedia{}).Where("user_id = ? AND watched = ? AND media_type = ?", userID, true, models.MediaTypeTV).Count(&totalTVShows)

	var totalEpisodes int64
	s.db.Model(&models.WatchedMedia{}).Where("user_id = ? AND watched = ? AND media_type = ?", userID, true, models.MediaTypeEpisode).Count(&totalEpisodes)

	var totalPhysical int64
	s.db.Model(&models.PhysicalMedia{}).Where("user_id = ?", userID).Count(&totalPhysical)

	// Physical by format
	var physicalMedia []models.PhysicalMedia
	s.db.Where("user_id = ?", userID).Find(&physicalMedia)
	physicalByFormat := make(map[string]int)
	for _, media := range physicalMedia {
		physicalByFormat[media.Format]++
	}

	// Most watched media
	var mostWatched []models.WatchedMedia
	s.db.Where("user_id = ?", userID).Order("watch_count DESC").Limit(10).Find(&mostWatched)
	mostWatchedResponses := make([]dto.WatchedMediaResponse, len(mostWatched))
	for i, media := range mostWatched {
		mostWatchedResponses[i] = *s.toWatchedMediaResponse(media)
	}

	return &dto.AnalyticsResponse{
		TotalWatched:     int(totalWatched),
		TotalMovies:      int(totalMovies),
		TotalTVShows:     int(totalTVShows),
		TotalEpisodes:    int(totalEpisodes),
		TotalPhysical:    int(totalPhysical),
		PhysicalByFormat: physicalByFormat,
		MostWatchedMedia: mostWatchedResponses,
	}, nil
}

func (s *WatchedService) toWatchedMediaResponse(media models.WatchedMedia) *dto.WatchedMediaResponse {
	return &dto.WatchedMediaResponse{
		ID:         media.ID,
		UserID:     media.UserID,
		TmdbID:     media.TmdbID,
		MediaType:  string(media.MediaType),
		SeasonNum:  media.SeasonNum,
		EpisodeNum: media.EpisodeNum,
		Watched:    media.Watched,
		WatchCount: media.WatchCount,
		CreatedAt:  media.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  media.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *WatchedService) toPhysicalMediaResponse(media models.PhysicalMedia) *dto.PhysicalMediaResponse {
	purchaseAt := ""
	if media.PurchaseAt != nil && media.PurchaseAt.Valid {
		purchaseAt = media.PurchaseAt.Time.Format(time.RFC3339)
	}

	return &dto.PhysicalMediaResponse{
		ID:         media.ID,
		UserID:     media.UserID,
		TmdbID:     media.TmdbID,
		MediaType:  string(media.MediaType),
		SeasonNum:  media.SeasonNum,
		Format:     media.Format,
		Edition:    media.Edition,
		Price:      media.Price,
		Store:      media.Store,
		PurchaseAt: purchaseAt,
		CreatedAt:  media.CreatedAt.Format(time.RFC3339),
		UpdatedAt:  media.UpdatedAt.Format(time.RFC3339),
	}
}
