package tmdbmodule

import (
	"fmt"
	"net/url"
	"strings"
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

func DiscoverUpcomingByYear(year, page int) (*PaginatedResult[Movie], error) {
	params := pageParams(page)

	params.Set("primary_release_year",
		fmt.Sprintf("%d", year))

	params.Set("sort_by", "popularity.desc")

	params.Set("include_adult", "false")
	params.Set("include_video", "false")

	var result PaginatedResult[Movie]

	err := get("/discover/movie", params, &result)

	return &result, err
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

func GetPopularSeries(page int) (*PaginatedResult[TVSeries], error) {
	var result PaginatedResult[TVSeries]
	if err := get("/tv/popular", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetPopularSeries: %w", err)
	}
	return &result, nil
}

func GetSeriesByID(tmdbID int) (*TVSeriesDetail, error) {
	var series TVSeriesDetail
	path := fmt.Sprintf("/tv/%d", tmdbID)

	err := getWithFallback(path, nil, &series, func() bool {
		return series.Overview == ""
	})
	if err != nil {
		return nil, fmt.Errorf("GetSeriesByID(%d): %w", tmdbID, err)
	}
	return &series, nil
}

func GetSeriesCredits(tmdbID int) (*Credits, error) {
	var credits Credits
	if err := get(fmt.Sprintf("/tv/%d/credits", tmdbID), nil, &credits); err != nil {
		return nil, fmt.Errorf("GetSeriesCredits(%d): %w", tmdbID, err)
	}
	return &credits, nil
}

func GetSeriesVideos(tmdbID int) ([]Video, error) {
	var result VideoResult
	path := fmt.Sprintf("/tv/%d/videos", tmdbID)

	err := getWithFallback(path, nil, &result, func() bool {
		return len(result.Results) == 0
	})
	if err != nil {
		return nil, fmt.Errorf("GetSeriesVideos(%d): %w", tmdbID, err)
	}
	return result.Results, nil
}

func GetSimilarSeries(tmdbID, page int) (*PaginatedResult[TVSeries], error) {
	var result PaginatedResult[TVSeries]
	path := fmt.Sprintf("/tv/%d/similar", tmdbID)

	if err := get(path, pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetSimilarSeries(%d): %w", tmdbID, err)
	}
	return &result, nil
}

func GetNowAiringSeries(page int) (*PaginatedResult[TVSeries], error) {
	var result PaginatedResult[TVSeries]
	if err := get("/tv/on_the_air", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetNowAiringSeries: %w", err)
	}
	return &result, nil
}

func GetTopRatedSeries(page int) (*PaginatedResult[TVSeries], error) {
	var result PaginatedResult[TVSeries]
	if err := get("/tv/top_rated", pageParams(page), &result); err != nil {
		return nil, fmt.Errorf("GetTopRatedSeries: %w", err)
	}
	return &result, nil
}

func SearchSeries(query string, page int) (*PaginatedResult[TVSeries], error) {
	params := pageParams(page)
	params.Set("query", query)

	var result PaginatedResult[TVSeries]
	if err := get("/search/tv", params, &result); err != nil {
		return nil, fmt.Errorf("SearchSeries(%q): %w", query, err)
	}
	return &result, nil
}

func GetSeriesGenres() ([]Genre, error) {
	var result struct {
		Genres []Genre `json:"genres"`
	}
	if err := get("/genre/tv/list", nil, &result); err != nil {
		return nil, fmt.Errorf("GetSeriesGenres: %w", err)
	}
	return result.Genres, nil
}

func DiscoverMovies(withGenres string, page int) (*PaginatedResult[Movie], error) {
	params := pageParams(page)
	orGenres := strings.ReplaceAll(withGenres, ",", "|")
	params.Set("with_genres", orGenres)
	params.Set("sort_by", "popularity.desc")

	var result PaginatedResult[Movie]
	if err := get("/discover/movie", params, &result); err != nil {
		return nil, fmt.Errorf("DiscoverMovies: %w", err)
	}
	return &result, nil
}

func SearchPerson(query string, page int) (*PaginatedResult[Person], error) {
	params := pageParams(page)
	params.Set("query", query)

	var result PaginatedResult[Person]
	if err := get("/search/person", params, &result); err != nil {
		return nil, fmt.Errorf("SearchPerson(%q): %w", query, err)
	}
	return &result, nil
}

func GetPersonMovieCredits(personID int) (*PersonMovieCredits, error) {
	var result PersonMovieCredits
	path := fmt.Sprintf("/person/%d/movie_credits", personID)

	if err := get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("GetPersonMovieCredits(%d): %w", personID, err)
	}
	return &result, nil
}

func GetPersonTVCredits(personID int) (*PersonTVCredits, error) {
	var result PersonTVCredits
	path := fmt.Sprintf("/person/%d/tv_credits", personID)

	if err := get(path, nil, &result); err != nil {
		return nil, fmt.Errorf("GetPersonTVCredits(%d): %w", personID, err)
	}
	return &result, nil
}
