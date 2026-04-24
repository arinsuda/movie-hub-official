package movie_module

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router) {
	h := NewHandler()

	movies := router.Group("/movies")

	movies.Get("/popular", h.GetPopular)
	movies.Get("/now-playing", h.GetNowPlaying)
	movies.Get("/top-rated", h.GetTopRated)
	movies.Get("/upcoming", h.GetUpcoming)
	movies.Get("/search", h.Search)
	movies.Get("/genres", h.GetGenres)
	movies.Get("/:id", h.GetByID)
	movies.Get("/:id/similar", h.GetSimilar)
}
