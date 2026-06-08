package user_stats_module

import "gorm.io/gorm"

// UserStat is a read-only projection from the user_stats DB view.
// It is never written to directly; counts come from the source tables.
type UserStat struct {
	UserID         uint `gorm:"column:user_id"`
	ReviewCount    int  `gorm:"column:review_count"`
	LikeCount      int  `gorm:"column:like_count"`
	WatchlistCount int  `gorm:"column:watchlist_count"`
	WatchedCount   int  `gorm:"column:watched_count"`
	FollowerCount  int  `gorm:"column:follower_count"`
	FollowingCount int  `gorm:"column:following_count"`
}

func (UserStat) TableName() string { return "user_stats" }

// UserStatus owns the level/exp progression for a user.
type UserStatus struct {
	gorm.Model
	UserID     uint `gorm:"uniqueIndex;not null"`
	Level      int  `gorm:"not null;default:1"`
	CurrentExp int  `gorm:"not null;default:0"`
}

func (UserStatus) TableName() string { return "user_statuses" }
