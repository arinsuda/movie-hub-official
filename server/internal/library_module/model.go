package library_module

import (
	"gorm.io/gorm"
	"time"
)

type MediaType string
type ListType string

const (
	MediaMovie  MediaType = "movie"
	MediaSeries MediaType = "tv"
)

const (
	ListWatchlist ListType = "watchlist"
	ListFavorite  ListType = "favorite"
	ListWatched   ListType = "watched"
)

type LibraryItem struct {
	gorm.Model
	UserID    uint      `gorm:"not null;index"`
	MediaID   int       `gorm:"not null"`
	MediaType MediaType `gorm:"type:varchar(10);not null"`
	ListType  ListType  `gorm:"type:varchar(20);not null"`
	WatchedAt *time.Time // วันที่ดูจบจริงๆ สำหรับ portfolio timeline
	Tags      *string    `gorm:"type:text"` // JSON array
	Note      *string    `gorm:"type:text"`
}

