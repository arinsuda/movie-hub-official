package media_stats_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
)

type StatsResponse struct {
	MediaID        int                    `json:"media_id"`
	MediaType      movie_module.MediaType `json:"media_type"`
	LikeCount      int                    `json:"like_count"`
	ViewCount      int                    `json:"view_count"`
	ReviewCount    int                    `json:"review_count"`
	WatchlistCount int                    `json:"watchlist_count"`
	AverageRating  float32                `json:"average_rating"`
	HasRating      bool                   `json:"has_rating"`
	LikedAt        *time.Time             `json:"liked_at,omitempty"`
	WatchlistedAt  *time.Time             `json:"watchlisted_at,omitempty"`
}

type IncrementField string

const (
	FieldView IncrementField = "view_count"
)
