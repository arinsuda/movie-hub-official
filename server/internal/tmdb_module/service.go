package tmdbmodule

import (
	"fmt"
	"net/url"
)

func pageParams(page int) url.Values {
	p := url.Values{}
	p.Set("page", fmt.Sprintf("%d", page))
	return p
}

func GetPopular(page int) (*PaginatedResult[Movie], error) {
	var result PaginatedResult[Movie]
	if err := get("/movie/popular", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetPopular: %w", err)
	}
	return &result, nil
}

func GetNowPlaying(page int) (*PaginatedResult[Movie], error) {
	var result PaginatedResult[Movie]
	if err := get("/movie/now_playing", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetNowPlaying: %w", err)
	}
	return &result, nil
}

func GetTopRated(page int) (*PaginatedResult[Movie], error) {
	var result PaginatedResult[Movie]
	if err := get("/movie/top_rated", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetTopRated: %w", err)
	}
	return &result, nil
}

func GetUpcoming(page int) (*PaginatedResult[Movie], error) {
	var result PaginatedResult[Movie]
	if err := get("/movie/upcoming", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetUpcoming: %w", err)
	}
	return &result, nil
}

func GetMovieByID(tmdbID int) (*MovieDetail, error) {
	var movie MovieDetail
	path := fmt.Sprintf("/movie/%d", tmdbID)

	err := getWithFallback(path, nil, &movie, func() bool {
		return movie.Overview == ""
	})
	if err != nil {
		return nil, fmt.Errorf("GetMovieByID(%d): %w", tmdbID, err)
	}
	return &movie, nil
}

func GetCredits(tmdbID int) (*Credits, error) {
	var credits Credits
	if err := get(fmt.Sprintf("/movie/%d/credits", tmdbID), nil, &credits); err != nil {
		return nil, fmt.Errorf("GetCredits(%d): %w", tmdbID, err)
	}
	return &credits, nil
}

func GetVideos(tmdbID int) ([]Video, error) {
	var result VideoResult
	path := fmt.Sprintf("/movie/%d/videos", tmdbID)

	err := getWithFallback(path, nil, &result, func() bool {
		return len(result.Results) == 0
	})
	if err != nil {
		return nil, fmt.Errorf("GetVideos(%d): %w", tmdbID, err)
	}
	return result.Results, nil
}

func GetSimilar(tmdbID, page int) (*PaginatedResult[Movie], error) {
	var result PaginatedResult[Movie]
	path := fmt.Sprintf("/movie/%d/similar", tmdbID)

	if err := get(path, pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetSimilar(%d): %w", tmdbID, err)
	}
	return &result, nil
}

func SearchMovies(query string, page int) (*PaginatedResult[Movie], error) {
	params := pageParams(page)
	params.Set("query", query)

	var result PaginatedResult[Movie]
	if err := get("/search/movie", params, &result); err != nil {
		return nil, fmt.Errorf("SearchMovies(%q): %w", query, err)
	}
	return &result, nil
}

func GetGenres() ([]Genre, error) {
	var result struct {
		Genres []Genre `json:"genres"`
	}
	if err := get("/genre/movie/list", nil, &result); err != nil {
		return nil, fmt.Errorf("GetGenres: %w", err)
	}
	return result.Genres, nil
}
