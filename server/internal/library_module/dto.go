package library_module

import (
	"time"

	"github.com/arinsuda/movie-hub/internal/movie_module"
	shared "github.com/arinsuda/movie-hub/internal/shared"
)

type AddItemRequest struct {
	MediaID   int                    `json:"media_id"`
	MediaType movie_module.MediaType `json:"media_type"`
	ListType  movie_module.ListType  `json:"list_type"`
	WatchedAt *string                `json:"watched_at"`
	Tags      []string               `json:"tags"`
	Note      *string                `json:"note"`
}

type UpdateItemRequest struct {
	WatchedAt *string  `json:"watched_at"`
	Tags      []string `json:"tags"`
	Note      *string  `json:"note"`
}

type LibraryItemResponse struct {
	ID        uint                  `json:"id"`
	Media     shared.MediaSummary   `json:"media"`
	ListType  movie_module.ListType `json:"list_type"`
	WatchedAt *time.Time            `json:"watched_at"`
	Tags      []string              `json:"tags"`
	Note      *string               `json:"note"`
	CreatedAt time.Time             `json:"created_at"`
}

type MediaStatusResponse struct {
	MediaID   int                    `json:"media_id"`
	MediaType movie_module.MediaType `json:"media_type"`
	InLists   []MediaItemStatus      `json:"in_lists"`
}

type MediaItemStatus struct {
	ListType movie_module.ListType `json:"list_type"`
	ItemID   uint                  `json:"item_id"`
}
