package like_module

import (
	"context"
	"errors"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/feed_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/shared"
	tmdbmodule "github.com/arinsuda/movie-hub/internal/tmdb_module"
	usermodule "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

type Service struct {
	repo       *repository
	db         *gorm.DB
	achieveSvc achievementsmodule.Service
	notifSvc   *notification_module.Service
	feedSvc    feed_module.Service
}

func NewService(db *gorm.DB, achieve achievementsmodule.Service, notif *notification_module.Service, feed feed_module.Service) *Service {
	return &Service{
		repo:       newRepository(db),
		db:         db,
		achieveSvc: achieve,
		notifSvc:   notif,
		feedSvc:    feed,
	}
}

func (s *Service) Like(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	if err := s.repo.Create(userID, mediaID, mediaType); err != nil {
		return err
	}

	var count int64
	s.db.Model(&MediaLike{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count)
	shared.TrackAndNotify(context.Background(), s.achieveSvc, s.notifSvc, userID, "like_count", int(count))

	// ── Feed: activity ให้คนที่ follow userID เห็นว่าไป like media นี้ ──
	// default ปิดไว้ใน ActivityPrivacySetting (media_liked) เพราะ like เกิดถี่ อาจ spam feed
	// ผู้ใช้เปิดเองได้ผ่าน PATCH /me/activity-settings
	if s.feedSvc != nil {
		mt := string(mediaType)
		_ = s.feedSvc.CreateActivity(userID, feed_module.ActivityMediaLiked, feed_module.ActivityPayload{
			MediaID:   &mediaID,
			MediaType: &mt,
		})
	}

	return nil
}

func (s *Service) Unlike(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	return s.repo.Delete(userID, mediaID, mediaType)
}

func (s *Service) GetLikes(ownerID, requesterID uint) ([]LikeResponse, error) {
	userRepo := usermodule.NewPublicRepository(s.db)
	owner, _, _, _, err := userRepo.FindByID(ownerID)
	if err != nil {
		return nil, err
	}

	if owner.IsPrivate && ownerID != requesterID {
		return nil, ErrForbidden
	}

	likes, err := s.repo.FindByUser(ownerID)
	if err != nil {
		return nil, err
	}

	responses := make([]LikeResponse, len(likes))
	for i, like := range likes {
		media := shared.MediaSummary{
			ID:        like.MediaID,
			MediaType: like.MediaType,
		}

		switch like.MediaType {
		case movie_module.MediaMovie:
			if details, err := tmdbmodule.GetMovieByID(like.MediaID); err == nil && details != nil {
				media.Title = details.Title
				media.PosterURL = details.PosterPath
				media.Genres = details.Genres
				media.VoteAverage = float32(details.VoteAverage)
			}
		case movie_module.MediaSeries:
			if details, err := tmdbmodule.GetSeriesByID(like.MediaID); err == nil && details != nil {
				media.Title = details.Name
				media.PosterURL = details.PosterPath
				media.Genres = details.Genres
				media.VoteAverage = float32(details.VoteAverage)
			}
		}

		responses[i] = LikeResponse{
			ID:        like.ID,
			UserID:    like.UserID,
			Media:     media,
			CreatedAt: like.CreatedAt,
		}
	}

	return responses, nil
}
