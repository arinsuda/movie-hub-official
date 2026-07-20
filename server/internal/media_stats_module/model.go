package media_stats_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

type MediaStat struct {
	gorm.Model
	MediaID   int                    `gorm:"not null;uniqueIndex:idx_media"`
	MediaType movie_module.MediaType `gorm:"not null;uniqueIndex:idx_media"`
	LikeCount int                    `gorm:"default:0"`
	ViewCount int                    `gorm:"default:0"`
}
