package movie_module

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	h := NewHandler()

	movies := router.Group("/movies")
	movies.Get("/popular", h.GetPopular)
	movies.Get("/now-playing", h.GetNowPlaying)
	movies.Get("/top-rated", h.GetTopRated)
	movies.Get("/upcoming", h.GetUpcoming)
	movies.Get("/upcoming/:year", h.GetUpcomingByYear)
	movies.Get("/search", h.Search)
	movies.Get("/genres", h.GetGenres)
	movies.Get("/recommended", h.GetRecommended)
	movies.Get("/:id", h.GetByID)
	movies.Get("/:id/similar", h.GetSimilar)
	
	series := router.Group("/tv")
	series.Get("/popular", h.GetPopularSeries)
	series.Get("/now-airing", h.GetNowAiringSeries)
	series.Get("/top-rated", h.GetTopRatedSeries)
	series.Get("/search", h.SearchSeries)
	series.Get("/genres", h.GetSeriesGenres)
	series.Get("/:id", h.GetSeriesByID)
	series.Get("/:id/similar", h.GetSimilarSeries)

	actors := router.Group("/actors")
	actors.Get("/search", h.SearchActor)
	actors.Get("/:id/movies", h.GetMoviesByActor)
	actors.Get("/:id/tv", h.GetSeriesByActor)
}
