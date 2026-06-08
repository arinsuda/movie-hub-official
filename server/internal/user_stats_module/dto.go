package user_stats_module

type UserStatsResponse struct {
	UserID         uint `json:"user_id"`
	Level          int  `json:"level"`
	CurrentExp     int  `json:"current_exp"`
	ReviewCount    int  `json:"review_count"`
	LikeCount      int  `json:"like_count"`
	WatchlistCount int  `json:"watchlist_count"`
	WatchedCount   int  `json:"watched_count"`
	FollowerCount  int  `json:"follower_count"`
	FollowingCount int  `json:"following_count"`
}
