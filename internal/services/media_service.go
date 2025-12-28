package services

import (
	"fmt"
	"time"

	"github.com/Cinelog-tv/backend/internal/dto"
	"github.com/Cinelog-tv/backend/pkg/tmdb"
	"github.com/patrickmn/go-cache"
)

type MediaService struct {
	tmdb  *tmdb.TmdbConfig
	cache *cache.Cache
}

func NewMediaService(tmdbClient *tmdb.TmdbConfig, c *cache.Cache) *MediaService {
	return &MediaService{
		tmdb:  tmdbClient,
		cache: c,
	}
}

func (s *MediaService) DiscoverMovies(req dto.DiscoverRequest) (*dto.DiscoverResponse, error) {
	discoverResult, err := s.tmdb.DiscoverMovie(req.Page, req.ReleaseDateLt, req.ReleaseDateGt, req.Genre, req.Translation)
	if err != nil {
		return nil, err
	}

	mediaResponses := make([]dto.MediaResponse, 0)
	for _, movie := range discoverResult.Results {
		mediaResponses = append(mediaResponses, s.mapDiscoverItemToMediaResponse(movie))
	}

	return &dto.DiscoverResponse{
		Page:         discoverResult.Page,
		TotalPages:   discoverResult.TotalPages,
		TotalResults: discoverResult.TotalResults,
		Results:      mediaResponses,
	}, nil
}

func (s *MediaService) DiscoverTV(req dto.DiscoverRequest) (*dto.DiscoverResponse, error) {
	discoverResult, err := s.tmdb.DiscoverTV(req.Genre, req.MinimumRating, req.ReleaseDateLt, req.ReleaseDateGt, req.Translation, req.Page)
	if err != nil {
		return nil, err
	}

	mediaResponses := make([]dto.MediaResponse, 0)
	for _, tv := range discoverResult.Results {
		mediaResponses = append(mediaResponses, s.mapTVItemToMediaResponse(tv))
	}

	return &dto.DiscoverResponse{
		Page:         discoverResult.Page,
		TotalPages:   discoverResult.TotalPages,
		TotalResults: discoverResult.TotalResults,
		Results:      mediaResponses,
	}, nil
}

func (s *MediaService) GetMovieDetails(movieID int, language string) (*dto.MediaResponse, error) {
	movieDetails, err := s.tmdb.GetMovieDetails(movieID, language)
	if err != nil {
		return nil, err
	}

	response := s.mapMovieDetailsToMediaResponse(*movieDetails)
	return &response, nil
}

func (s *MediaService) GetTVDetails(tvID int, language string) (*dto.MediaResponse, error) {
	tvDetails, err := s.tmdb.GetTVDetails(tvID, language)
	if err != nil {
		return nil, err
	}

	response := s.mapTVDetailsToMediaResponse(*tvDetails)
	return &response, nil
}

func (s *MediaService) GetMovieCast(movieID int, language string) ([]dto.CastResponse, error) {
	credits, err := s.tmdb.GetMovieCredits(movieID)
	if err != nil {
		return nil, err
	}

	castResponses := make([]dto.CastResponse, 0)
	for _, cast := range credits.Cast {
		profilePath := ""
		if cast.ProfilePath != "" {
			profilePath = "https://image.tmdb.org/t/p/w500/" + cast.ProfilePath
		}

		castResponses = append(castResponses, dto.CastResponse{
			ID:          cast.ID,
			Name:        cast.Name,
			Character:   cast.Character,
			ProfilePath: profilePath,
			Order:       cast.Order,
		})
	}

	return castResponses, nil
}

func (s *MediaService) GetTVCast(tvID int, language string) ([]dto.CastResponse, error) {
	credits, err := s.tmdb.GetTVCredits(tvID)
	if err != nil {
		return nil, err
	}

	castResponses := make([]dto.CastResponse, 0)
	for _, cast := range credits.Cast {
		profilePath := ""
		if cast.ProfilePath != "" {
			profilePath = "https://image.tmdb.org/t/p/w500/" + cast.ProfilePath
		}

		castResponses = append(castResponses, dto.CastResponse{
			ID:          cast.ID,
			Name:        cast.Name,
			Character:   cast.Character,
			ProfilePath: profilePath,
			Order:       cast.Order,
		})
	}

	return castResponses, nil
}

func (s *MediaService) GetSeasonDetails(tvID int, seasonNumber int, language string) (*tmdb.SeasonDetails, error) {
	return s.tmdb.GetSeasonDetails(tvID, seasonNumber, language)
}

func (s *MediaService) mapDiscoverItemToMediaResponse(item tmdb.TmdbDiscoverItem) dto.MediaResponse {
	year := item.GetYear()

	var genre tmdb.GenreItem
	if len(item.GenreIDs) > 0 {
		genreName, found := s.cache.Get("genre-" + fmt.Sprintf("%d", item.GenreIDs[0]))
		if found {
			genre = tmdb.GenreItem{ID: item.GenreIDs[0], Name: genreName.(string)}
		} else {
			genre = tmdb.GenreItem{ID: item.GenreIDs[0], Name: "Unknown"}
		}
	} else {
		genre = tmdb.GenreItem{ID: 0, Name: "Unknown"}
	}

	poster := ""
	if item.PosterPath != "" {
		poster = "https://image.tmdb.org/t/p/w500/" + item.PosterPath
	}

	return dto.MediaResponse{
		ID:     item.ID,
		Title:  item.GetTitle(),
		Year:   year,
		Type:   dto.Movie,
		Genre:  genre,
		Rating: item.VoteAverage,
		Poster: poster,
	}
}

func (s *MediaService) mapTVItemToMediaResponse(item tmdb.TVItem) dto.MediaResponse {
	year := 0
	if item.FirstAirDate != "" && len(item.FirstAirDate) >= 4 {
		fmt.Sscanf(item.FirstAirDate[:4], "%d", &year)
	}

	var genre tmdb.GenreItem
	if len(item.GenreIds) > 0 {
		genreName, found := s.cache.Get("genre-" + fmt.Sprintf("%d", item.GenreIds[0]))
		if found {
			genre = tmdb.GenreItem{ID: item.GenreIds[0], Name: genreName.(string)}
		} else {
			genre = tmdb.GenreItem{ID: item.GenreIds[0], Name: "Unknown"}
		}
	} else {
		genre = tmdb.GenreItem{ID: 0, Name: "Unknown"}
	}

	poster := ""
	if item.PosterPath != "" {
		poster = "https://image.tmdb.org/t/p/w500/" + item.PosterPath
	}

	return dto.MediaResponse{
		ID:     item.ID,
		Title:  item.Name,
		Year:   year,
		Type:   dto.TV,
		Genre:  genre,
		Rating: item.VoteAverage,
		Poster: poster,
	}
}

func (s *MediaService) mapMovieDetailsToMediaResponse(movie tmdb.MovieDetails) dto.MediaResponse {
	year := 0
	if len(movie.ReleaseDate) >= 4 {
		fmt.Sscanf(movie.ReleaseDate[:4], "%d", &year)
	}

	var mainGenre tmdb.GenreItem
	if len(movie.Genres) > 0 {
		mainGenre = movie.Genres[0]
	} else {
		mainGenre = tmdb.GenreItem{ID: 0, Name: "Unknown"}
	}

	poster := ""
	if movie.PosterPath != "" {
		poster = "https://image.tmdb.org/t/p/w500/" + movie.PosterPath
	}

	backdrop := ""
	if movie.BackdropPath != "" {
		backdrop = "https://image.tmdb.org/t/p/original/" + movie.BackdropPath
	}

	return dto.MediaResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Year:        year,
		Type:        dto.Movie,
		Genre:       mainGenre,
		Rating:      movie.VoteAverage,
		Poster:      poster,
		Backdrop:    backdrop,
		Overview:    movie.Overview,
		ReleaseDate: movie.ReleaseDate,
		Runtime:     movie.Runtime,
		Genres:      movie.Genres,
	}
}

func (s *MediaService) mapTVDetailsToMediaResponse(tv tmdb.TVDetails) dto.MediaResponse {
	year := 0
	if len(tv.FirstAirDate) >= 4 {
		fmt.Sscanf(tv.FirstAirDate[:4], "%d", &year)
	}

	var mainGenre tmdb.GenreItem
	if len(tv.Genres) > 0 {
		mainGenre = tv.Genres[0]
	} else {
		mainGenre = tmdb.GenreItem{ID: 0, Name: "Unknown"}
	}

	poster := ""
	if tv.PosterPath != "" {
		poster = "https://image.tmdb.org/t/p/w500/" + tv.PosterPath
	}

	backdrop := ""
	if tv.BackdropPath != "" {
		backdrop = "https://image.tmdb.org/t/p/original/" + tv.BackdropPath
	}

	seasons := make([]dto.SeasonResponse, 0)
	for _, season := range tv.Seasons {
		seasonPoster := ""
		if season.PosterPath != "" {
			seasonPoster = "https://image.tmdb.org/t/p/w500/" + season.PosterPath
		}

		airDate := ""
		if !season.AirDate.IsZero() {
			airDate = season.AirDate.ToTime().Format(time.RFC3339)
		}

		seasons = append(seasons, dto.SeasonResponse{
			ID:           season.ID,
			SeasonNumber: season.SeasonNumber,
			EpisodeCount: season.EpisodeCount,
			Name:         season.Name,
			Overview:     season.Overview,
			PosterPath:   seasonPoster,
			AirDate:      airDate,
		})
	}

	return dto.MediaResponse{
		ID:               tv.ID,
		Title:            tv.Name,
		Year:             year,
		Type:             dto.TV,
		Genre:            mainGenre,
		Rating:           tv.VoteAverage,
		Poster:           poster,
		Backdrop:         backdrop,
		Overview:         tv.Overview,
		ReleaseDate:      tv.FirstAirDate,
		Genres:           tv.Genres,
		Seasons:          seasons,
		NumberOfSeasons:  tv.NumberOfSeasons,
		NumberOfEpisodes: tv.NumberOfEpisodes,
	}
}
