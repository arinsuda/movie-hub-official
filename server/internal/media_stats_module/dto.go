package media_stats_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
)

// StatsResponse คือ aggregated stats ของ media หนึ่งเรื่อง
//
//   - LikeCount      → COUNT จาก media_likes table
//   - ViewCount      → counter จาก media_stats table
//   - ReviewCount    → COUNT จาก reviews table (soft-delete safe)
//   - WatchlistCount → COUNT จาก library_items WHERE list_type = 'watchlist'
//   - LikedAt        → timestamp ที่ requester กด like media นี้ (ถ้ามี)
//   - WatchlistedAt  → timestamp ที่ requester กด addToWatchlist media นี้ (ถ้ามี)
type StatsResponse struct {
	MediaID        int                    `json:"media_id"`
	MediaType      movie_module.MediaType `json:"media_type"`
	LikeCount      int                    `json:"like_count"`
	ViewCount      int                    `json:"view_count"`
	ReviewCount    int                    `json:"review_count"`
	WatchlistCount int                    `json:"watchlist_count"`
	LikedAt        *time.Time             `json:"liked_at,omitempty"`
	WatchlistedAt  *time.Time             `json:"watchlisted_at,omitempty"`
}

type IncrementField string

const (
	FieldView IncrementField = "view_count"
)
