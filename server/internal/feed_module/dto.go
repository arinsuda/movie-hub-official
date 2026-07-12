package feed_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/shared"
)

// ── Activity creation (ใช้ภายใน — module อื่นเรียกผ่าน Service.CreateActivity) ──

// ActivityPayload คือ argument ที่ module อื่น (review, library, like, achievement)
// ส่งเข้ามาตอนสร้าง activity event ใหม่
type ActivityPayload struct {
	MediaID   *int
	MediaType *string

	ReviewID      *uint
	CommentID     *uint
	AchievementID *uint
	LibraryItemID *uint

	Message string
}

// ── Pagination ───────────────────────────────────────────────────

type PaginationQuery struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

func (p *PaginationQuery) Normalize() {
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit < 1 || p.Limit > 100 {
		p.Limit = 20
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func newPaginationMeta(page, limit int, total int64) PaginationMeta {
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}
	return PaginationMeta{Page: page, Limit: limit, Total: total, TotalPages: totalPages}
}

// ── Feed responses ───────────────────────────────────────────────

type ActorSummary struct {
	ID          uint    `json:"id"`
	Username    string  `json:"username"`
	DisplayName *string `json:"display_name"`
	AvatarURL   *string `json:"avatar_url"`
}

type FeedItemResponse struct {
	ID uint `json:"id"`

	Type ActivityType `json:"type"`

	Actor ActorSummary `json:"actor"`

	Media *shared.MediaSummary `json:"media,omitempty"`

	ReviewID      *uint `json:"review_id,omitempty"`
	CommentID     *uint `json:"comment_id,omitempty"`
	AchievementID *uint `json:"achievement_id,omitempty"`
	LibraryItemID *uint `json:"library_item_id,omitempty"`

	Message string `json:"message"`

	CreatedAt time.Time `json:"created_at"`
}

type FeedListResponse struct {
	Items      []FeedItemResponse `json:"items"`
	Pagination PaginationMeta     `json:"pagination"`
}

// ── Visibility ───────────────────────────────────────────────────

type UpdateVisibilityRequest struct {
	IsVisible bool `json:"is_visible"`
}

// ── Settings ─────────────────────────────────────────────────────

type ActivitySettingsResponse struct {
	ReviewCreated       bool `json:"review_created"`
	ReviewCommented     bool `json:"review_commented"`
	ReviewLiked         bool `json:"review_liked"`
	MediaLiked          bool `json:"media_liked"`
	WatchlistAdded      bool `json:"watchlist_added"`
	WatchedAdded        bool `json:"watched_added"`
	AchievementUnlocked bool `json:"achievement_unlocked"`
}

type UpdateActivitySettingsRequest struct {
	ReviewCreated       *bool `json:"review_created"`
	ReviewCommented     *bool `json:"review_commented"`
	ReviewLiked         *bool `json:"review_liked"`
	MediaLiked          *bool `json:"media_liked"`
	WatchlistAdded      *bool `json:"watchlist_added"`
	WatchedAdded        *bool `json:"watched_added"`
	AchievementUnlocked *bool `json:"achievement_unlocked"`
}
