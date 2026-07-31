package movie_module

import (
	"github.com/arinsuda/movie-hub/internal/shared"
	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
)

type MediaType string
type ListType string

const (
	MediaMovie  MediaType = "movie"
	MediaSeries MediaType = "tv"
)

const (
	ListWatchlist ListType = "watchlist"
	ListWatched   ListType = "watched"
)

type RatingSourceResponse struct {
	Average   *float64 `json:"average"`
	Count     int      `json:"count"`
	Available bool     `json:"available"`
	Scale     float64  `json:"scale"`
}

type MediaRatingsResponse struct {
	REMOVY RatingSourceResponse `json:"removy"`
	TMDB   RatingSourceResponse `json:"tmdb"`
}

type MovieDTO struct {
	tmdb.Movie
	Ratings MediaRatingsResponse `json:"ratings"`
}

type MovieDetailDTO struct {
	tmdb.MovieDetail
	Ratings MediaRatingsResponse `json:"ratings"`
}

type TVSeriesDTO struct {
	tmdb.TVSeries
	Ratings MediaRatingsResponse `json:"ratings"`
}

func NewMediaRatingsResponse(removyStats shared.RatingStats, tmdbAvg float32, tmdbVotes int) MediaRatingsResponse {
	var removyAvg *float64
	if removyStats.Average != nil {
		removyAvg = removyStats.Average
	}
	removyAvailable := removyAvg != nil && removyStats.Count > 0

	var tmdbAvgPtr *float64
	if tmdbVotes > 0 {
		v := float64(tmdbAvg)
		tmdbAvgPtr = &v
	}
	tmdbAvailable := tmdbAvgPtr != nil && tmdbVotes > 0

	return MediaRatingsResponse{
		REMOVY: RatingSourceResponse{
			Average:   removyAvg,
			Count:     removyStats.Count,
			Available: removyAvailable,
			Scale:     5.0,
		},
		TMDB: RatingSourceResponse{
			Average:   tmdbAvgPtr,
			Count:     tmdbVotes,
			Available: tmdbAvailable,
			Scale:     10.0,
		},
	}
}
