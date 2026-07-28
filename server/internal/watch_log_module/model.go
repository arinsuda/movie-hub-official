package watch_log_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MediaWatchLog struct {
	ID         uint             `gorm:"primarykey;autoIncrement"`
	UserID     uint             `gorm:"not null;index"`
	User       user_module.User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	MediaType  string           `gorm:"type:varchar(10);not null"`
	MediaID    int              `gorm:"not null"`
	WatchedOn  datatypes.Date   `gorm:"type:date;not null"`
	Visibility string           `gorm:"type:varchar(20);not null;default:'public'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

type WatchSummary struct {
	WatchCount     int
	FirstWatchedOn *string // YYYY-MM-DD
	LastWatchedOn  *string // YYYY-MM-DD
	HasWatched     bool
}
