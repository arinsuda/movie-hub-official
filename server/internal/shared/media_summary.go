package shared

import (
	"github.com/arinsuda/movie-hub/internal/movie_module"
	tmdbmodule "github.com/arinsuda/movie-hub/internal/tmdb_module"
)

type MediaSummary struct {
	ID          int                    `json:"id"`
	Title       string                 `json:"title"`
	PosterURL   string                 `json:"poster_url"`
	MediaType   movie_module.MediaType `json:"media_type"`
	Genres      []tmdbmodule.Genre     `json:"genres"`
	VoteAverage float32                `json:"vote_average"`
}
