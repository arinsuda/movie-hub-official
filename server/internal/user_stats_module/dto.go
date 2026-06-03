package user_stats_module

// UserStatsResponse คือ stats ของ user ที่แสดงบนหน้า profile
//
//   - ReviewCount   : จำนวน reviews ที่ user เขียน
//   - LikeCount     : จำนวน media ที่ user กด like
//   - WatchlistCount: จำนวน media ที่ user เก็บใน watchlist
//   - WatchedCount  : จำนวน media ที่ user ดูแล้ว
//   - FollowerCount : จำนวนคนที่ follow user นี้
//   - FollowingCount: จำนวนคนที่ user นี้ follow
type UserStatsResponse struct {
	UserID         uint `json:"user_id"`
	ReviewCount    int  `json:"review_count"`
	LikeCount      int  `json:"like_count"`
	WatchlistCount int  `json:"watchlist_count"`
	WatchedCount   int  `json:"watched_count"`
	FollowerCount  int  `json:"follower_count"`
	FollowingCount int  `json:"following_count"`
}
