package media_stats_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

type Service struct {
	repo *repository
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: newRepository(db)}
}

// GetStats รวม stats จากหลาย source
// requesterID = 0 → anonymous (IsLiked = false เสมอ)
func (s *Service) GetStats(mediaID int, mediaType movie_module.MediaType, requesterID uint) (*StatsResponse, error) {
	viewCount, err := s.repo.GetViewCount(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	likeCount, err := s.repo.CountLikes(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	reviewCount, err := s.repo.CountReviews(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	watchlistCount, err := s.repo.CountWatchlist(mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	var likedAt *time.Time
	if requesterID > 0 {
		likedAt, err = s.repo.GetLikedAt(requesterID, mediaID, mediaType)
		if err != nil {
			return nil, err
		}
	}

	var watchlistedAt *time.Time
	if requesterID > 0 {
		watchlistedAt, err = s.repo.GetWatchlistedAt(requesterID, mediaID, mediaType)
		if err != nil {
			return nil, err
		}
	}

	return &StatsResponse{
		MediaID:        mediaID,
		MediaType:      mediaType,
		LikeCount:      likeCount,
		ViewCount:      viewCount,
		ReviewCount:    reviewCount,
		WatchlistCount: watchlistCount,
		LikedAt:        likedAt,
		WatchlistedAt:  watchlistedAt, // 💡 แตกข้อมูลส่งต่อไปให้เรียบร้อย
	}, nil
}

// RecordView บันทึกการเข้าชม
func (s *Service) RecordView(mediaID int, mediaType movie_module.MediaType) error {
	if err := s.repo.UpsertStat(mediaID, mediaType); err != nil {
		return err
	}
	return s.repo.IncrementView(mediaID, mediaType)
}

// LikeMedia เพิ่ม like
func (s *Service) LikeMedia(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	return s.repo.CreateLike(userID, mediaID, mediaType)
}

// UnlikeMedia ลบ like
func (s *Service) UnlikeMedia(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	return s.repo.DeleteLike(userID, mediaID, mediaType)
}
