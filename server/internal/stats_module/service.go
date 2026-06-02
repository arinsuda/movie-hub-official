package statsmodule

import (
	"github.com/arinsuda/movie-hub/internal/library_module"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

func (s *Service) GetStats(mediaID int, mediaType library_module.MediaType) (*StatsResponse, error) {
	stat, err := s.repo.GetStats(mediaID, mediaType)
	if err != nil {
		return nil, err
	}
	return toResponse(stat), nil
}

func (s *Service) IncrementStat(mediaID int, mediaType library_module.MediaType, field IncrementField, delta int) (*StatsResponse, error) {
	if err := s.repo.Upsert(mediaID, mediaType); err != nil {
		return nil, err
	}
	if err := s.repo.Increment(mediaID, mediaType, field, delta); err != nil {
		return nil, err
	}
	return s.GetStats(mediaID, mediaType)
}

func (s *Service) IncrementReviewCount(mediaID int, mediaType string, delta int) error {
	mt := library_module.MediaType(mediaType)
	if mt != library_module.MediaMovie && mt != library_module.MediaSeries {
		return nil
	}
	_, err := s.IncrementStat(mediaID, mt, FieldReview, delta)
	return err
}

func toResponse(stat *MediaStat) *StatsResponse {
	return &StatsResponse{
		MediaID:        stat.MediaID,
		MediaType:      stat.MediaType,
		LikeCount:      stat.LikeCount,
		ViewCount:      stat.ViewCount,
		ReviewCount:    stat.ReviewCount,
		WatchlistCount: stat.WatchlistCount,
	}
}
