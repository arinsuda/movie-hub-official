package user_stats_module

import (
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

type UserStat struct {
	UserID         uint       `gorm:"column:user_id"`
	User           users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ReviewCount    int        `gorm:"column:review_count"`
	LikeCount      int        `gorm:"column:like_count"`
	WatchlistCount int        `gorm:"column:watchlist_count"`
	WatchedCount   int        `gorm:"column:watched_count"`
	FollowerCount  int        `gorm:"column:follower_count"`
	FollowingCount int        `gorm:"column:following_count"`
}

func (UserStat) TableName() string { return "user_stats" }

type UserStatus struct {
	gorm.Model
	UserID     uint       `gorm:"uniqueIndex;not null"`
	User       users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Level      int        `gorm:"not null;default:1"`
	CurrentExp int        `gorm:"not null;default:0"`
}

func (UserStatus) TableName() string { return "user_statuses" }
