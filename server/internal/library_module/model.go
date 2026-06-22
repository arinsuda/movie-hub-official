package library_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
	"time"
)

type LibraryItem struct {
	gorm.Model
	UserID    uint                   `gorm:"not null;index"`
	MediaID   int                    `gorm:"not null"`
	MediaType movie_module.MediaType `gorm:"type:varchar(10);not null"`
	ListType  movie_module.ListType  `gorm:"type:varchar(20);not null"`
	WatchedAt *time.Time             // วันที่ดูจบจริงๆ สำหรับ portfolio timeline
	Tags      string                 `gorm:"type:text;default:'[]'"`
	Note      *string                `gorm:"type:text"`
}
