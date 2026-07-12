package feed_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"github.com/arinsuda/movie-hub/internal/shared"
	tmdbmodule "github.com/arinsuda/movie-hub/internal/tmdb_module"
	"gorm.io/gorm"
)

type Service interface {
	// CreateActivity ให้ module อื่น (review/library/like/achievement) เรียกหลังทำ action สำเร็จ
	// จะเช็ก ActivityPrivacySetting ของ actor ก่อนเสมอ — ถ้าปิดแชร์ activity type นี้ไว้ จะไม่สร้าง event ให้
	// (ทำให้ user ยังใช้งาน watchlist/watched/like ได้ตามปกติ แค่ไม่โผล่ใน feed คนอื่น)
	CreateActivity(actorID uint, activityType ActivityType, payload ActivityPayload) error

	GetFeed(userID uint, pq PaginationQuery) (*FeedListResponse, error)
	GetUserActivities(targetUserID, requesterID uint, pq PaginationQuery) (*FeedListResponse, error)

	UpdateVisibility(activityID, requesterID uint, isVisible bool) error

	GetSettings(userID uint) (*ActivitySettingsResponse, error)
	UpdateSettings(userID uint, req UpdateActivitySettingsRequest) (*ActivitySettingsResponse, error)
}

type service struct {
	repo *repository
}

func NewService(db *gorm.DB) Service {
	return &service{repo: newRepository(db)}
}

func (s *service) CreateActivity(actorID uint, activityType ActivityType, payload ActivityPayload) error {
	enabled, err := s.repo.IsEnabled(actorID, activityType)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}

	event := &ActivityEvent{
		ActorID:       actorID,
		Type:          activityType,
		MediaID:       payload.MediaID,
		MediaType:     payload.MediaType,
		ReviewID:      payload.ReviewID,
		CommentID:     payload.CommentID,
		AchievementID: payload.AchievementID,
		LibraryItemID: payload.LibraryItemID,
		Message:       payload.Message,
		IsVisible:     true,
	}
	return s.repo.Create(event)
}

func (s *service) GetFeed(userID uint, pq PaginationQuery) (*FeedListResponse, error) {
	pq.Normalize()
	rows, total, err := s.repo.FindFeed(userID, pq)
	if err != nil {
		return nil, err
	}
	return toFeedListResponse(rows, pq, total), nil
}

func (s *service) GetUserActivities(targetUserID, requesterID uint, pq PaginationQuery) (*FeedListResponse, error) {
	canView, err := s.repo.canViewUserActivities(requesterID, targetUserID)
	if err != nil {
		return nil, err
	}
	if !canView {
		return nil, ErrForbidden
	}

	pq.Normalize()
	rows, total, err := s.repo.FindUserActivities(targetUserID, pq)
	if err != nil {
		return nil, err
	}
	return toFeedListResponse(rows, pq, total), nil
}

func (s *service) UpdateVisibility(activityID, requesterID uint, isVisible bool) error {
	event, err := s.repo.FindByID(activityID)
	if err != nil {
		return err
	}
	if event.ActorID != requesterID {
		return ErrForbidden
	}
	return s.repo.UpdateVisibility(activityID, isVisible)
}

func (s *service) GetSettings(userID uint) (*ActivitySettingsResponse, error) {
	rows, err := s.repo.FindSettings(userID)
	if err != nil {
		return nil, err
	}
	return toSettingsResponse(rows), nil
}

func (s *service) UpdateSettings(userID uint, req UpdateActivitySettingsRequest) (*ActivitySettingsResponse, error) {
	updates := map[ActivityType]*bool{
		ActivityReviewCreated:       req.ReviewCreated,
		ActivityReviewCommented:     req.ReviewCommented,
		ActivityReviewLiked:         req.ReviewLiked,
		ActivityMediaLiked:          req.MediaLiked,
		ActivityWatchlistAdded:      req.WatchlistAdded,
		ActivityWatchedAdded:        req.WatchedAdded,
		ActivityAchievementUnlocked: req.AchievementUnlocked,
	}

	for activityType, value := range updates {
		if value == nil {
			continue
		}
		if err := s.repo.UpsertSetting(userID, activityType, *value); err != nil {
			return nil, err
		}
	}

	return s.GetSettings(userID)
}

// ── Helpers ──────────────────────────────────────────────────────

func toSettingsResponse(rows []ActivityPrivacySetting) *ActivitySettingsResponse {
	enabled := make(map[ActivityType]bool, len(AllActivityTypes))
	for _, t := range AllActivityTypes {
		enabled[t] = defaultEnabled[t]
	}
	for _, row := range rows {
		enabled[row.ActivityType] = row.Enabled
	}

	return &ActivitySettingsResponse{
		ReviewCreated:       enabled[ActivityReviewCreated],
		ReviewCommented:     enabled[ActivityReviewCommented],
		ReviewLiked:         enabled[ActivityReviewLiked],
		MediaLiked:          enabled[ActivityMediaLiked],
		WatchlistAdded:      enabled[ActivityWatchlistAdded],
		WatchedAdded:        enabled[ActivityWatchedAdded],
		AchievementUnlocked: enabled[ActivityAchievementUnlocked],
	}
}

func toFeedListResponse(rows []feedRow, pq PaginationQuery, total int64) *FeedListResponse {
	items := make([]FeedItemResponse, len(rows))
	for i, row := range rows {
		items[i] = toFeedItemResponse(row)
	}
	return &FeedListResponse{
		Items:      items,
		Pagination: newPaginationMeta(pq.Page, pq.Limit, total),
	}
}

func toFeedItemResponse(row feedRow) FeedItemResponse {
	item := FeedItemResponse{
		ID:   row.ID,
		Type: row.Type,
		Actor: ActorSummary{
			ID:          row.ActorID,
			Username:    row.ActorUsername,
			DisplayName: row.ActorDisplayName,
			AvatarURL:   row.ActorAvatarURL,
		},
		ReviewID:      row.ReviewID,
		CommentID:     row.CommentID,
		AchievementID: row.AchievementID,
		LibraryItemID: row.LibraryItemID,
		Message:       row.Message,
		CreatedAt:     row.CreatedAt,
	}

	if row.MediaID != nil && row.MediaType != nil {
		media := shared.MediaSummary{
			ID:        *row.MediaID,
			MediaType: movie_module.MediaType(*row.MediaType),
		}
		switch movie_module.MediaType(*row.MediaType) {
		case movie_module.MediaMovie:
			if details, err := tmdbmodule.GetMovieByID(*row.MediaID); err == nil && details != nil {
				media.Title = details.Title
				media.PosterURL = details.PosterPath
				media.Genres = details.Genres
				media.VoteAverage = float32(details.VoteAverage)
			}
		case movie_module.MediaSeries:
			if details, err := tmdbmodule.GetSeriesByID(*row.MediaID); err == nil && details != nil {
				media.Title = details.Name
				media.PosterURL = details.PosterPath
				media.Genres = details.Genres
				media.VoteAverage = float32(details.VoteAverage)
			}
		}
		item.Media = &media
	}

	return item
}
