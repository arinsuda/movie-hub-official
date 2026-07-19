package movie_module

import (
	"context"

	"github.com/arinsuda/movie-hub/internal/shared"
	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
)

type MovieService struct {
	ratingRepo shared.RatingStatsReader
}

func NewMovieService(repo shared.RatingStatsReader) *MovieService {
	return &MovieService{ratingRepo: repo}
}

func (s *MovieService) GetPopular(ctx context.Context, page int) (*tmdb.PaginatedResult[MovieDTO], error) {
	res, err := tmdb.GetPopular(page)
	if err != nil {
		return nil, err
	}
	return s.enrichMovies(ctx, res)
}

func (s *MovieService) GetNowPlaying(ctx context.Context, page int) (*tmdb.PaginatedResult[MovieDTO], error) {
	res, err := tmdb.GetNowPlaying(page)
	if err != nil {
		return nil, err
	}
	return s.enrichMovies(ctx, res)
}

func (s *MovieService) GetTopRated(ctx context.Context, page int) (*tmdb.PaginatedResult[MovieDTO], error) {
	res, err := tmdb.GetTopRated(page)
	if err != nil {
		return nil, err
	}
	return s.enrichMovies(ctx, res)
}

func (s *MovieService) GetUpcoming(ctx context.Context, page int) (*tmdb.PaginatedResult[MovieDTO], error) {
	res, err := tmdb.GetUpcoming(page)
	if err != nil {
		return nil, err
	}
	return s.enrichMovies(ctx, res)
}

func (s *MovieService) SearchMovies(ctx context.Context, query string, page int) (*tmdb.PaginatedResult[MovieDTO], error) {
	res, err := tmdb.SearchMovies(query, page)
	if err != nil {
		return nil, err
	}
	return s.enrichMovies(ctx, res)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id int) (*MovieDetailDTO, error) {
	detail, err := tmdb.GetMovieByID(id)
	if err != nil {
		return nil, err
	}

	stats, err := s.ratingRepo.GetMediaRating(ctx, shared.MediaIdentity{ID: id, Type: shared.MediaTypeMovie})
	if err != nil {
		return nil, err
	}

	return &MovieDetailDTO{
		MovieDetail: *detail,
		Ratings:     NewMediaRatingsResponse(stats, detail.VoteAverage, detail.VoteCount),
	}, nil
}

func (s *MovieService) enrichMovies(ctx context.Context, res *tmdb.PaginatedResult[tmdb.Movie]) (*tmdb.PaginatedResult[MovieDTO], error) {
	if len(res.Results) == 0 {
		return &tmdb.PaginatedResult[MovieDTO]{
			Page:         res.Page,
			Results:      []MovieDTO{},
			TotalPages:   res.TotalPages,
			TotalResults: res.TotalResults,
		}, nil
	}

	ids := make([]shared.MediaIdentity, len(res.Results))
	for i, m := range res.Results {
		ids[i] = shared.MediaIdentity{ID: m.ID, Type: shared.MediaTypeMovie}
	}

	ratings, err := s.ratingRepo.GetBatchMediaRatings(ctx, ids)
	if err != nil {
		return nil, err
	}

	dtos := make([]MovieDTO, len(res.Results))
	for i, m := range res.Results {
		key := shared.MediaKey{ID: m.ID, Type: shared.MediaTypeMovie}
		dtos[i] = MovieDTO{
			Movie:   m,
			Ratings: NewMediaRatingsResponse(ratings[key], m.VoteAverage, m.VoteCount),
		}
	}

	return &tmdb.PaginatedResult[MovieDTO]{
		Page:         res.Page,
		Results:      dtos,
		TotalPages:   res.TotalPages,
		TotalResults: res.TotalResults,
	}, nil
}
