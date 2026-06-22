package like_module

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
)

type  MediaLike struct {
	gorm.Model
	UserID    uint                     `gorm:"not null;index;uniqueIndex:idx_user_media_like"`
	MediaID   int                      `gorm:"not null;index;uniqueIndex:idx_user_media_like"`
	MediaType movie_module.MediaType `gorm:"not null;uniqueIndex:idx_user_media_like"`
}
