package statsmodule

import (
	"github.com/arinsuda/movie-hub/internal/library_module"
	"gorm.io/gorm"
)

type MediaStat struct {
	gorm.Model
	MediaID        int                      `gorm:"not null;uniqueIndex:idx_media"`
	MediaType      library_module.MediaType `gorm:"not null;uniqueIndex:idx_media"`
	LikeCount      int                      `gorm:"default:0"`
	ViewCount      int                      `gorm:"default:0"`
	ReviewCount    int                      `gorm:"default:0"`
	WatchlistCount int                      `gorm:"default:0"`
}
