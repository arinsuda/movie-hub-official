package analytics_module

// MediaAnalyticsResponse คือ stats ของ media หนึ่ง ที่รวมข้อมูลจาก library + review
type MediaAnalyticsResponse struct {
	MediaID        int     `json:"media_id"`
	MediaType      string  `json:"media_type"`
	WatchlistCount int     `json:"watchlist_count"`
	FavoriteCount  int     `json:"favorite_count"`
	WatchedCount   int     `json:"watched_count"`
	ReviewCount    int     `json:"review_count"`
	AverageRating  float32 `json:"average_rating"` // in-app rating (0 = ยังไม่มี review)
	TrendingScore  float64 `json:"trending_score"`
	IsTrending     bool    `json:"is_trending"` // true ถ้า score อยู่ใน top trending
}

// TrendingItem ใช้ใน list trending
type TrendingItem struct {
	MediaID        int     `json:"media_id"`
	MediaType      string  `json:"media_type"`
	TrendingScore  float64 `json:"trending_score"`
	WatchlistCount int     `json:"watchlist_count"`
	FavoriteCount  int     `json:"favorite_count"`
	WatchedCount   int     `json:"watched_count"`
}
