package bmol_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
	shared "github.com/arinsuda/movie-hub/internal/shared"
)

type CreateBMOLRequest struct {
	MediaID   int                    `json:"media_id" binding:"required"`
	MediaType movie_module.MediaType `json:"media_type" binding:"required"`
	Rank      int                    `json:"rank" binding:"required,min=1"`
}

type UpdateBMOLRequest struct {
	Rank int `json:"rank" binding:"required,min=1"`
}

type BMOLItemResponse struct {
	ID        uint                   `json:"id"`
	Media     shared.MediaSummary    `json:"media"`
	MediaType movie_module.MediaType `json:"media_type"`
	Rank      int                    `json:"rank"`
	CreatedAt time.Time              `json:"created_at"`
}
