package like_module

import (
	"errors"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
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
}

func NewService(db *gorm.DB, achieve achievementsmodule.Service) *Service {
	return &Service{
		repo:       newRepository(db),
		db:         db,
		achieveSvc: achieve,
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
	_, _ = s.achieveSvc.Track(userID, "like_count", int(count))

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
