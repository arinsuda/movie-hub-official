package notification_module

import (
	"time"

	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotifFollowedYou NotificationType = "followed_you"

	NotifMovieNowPlaying NotificationType = "movie_now_playing"

	NotifFollowingReviewed       NotificationType = "following_reviewed"
	NotifFollowingLikedReview    NotificationType = "following_liked_review"
	NotifFollowingAddedWatchlist NotificationType = "following_added_watchlist"
	NotifFollowingAddedWatched   NotificationType = "following_added_watched"
)

type Notification struct {
	gorm.Model
	UserID    uint             `gorm:"not null;index"`
	User      users.User       `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActorID   *uint            `gorm:"index"`
	Type      NotificationType `gorm:"type:varchar(50);not null;index"`
	TargetID  *uint
	TargetRef *string `gorm:"type:varchar(50)"`
	Message   string  `gorm:"type:text;not null"`
	IsRead    bool    `gorm:"default:false;index"`
	ReadAt    *time.Time
}
