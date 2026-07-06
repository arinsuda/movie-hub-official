package review_module

import (
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/gorm"
	"time"
)

type Review struct {
	gorm.Model
	UserID uint       `gorm:"not null;index"`
	User   users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	MediaID   int    `gorm:"not null;index"`
	MediaType string `gorm:"type:varchar(10);not null"`

	Rating    float32 `gorm:"type:decimal(3,1);not null"`
	Body      string  `gorm:"type:text;not null"`
	IsPublic  bool    `gorm:"default:true"`
	WatchedAt *time.Time

	LikeCount    int `gorm:"default:0"`
	CommentCount int `gorm:"default:0"`
	HelpfulCount int `gorm:"default:0"` // NEW: จำนวนคนโหวตว่ารีวิวนี้มีประโยชน์
}

type ReviewLike struct {
	ID        uint       `gorm:"primarykey;autoIncrement"`
	ReviewID  uint       `gorm:"not null;index"`
	UserID    uint       `gorm:"not null;index"`
	User      users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
}

// ReviewHelpful เก็บว่า user คนไหนโหวตว่า review ไหน "มีประโยชน์"
// แยกตารางจาก ReviewLike เพราะเป็นคนละความหมาย (like = ชอบ, helpful = มีประโยชน์)
// และอนาคตอาจมี business rule ต่างกัน (เช่น เฉพาะคนที่ดูหนังเรื่องนี้แล้วโหวตได้)
type ReviewHelpful struct {
	ID        uint       `gorm:"primarykey;autoIncrement"`
	ReviewID  uint       `gorm:"not null;index:idx_review_helpful_review_user,unique,composite:review_user"`
	UserID    uint       `gorm:"not null;index:idx_review_helpful_review_user,unique,composite:review_user"`
	User      users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time
}

type ReviewComment struct {
	gorm.Model
	ReviewID uint       `gorm:"not null;index"`
	UserID   uint       `gorm:"not null;index"`
	User     users.User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Body     string     `gorm:"type:text;not null"`
}
