package review_module

import (
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"time"
)

// ── Review ────────────────────────────────────────────────────────

type CreateReviewRequest struct {
	MediaID   int     `json:"media_id"`
	MediaType string  `json:"media_type"` // "movie" | "tv"
	Rating    float32 `json:"rating"`     // 0.0 - 10.0
	Body      string  `json:"body"`
	IsPublic  bool    `json:"is_public"`
	WatchedAt *string `json:"watched_at"` // "2026-05-01"
}

type UpdateReviewRequest struct {
	Rating    *float32 `json:"rating"`
	Body      *string  `json:"body"`
	IsPublic  *bool    `json:"is_public"`
	WatchedAt *string  `json:"watched_at"`
}

type ReviewResponse struct {
	ID           uint                      `json:"id"`
	User         users.UserSummaryResponse `json:"user"`
	MediaID      int                       `json:"media_id"`
	MediaType    string                    `json:"media_type"`
	Rating       float32                   `json:"rating"`
	Body         string                    `json:"body"`
	IsPublic     bool                      `json:"is_public"`
	WatchedAt    *time.Time                `json:"watched_at"`
	LikeCount    int                       `json:"like_count"`
	CommentCount int                       `json:"comment_count"`
	IsLiked      bool                      `json:"is_liked"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
}

// ── In-app Rating (aggregate) ─────────────────────────────────────

// RatingResponse คือ aggregate rating ที่คำนวณจาก reviews ภายใน app เท่านั้น
// ไม่เกี่ยวกับ TMDB
type RatingResponse struct {
	MediaID       int     `json:"media_id"`
	MediaType     string  `json:"media_type"`
	AverageRating float32 `json:"average_rating"` // 0.0 ถ้าไม่มี review
	ReviewCount   int     `json:"review_count"`
	HasRating     bool    `json:"has_rating"` // false = ยังไม่มีใครรีวิว
}

// ── Comment ───────────────────────────────────────────────────────

type CreateCommentRequest struct {
	Body string `json:"body"`
}

type UpdateCommentRequest struct {
	Body string `json:"body"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	ReviewID  uint      `json:"review_id"`
	UserID    uint      `json:"user_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
