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

	// Fan-out — แจ้ง follower ว่าคนที่ตามอยู่มี activity
	NotifFollowingReviewed       NotificationType = "following_reviewed"
	NotifFollowingLikedReview    NotificationType = "following_liked_review"
	NotifFollowingAddedWatchlist NotificationType = "following_added_watchlist"
	NotifFollowingAddedWatched   NotificationType = "following_added_watched"
	NotifFollowingMarkedHelpful  NotificationType = "following_marked_helpful" // NEW
	NotifFollowingCommented      NotificationType = "following_commented"      // NEW

	// Direct — แจ้งเจ้าของ content ตรงๆ
	NotifReviewLiked         NotificationType = "review_liked"          // NEW
	NotifReviewCommented     NotificationType = "review_commented"      // NEW
	NotifReviewMarkedHelpful NotificationType = "review_marked_helpful" // NEW

	NotifAchievementUnlocked NotificationType = "achievement_unlocked" // NEW

	NotifEmailVerified   NotificationType = "email_verified"   // NEW
	NotifPasswordChanged NotificationType = "password_changed" // NEW
)

type NotificationCategory string

const (
	CategorySystem      NotificationCategory = "system"
	CategorySocial      NotificationCategory = "social"
	CategoryMedia       NotificationCategory = "media"
	CategoryAchievement NotificationCategory = "achievement"
)

var categoryByType = map[NotificationType]NotificationCategory{
	NotifFollowedYou:             CategorySocial,
	NotifFollowingReviewed:       CategorySocial,
	NotifFollowingLikedReview:    CategorySocial,
	NotifFollowingAddedWatchlist: CategorySocial,
	NotifFollowingAddedWatched:   CategorySocial,
	NotifFollowingMarkedHelpful:  CategorySocial,
	NotifFollowingCommented:      CategorySocial,
	NotifReviewLiked:             CategorySocial,
	NotifReviewCommented:         CategorySocial,
	NotifReviewMarkedHelpful:     CategorySocial,

	NotifMovieNowPlaying: CategoryMedia,

	NotifAchievementUnlocked: CategoryAchievement,

	NotifEmailVerified:   CategorySystem,
	NotifPasswordChanged: CategorySystem,
}

type Notification struct {
	gorm.Model
	UserID    uint                 `gorm:"not null;index"`
	User      users.User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActorID   *uint                `gorm:"index"`
	Type      NotificationType     `gorm:"type:varchar(50);not null;index"`
	Category  NotificationCategory `gorm:"type:varchar(20);not null;index"` // NEW
	TargetID  *uint
	TargetRef *string `gorm:"type:varchar(50)"`
	Message   string  `gorm:"type:text;not null"`
	IsRead    bool    `gorm:"default:false;index"`
	ReadAt    *time.Time
}
