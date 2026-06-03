package user_stats_module

// UserStat map กับ VIEW user_stats
// ไม่ใช่ table จริง — อ่านได้อย่างเดียว ห้าม AutoMigrate
type UserStat struct {
	UserID         uint `gorm:"column:user_id"`
	ReviewCount    int  `gorm:"column:review_count"`
	LikeCount      int  `gorm:"column:like_count"`
	WatchlistCount int  `gorm:"column:watchlist_count"`
	WatchedCount   int  `gorm:"column:watched_count"`
	FollowerCount  int  `gorm:"column:follower_count"`
	FollowingCount int  `gorm:"column:following_count"`
}

// TableName ชี้ไปที่ VIEW แทน table
func (UserStat) TableName() string {
	return "user_stats"
}
