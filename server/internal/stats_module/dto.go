package statsmodule

import "github.com/arinsuda/movie-hub/internal/library_module"

type StatsResponse struct {
	MediaID        int                      `json:"media_id"`
	MediaType      library_module.MediaType `json:"media_type"`
	LikeCount      int                      `json:"like_count"`
	ViewCount      int                      `json:"view_count"`
	ReviewCount    int                      `json:"review_count"`
	WatchlistCount int                      `json:"watchlist_count"`
}

type IncrementField string

const (
	FieldLike      IncrementField = "like_count"
	FieldView      IncrementField = "view_count"
	FieldReview    IncrementField = "review_count"
	FieldWatchlist IncrementField = "watchlist_count"
)
