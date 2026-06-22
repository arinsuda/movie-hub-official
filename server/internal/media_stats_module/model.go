package media_stats_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

// MediaStat เก็บเฉพาะ counter ที่ไม่มี source of truth table อื่น
//
//   - LikeCount    → sync กับ media_likes table (insert/delete)
//   - ViewCount    → counter ล้วนๆ ไม่มี table เก็บ raw events
//
// review_count และ watchlist_count ถูกย้ายไป query จาก
// reviews และ library_items table โดยตรงแทน (ดู StatsResponse)
type MediaStat struct {
	gorm.Model
	MediaID   int                    `gorm:"not null;uniqueIndex:idx_media"`
	MediaType movie_module.MediaType `gorm:"not null;uniqueIndex:idx_media"`
	LikeCount int                    `gorm:"default:0"`
	ViewCount int                    `gorm:"default:0"`
}