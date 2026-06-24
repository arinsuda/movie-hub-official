package notification_module

import (
	"time"

	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
)

// NotificationType กำหนดประเภทของ notification ทั้งหมดในระบบ
type NotificationType string

const (
	// Social notifications
	NotifFollowedYou NotificationType = "followed_you" // มีคน follow เรา

	// Watchlist notifications
	NotifMovieNowPlaying NotificationType = "movie_now_playing" // หนังใน watchlist เข้าฉายแล้ว

	// Following activity notifications
	NotifFollowingReviewed       NotificationType = "following_reviewed"        // คนที่เราติดตาม review หนัง
	NotifFollowingLikedReview    NotificationType = "following_liked_review"    // คนที่เราติดตาม like review
	NotifFollowingAddedWatchlist NotificationType = "following_added_watchlist" // คนที่เราติดตามเพิ่มหนังใน watchlist
	NotifFollowingAddedWatched   NotificationType = "following_added_watched"
)

// Notification เก็บ notification แต่ละรายการของ user
//
// ActorID   = user ที่เป็นต้นเหตุของ notification (อาจเป็น nil สำหรับ system notification เช่น movie release)
// TargetID  = id ของ entity ที่เกี่ยวข้อง เช่น movie_id, review_id
// TargetRef = ชนิดของ entity ("movie", "review", "user") ใช้คู่กับ TargetID เพื่อสร้าง deep-link
type Notification struct {
	gorm.Model
	UserID    uint             `gorm:"not null;index"` // เจ้าของ notification
	User      users.User       `gorm:"foreignKey:UserID;contraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ActorID   *uint            `gorm:"index"` // ผู้กระทำ (nil = system)
	Type      NotificationType `gorm:"type:varchar(50);not null;index"`
	TargetID  *uint            // id ของ entity ที่เกี่ยวข้อง
	TargetRef *string          `gorm:"type:varchar(50)"` // "movie" | "review" | "user"
	Message   string           `gorm:"type:text;not null"`
	IsRead    bool             `gorm:"default:false;index"`
	ReadAt    *time.Time
}
