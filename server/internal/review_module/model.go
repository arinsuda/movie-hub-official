package review_module

import (
	"gorm.io/gorm"
	"time"
)

type Review struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`
	MediaID   int    `gorm:"not null;index"`
	MediaType string `gorm:"type:varchar(10);not null"`

	Rating    float32 `gorm:"type:decimal(3,1);not null"`
	Body      string  `gorm:"type:text;not null"`
	IsPublic  bool    `gorm:"default:true"`
	WatchedAt *time.Time

	LikeCount    int `gorm:"default:0"`
	CommentCount int `gorm:"default:0"`
}

type ReviewLike struct {
	ID        uint `gorm:"primarykey;autoIncrement"`
	ReviewID  uint `gorm:"not null;index"`
	UserID    uint `gorm:"not null;index"`
	CreatedAt time.Time
	// UNIQUE: (review_id, user_id)
}

type ReviewComment struct {
	gorm.Model
	ReviewID uint   `gorm:"not null;index"`
	UserID   uint   `gorm:"not null;index"`
	Body     string `gorm:"type:text;not null"`
}
