package tmdbmodule

import (
	"fmt"
	"net/url"
)

func GetPopular(page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get("/movie/popular", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetNowPlaying(page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get("/movie/now_playing", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetTopRated(page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get("/movie/top_rated", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetUpcoming(page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get("/movie/upcoming", params, &result); err != nil {
		return nil, err
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
		return nil, err
	}
	return &movie, nil
}

func GetCredits(tmdbID int) (*Credits, error) {
	var credits Credits
	if err := get(fmt.Sprintf("/movie/%d/credits", tmdbID), nil, &credits); err != nil {
		return nil, err
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
		return nil, err
	}
	return result.Results, nil
}

func GetSimilar(tmdbID int, page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get(fmt.Sprintf("/movie/%d/similar", tmdbID), params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func SearchMovies(query string, page int) (*PaginatedResult[Movie], error) {
	params := url.Values{}
	params.Set("query", query)
	params.Set("page", fmt.Sprintf("%d", page))

	var result PaginatedResult[Movie]
	if err := get("/search/movie", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetGenres() ([]Genre, error) {
	var result struct {
		Genres []Genre `json:"genres"`
	}
	if err := get("/genre/movie/list", nil, &result); err != nil {
		return nil, err
	}
	return result.Genres, nil
}
