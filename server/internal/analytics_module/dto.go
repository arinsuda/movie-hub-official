package analytics_module

type MediaAnalyticsResponse struct {
	MediaID        int     `json:"media_id"`
	MediaType      string  `json:"media_type"`
	WatchlistCount int     `json:"watchlist_count"`
	FavoriteCount  int     `json:"favorite_count"`
	WatchedCount   int     `json:"watched_count"`
	ReviewCount    int     `json:"review_count"`
	AverageRating  float32 `json:"average_rating"`
	TrendingScore  float64 `json:"trending_score"`
	IsTrending     bool    `json:"is_trending"`
}

type TrendingItem struct {
	MediaID        int     `json:"media_id"`
	MediaType      string  `json:"media_type"`
	TrendingScore  float64 `json:"trending_score"`
	WatchlistCount int     `json:"watchlist_count"`
	FavoriteCount  int     `json:"favorite_count"`
	WatchedCount   int     `json:"watched_count"`
}
