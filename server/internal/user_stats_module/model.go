package user_stats_module

type UserStat struct {
	UserID         uint `gorm:"column:user_id"`
	ReviewCount    int  `gorm:"column:review_count"`
	LikeCount      int  `gorm:"column:like_count"`
	WatchlistCount int  `gorm:"column:watchlist_count"`
	WatchedCount   int  `gorm:"column:watched_count"`
	FollowerCount  int  `gorm:"column:follower_count"`
	FollowingCount int  `gorm:"column:following_count"`
}

func (UserStat) TableName() string {
	return "user_stats"
}
