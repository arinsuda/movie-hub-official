package watch_log_module

import (
	"context"
	"errors"
	"time"

	"github.com/arinsuda/movie-hub/internal/feed_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	"github.com/arinsuda/movie-hub/internal/shared"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Service interface {
	CreateWatchLog(ctx context.Context, userID uint, mediaType string, mediaID int, req CreateWatchLogRequest) (*WatchLogResponse, error)
	GetMyWatchLogs(ctx context.Context, userID uint, mediaType string, mediaID int) ([]WatchLogResponse, error)
	GetUserWatchHistory(ctx context.Context, userID, requesterID uint, pq PaginationQuery) (*WatchLogListResponse, error)
	UpdateWatchLog(ctx context.Context, watchLogID, userID uint, req UpdateWatchLogRequest) (*WatchLogResponse, error)
	DeleteWatchLog(ctx context.Context, watchLogID, userID uint) error
	GetWatchSummary(ctx context.Context, userID uint, mediaType string, mediaID int, viewerID uint) (*WatchSummaryResponse, error)
}

type service struct {
	repo    Repository
	db      *gorm.DB
	feedSvc feed_module.Service
	libSvc  *library_module.Service
	policy  privacy_policy.UserAccessPolicy
}

func NewService(db *gorm.DB, feedSvc feed_module.Service, libSvc *library_module.Service, policy privacy_policy.UserAccessPolicy) Service {
	return &service{
		repo:    NewRepository(db),
		db:      db,
		feedSvc: feedSvc,
		libSvc:  libSvc,
		policy:  policy,
	}
}

func (s *service) CreateWatchLog(ctx context.Context, userID uint, mediaType string, mediaID int, req CreateWatchLogRequest) (*WatchLogResponse, error) {
	if mediaType != "movie" && mediaType != "tv" {
		return nil, ErrInvalidMediaType
	}

	if req.Visibility != "public" && req.Visibility != "followers" && req.Visibility != "private" {
		return nil, ErrInvalidVisibility
	}

	watchedOn, err := time.Parse("2006-01-02", req.WatchedOn)
	if err != nil {
		return nil, ErrInvalidDate
	}

	// Reject future dates > today + 1 day UTC tolerance
	if watchedOn.After(time.Now().UTC().Add(24 * time.Hour)) {
		return nil, ErrFutureDate
	}

	var watchLog *MediaWatchLog
	var isRewatch bool

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)

		count, err := txRepo.CountByUserAndMedia(userID, mediaType, mediaID)
		if err != nil {
			return err
		}
		isRewatch = count > 0

		watchLog = &MediaWatchLog{
			UserID:     userID,
			MediaType:  mediaType,
			MediaID:    mediaID,
			WatchedOn:  datatypes.Date(watchedOn),
			Visibility: req.Visibility,
		}

		if err := txRepo.Create(watchLog); err != nil {
			return err
		}

		// Ensure derived LibraryItem projection exists for 'watched'
		var existingItem library_module.LibraryItem
		err = tx.Unscoped().Where("user_id = ? AND media_id = ? AND media_type = ? AND list_type = ?",
			userID, mediaID, mediaType, movie_module.ListWatched).First(&existingItem).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new watched library item
			newItem := library_module.LibraryItem{
				UserID:    userID,
				MediaID:   mediaID,
				MediaType: movie_module.MediaType(mediaType),
				ListType:  movie_module.ListWatched,
			}
			if err := tx.Create(&newItem).Error; err != nil {
				return err
			}
		} else if err == nil && existingItem.DeletedAt.Valid {
			// Restore soft-deleted library item
			if err := tx.Unscoped().Model(&existingItem).Update("deleted_at", nil).Error; err != nil {
				return err
			}
		}

		// Feed activity creation
		if s.feedSvc != nil && (req.Visibility == "public" || req.Visibility == "followers") {
			activityType := privacy_policy.ActivityWatchLogCreated
			if isRewatch {
				activityType = privacy_policy.ActivityRewatched
			}
			_ = s.feedSvc.CreateActivity(ctx, userID, activityType, feed_module.ActivityPayload{
				MediaID:    &mediaID,
				MediaType:  &mediaType,
				WatchLogID: &watchLog.ID,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return toWatchLogResponse(watchLog), nil
}

func (s *service) GetMyWatchLogs(ctx context.Context, userID uint, mediaType string, mediaID int) ([]WatchLogResponse, error) {
	logs, err := s.repo.FindByUserAndMedia(userID, mediaType, mediaID)
	if err != nil {
		return nil, err
	}
	responses := make([]WatchLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = *toWatchLogResponse(&log)
	}
	return responses, nil
}

func (s *service) GetUserWatchHistory(ctx context.Context, userID, requesterID uint, pq PaginationQuery) (*WatchLogListResponse, error) {
	canView, err := s.policy.CanViewProfileSection(ctx, requesterID, userID, privacy_policy.SectionLibrary)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrForbidden
	}

	logs, total, err := s.repo.FindByUser(userID, pq)
	if err != nil {
		return nil, err
	}

	var visibleLogs []MediaWatchLog
	for _, log := range logs {
		if log.UserID == requesterID {
			visibleLogs = append(visibleLogs, log)
		} else if log.Visibility == "public" {
			visibleLogs = append(visibleLogs, log)
		} else if log.Visibility == "followers" {
			visibleLogs = append(visibleLogs, log)
		}
	}

	responses := make([]WatchLogResponse, len(visibleLogs))
	for i, log := range visibleLogs {
		responses[i] = *toWatchLogResponse(&log)
	}

	return &WatchLogListResponse{
		Items:      responses,
		Pagination: newPaginationMeta(pq.Page, pq.Limit, total),
	}, nil
}

func (s *service) UpdateWatchLog(ctx context.Context, watchLogID, userID uint, req UpdateWatchLogRequest) (*WatchLogResponse, error) {
	log, err := s.repo.FindByID(watchLogID)
	if err != nil {
		return nil, err
	}
	if log == nil {
		return nil, ErrWatchLogNotFound
	}
	if log.UserID != userID {
		return nil, ErrForbidden
	}

	updates := map[string]any{}
	if req.Visibility != nil {
		if *req.Visibility != "public" && *req.Visibility != "followers" && *req.Visibility != "private" {
			return nil, ErrInvalidVisibility
		}
		updates["visibility"] = *req.Visibility
	}
	if req.WatchedOn != nil {
		watchedOn, err := time.Parse("2006-01-02", *req.WatchedOn)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if watchedOn.After(time.Now().UTC().Add(24 * time.Hour)) {
			return nil, ErrFutureDate
		}
		updates["watched_on"] = datatypes.Date(watchedOn)
	}

	if len(updates) == 0 {
		return toWatchLogResponse(log), nil
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)
		if err := txRepo.Update(watchLogID, updates); err != nil {
			return err
		}
		if req.Visibility != nil && s.feedSvc != nil {
			_ = s.feedSvc.UpdateActivityVisibility(ctx, watchLogID, userID, privacy_policy.ActivityVisibility(*req.Visibility))
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	updatedLog, _ := s.repo.FindByID(watchLogID)
	return toWatchLogResponse(updatedLog), nil
}

func (s *service) DeleteWatchLog(ctx context.Context, watchLogID, userID uint) error {
	log, err := s.repo.FindByID(watchLogID)
	if err != nil {
		return err
	}
	if log == nil {
		return ErrWatchLogNotFound
	}
	if log.UserID != userID {
		return ErrForbidden
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		txRepo := NewRepository(tx)
		if err := txRepo.Delete(watchLogID); err != nil {
			return err
		}

		count, err := txRepo.CountByUserAndMedia(userID, log.MediaType, log.MediaID)
		if err == nil && count == 0 {
			// Delete derived watched library item when final watch log is deleted
			_ = tx.Where("user_id = ? AND media_id = ? AND media_type = ? AND list_type = ?",
				userID, log.MediaID, log.MediaType, movie_module.ListWatched).
				Delete(&library_module.LibraryItem{}).Error
		}

		if s.feedSvc != nil {
			_ = s.feedSvc.DeleteWatchLogActivity(ctx, userID, watchLogID)
		}
		return nil
	})

	return err
}

func (s *service) GetWatchSummary(ctx context.Context, userID uint, mediaType string, mediaID int, viewerID uint) (*WatchSummaryResponse, error) {
	summary, err := s.repo.ComputeSummary(userID, mediaType, mediaID)
	if err != nil {
		return nil, err
	}
	return &WatchSummaryResponse{
		WatchCount:     summary.WatchCount,
		FirstWatchedOn: summary.FirstWatchedOn,
		LastWatchedOn:  summary.LastWatchedOn,
		HasWatched:     summary.HasWatched,
	}, nil
}

func toWatchLogResponse(log *MediaWatchLog) *WatchLogResponse {
	if log == nil {
		return nil
	}
	return &WatchLogResponse{
		ID:         log.ID,
		WatchedOn:  time.Time(log.WatchedOn).Format("2006-01-02"),
		Visibility: log.Visibility,
		CreatedAt:  log.CreatedAt,
		UpdatedAt:  log.UpdatedAt,
		Media: shared.MediaSummary{
			ID:        log.MediaID,
			MediaType: shared.MediaType(log.MediaType),
		},
	}
}
