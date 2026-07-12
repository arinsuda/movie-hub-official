package feed_module

import (
	"time"

	users "github.com/arinsuda/movie-hub/internal/user_module"
)

type ActivityType string

const (
	ActivityReviewCreated       ActivityType = "review_created"
	ActivityReviewCommented     ActivityType = "review_commented"
	ActivityReviewLiked         ActivityType = "review_liked"
	ActivityMediaLiked          ActivityType = "media_liked"
	ActivityWatchlistAdded      ActivityType = "watchlist_added"
	ActivityWatchedAdded        ActivityType = "watched_added"
	ActivityAchievementUnlocked ActivityType = "achievement_unlocked"
)

// AllActivityTypes ใช้ตอนคืนค่า activity settings ให้ FE ครบทุกประเภท
// แม้ user จะยังไม่เคยตั้งค่าเอง (ไม่มี row ใน activity_privacy_settings) ก็ตาม
var AllActivityTypes = []ActivityType{
	ActivityReviewCreated,
	ActivityReviewCommented,
	ActivityReviewLiked,
	ActivityMediaLiked,
	ActivityWatchlistAdded,
	ActivityWatchedAdded,
	ActivityAchievementUnlocked,
}

// defaultEnabled คือค่า default ตอนยังไม่มี ActivityPrivacySetting ของ user คนนั้นเลย
// like อาจ spam feed ได้ง่าย เลย default ปิดไว้ก่อน ส่วนที่เหลือ default เปิด
var defaultEnabled = map[ActivityType]bool{
	ActivityReviewCreated:       true,
	ActivityReviewCommented:     true,
	ActivityReviewLiked:         false,
	ActivityMediaLiked:          false,
	ActivityWatchlistAdded:      true,
	ActivityWatchedAdded:        true,
	ActivityAchievementUnlocked: true,
}

// ActivityEvent คือหัวใจของ feed — หนึ่ง row ต่อหนึ่งการกระทำที่ user ทำแล้วเลือกจะแชร์
// (สร้างผ่าน Service.CreateActivity เท่านั้น อย่า insert ตรงจาก module อื่น)
type ActivityEvent struct {
	ID uint `gorm:"primarykey;autoIncrement"`

	ActorID uint       `gorm:"not null;index"` // คนที่ทำ action เช่น user 2 รีวิวหนัง
	Actor   users.User `gorm:"foreignKey:ActorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	Type ActivityType `gorm:"type:varchar(40);not null;index"`

	MediaID   *int    `gorm:"index"`
	MediaType *string `gorm:"type:varchar(10)"`

	ReviewID      *uint `gorm:"index"`
	CommentID     *uint `gorm:"index"`
	AchievementID *uint `gorm:"index"`
	LibraryItemID *uint `gorm:"index"`

	Message string `gorm:"type:text"`

	// IsVisible คือสวิตช์ให้เจ้าของ activity ซ่อนรายการนี้จาก feed คนอื่นได้ทีละอัน
	// (ไม่ใช่การลบ record — เก็บไว้เผื่อ audit/analytics)
	IsVisible bool `gorm:"default:true;index"`

	CreatedAt time.Time `gorm:"index"`
}

// ActivityPrivacySetting คือ setting รายประเภทของ user ว่าจะให้ activity type ไหน
// สร้างเป็น ActivityEvent เข้า feed คนอื่นได้บ้าง (คล้าย activity privacy ของ Facebook)
// ไม่มี row ของ (user, type) คู่ไหน = ยังไม่เคยตั้งค่า -> ใช้ defaultEnabled แทน
type ActivityPrivacySetting struct {
	ID uint `gorm:"primarykey;autoIncrement"`

	UserID uint       `gorm:"not null;uniqueIndex:idx_user_activity_type"`
	User   users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	ActivityType ActivityType `gorm:"type:varchar(40);not null;uniqueIndex:idx_user_activity_type"`

	Enabled bool `gorm:"default:true"`

	UpdatedAt time.Time
	CreatedAt time.Time
}
