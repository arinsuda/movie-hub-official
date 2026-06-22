package like_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/shared"
)

type LikeResponse struct {
	ID        uint                `json:"id"`
	UserID    uint                `json:"user_id"`
	Media     shared.MediaSummary `json:"media"`
	CreatedAt time.Time           `json:"created_at"`
}
