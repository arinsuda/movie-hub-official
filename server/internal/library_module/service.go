package library_module

import (
	"context"
	"encoding/json"
	"time"

	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/feed_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	shared "github.com/arinsuda/movie-hub/internal/shared"
	tmdbmodule "github.com/arinsuda/movie-hub/internal/tmdb_module"
	users "github.com/arinsuda/movie-hub/internal/user_module"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"gorm.io/gorm"
)

type Service struct {
	repo       *repository
	db         *gorm.DB
	expPort    stats.ExpAdder
	achieveSvc achievementsmodule.Service
	notifSvc   *notification_module.Service
	feedSvc    feed_module.Service
}

func NewService(
	db *gorm.DB,
	exp stats.ExpAdder,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
	feed feed_module.Service,
) *Service {
	return &Service{
		repo:       newRepository(db),
		db:         db,
		expPort:    exp,
		achieveSvc: achieve,
		notifSvc:   notif,
		feedSvc:    feed,
	}
}

func (s *Service) AddItem(userID uint, req AddItemRequest) (*LibraryItemResponse, error) {
	if err := validateAddItemRequest(req); err != nil {
		return nil, err
	}

	existing, err := s.repo.FindMediaStatus(userID, req.MediaID, req.MediaType)
	if err == nil {
		for _, item := range existing {
			if item.ListType == req.ListType {
				return nil, ErrDuplicate
			}
		}
	}

	item := &LibraryItem{
		UserID:    userID,
		MediaID:   req.MediaID,
		MediaType: req.MediaType,
		ListType:  req.ListType,
		Note:      req.Note,
	}
	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		item.WatchedAt = &t
	}

	if err := s.repo.Create(item); err != nil {
		return nil, err
	}

	switch req.ListType {

	case movie_module.ListWatched:

		_ = s.expPort.AddExperience(userID, stats.ExpPerWatched)

		var watchedCount int64
		s.db.Model(&LibraryItem{}).
			Where("user_id = ? AND list_type = 'watched' AND deleted_at IS NULL", userID).
			Count(&watchedCount)
		shared.TrackAndNotify(context.Background(), s.achieveSvc, s.notifSvc, userID, "watched_count", int(watchedCount))

		s.trackLibraryTotal(userID)

		// ── Feed: activity ให้คนที่ follow userID เห็นว่ามาร์คดูแล้ว ──
		if s.feedSvc != nil {
			mediaType := string(req.MediaType)
			_ = s.feedSvc.CreateActivity(userID, feed_module.ActivityWatchedAdded, feed_module.ActivityPayload{
				MediaID:       &req.MediaID,
				MediaType:     &mediaType,
				LibraryItemID: &item.ID,
			})
		}

	case movie_module.ListWatchlist:

		var watchlistCount int64
		s.db.Model(&LibraryItem{}).
			Where("user_id = ? AND list_type = 'watchlist' AND deleted_at IS NULL", userID).
			Count(&watchlistCount)
		shared.TrackAndNotify(context.Background(), s.achieveSvc, s.notifSvc, userID, "watchlist_count", int(watchlistCount))
		s.trackLibraryTotal(userID)

		if s.notifSvc != nil {
			if actor, err := s.getUserSummary(userID); err == nil {
				title, _ := s.fetchTitle(req.MediaID, req.MediaType)
				_ = s.notifSvc.PushFollowingAddedWatchlist(
					context.Background(),
					userID,
					actor.Username,
					uint(req.MediaID),
					title,
				)
			}
		}

		// ── Feed: activity ให้คนที่ follow userID เห็นว่าเพิ่มเข้า watchlist ──
		if s.feedSvc != nil {
			mediaType := string(req.MediaType)
			_ = s.feedSvc.CreateActivity(userID, feed_module.ActivityWatchlistAdded, feed_module.ActivityPayload{
				MediaID:       &req.MediaID,
				MediaType:     &mediaType,
				LibraryItemID: &item.ID,
			})
		}
	}

	return toResponse(item), nil
}

func (s *Service) GetLibrary(userID uint, listType *movie_module.ListType, mediaType *movie_module.MediaType) ([]LibraryItemResponse, error) {
	items, err := s.repo.FindByUser(userID, listType, mediaType)
	if err != nil {
		return nil, err
	}

	responses := make([]LibraryItemResponse, len(items))
	for i, item := range items {
		responses[i] = *toResponse(&item)
	}
	return responses, nil
}

func (s *Service) RemoveItem(itemID, requesterID uint) error {
	item, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return err
	}
	if item.UserID != requesterID {
		return ErrForbidden
	}

	if err := s.repo.Delete(itemID); err != nil {
		return err
	}

	if item.ListType == movie_module.ListWatched {
		_ = s.expPort.AddExperience(requesterID, -stats.ExpPerWatched)
	}
	return nil
}

func (s *Service) UpdateItem(itemID, requesterID uint, req UpdateItemRequest) (*LibraryItemResponse, error) {
	item, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return nil, err
	}
	if item.UserID != requesterID {
		return nil, ErrForbidden
	}

	updates := map[string]any{}
	if req.WatchedAt != nil {
		t, err := time.Parse("2006-01-02", *req.WatchedAt)
		if err != nil {
			return nil, ErrInvalidWatchedAt
		}
		updates["watched_at"] = t
	}
	if req.Tags != nil {
		b, _ := json.Marshal(req.Tags)
		updates["tags"] = string(b)
	}
	if req.Note != nil {
		updates["note"] = req.Note
	}
	if len(updates) == 0 {
		return toResponse(item), nil
	}

	if err := s.repo.Update(itemID, updates); err != nil {
		return nil, err
	}

	updated, err := s.repo.FindOne(itemID, requesterID)
	if err != nil {
		return nil, err
	}
	return toResponse(updated), nil
}

func (s *Service) GetMediaStatus(userID uint, mediaID int, mediaType movie_module.MediaType) (*MediaStatusResponse, error) {
	items, err := s.repo.FindMediaStatus(userID, mediaID, mediaType)
	if err != nil {
		return nil, err
	}

	inLists := make([]MediaItemStatus, len(items))
	for i, item := range items {
		inLists[i] = MediaItemStatus{ListType: item.ListType, ItemID: item.ID}
	}

	return &MediaStatusResponse{MediaID: mediaID, MediaType: mediaType, InLists: inLists}, nil
}

func (s *Service) trackLibraryTotal(userID uint) {
	var total int64
	s.db.Model(&LibraryItem{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&total)
	shared.TrackAndNotify(context.Background(), s.achieveSvc, s.notifSvc, userID, "library_total", int(total))
}

func (s *Service) getUserSummary(userID uint) (*users.User, error) {
	var u users.User
	err := s.db.First(&u, userID).Error
	return &u, err
}

func (s *Service) fetchTitle(mediaID int, mediaType movie_module.MediaType) (string, error) {
	switch mediaType {
	case movie_module.MediaMovie:
		if d, err := tmdbmodule.GetMovieByID(mediaID); err == nil {
			return d.Title, nil
		}
	case movie_module.MediaSeries:
		if d, err := tmdbmodule.GetSeriesByID(mediaID); err == nil {
			return d.Name, nil
		}
	}
	return "", nil
}

func validateAddItemRequest(req AddItemRequest) error {
	if req.MediaType != movie_module.MediaMovie && req.MediaType != movie_module.MediaSeries {
		return ErrInvalidMediaType
	}
	if req.ListType != movie_module.ListWatchlist &&
		req.ListType != movie_module.ListWatched {
		return ErrInvalidListType
	}
	if req.MediaID <= 0 {
		return ErrInvalidMediaID
	}
	return nil
}

func toResponse(item *LibraryItem) *LibraryItemResponse {
	tags := []string{}
	if item.Tags != "" {
		_ = json.Unmarshal([]byte(item.Tags), &tags)
	}

	media := shared.MediaSummary{
		ID:        item.MediaID,
		MediaType: item.MediaType,
	}

	switch item.MediaType {
	case movie_module.MediaMovie:
		if details, err := tmdbmodule.GetMovieByID(item.MediaID); err == nil && details != nil {
			media.Title = details.Title
			media.PosterURL = details.PosterPath
			media.Genres = details.Genres
			media.VoteAverage = float32(details.VoteAverage)
		}
	case movie_module.MediaSeries:
		if details, err := tmdbmodule.GetSeriesByID(item.MediaID); err == nil && details != nil {
			media.Title = details.Name
			media.PosterURL = details.PosterPath
			media.Genres = details.Genres
			media.VoteAverage = float32(details.VoteAverage)
		}
	}

	return &LibraryItemResponse{
		ID:        item.ID,
		Media:     media,
		ListType:  item.ListType,
		Tags:      tags,
		WatchedAt: item.WatchedAt,
		Note:      item.Note,
		CreatedAt: item.CreatedAt,
	}
}
