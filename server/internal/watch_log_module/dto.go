package watch_log_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/shared"
)

type CreateWatchLogRequest struct {
	WatchedOn  string `json:"watched_on" validate:"required"` // YYYY-MM-DD
	Visibility string `json:"visibility" validate:"required,oneof=public followers private"`
}

type UpdateWatchLogRequest struct {
	WatchedOn  *string `json:"watched_on,omitempty"`
	Visibility *string `json:"visibility,omitempty" validate:"omitempty,oneof=public followers private"`
}

type WatchLogResponse struct {
	ID         uint                `json:"id"`
	Media      shared.MediaSummary `json:"media"`
	WatchedOn  string              `json:"watched_on"`
	Visibility string              `json:"visibility"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

type MyWatchLogsResponse struct {
	Logs    []WatchLogResponse   `json:"logs"`
	Summary WatchSummaryResponse `json:"summary"`
}

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

type WatchLogListResponse struct {
	Items      []WatchLogResponse `json:"items"`
	Pagination PaginationMeta     `json:"pagination"`
}

type WatchSummaryResponse struct {
	WatchCount     int     `json:"watch_count"`
	FirstWatchedOn *string `json:"first_watched_on"`
	LastWatchedOn  *string `json:"last_watched_on"`
	HasWatched     bool    `json:"has_watched"`
}
