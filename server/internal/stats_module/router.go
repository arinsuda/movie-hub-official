package statsmodule

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	stats := router.Group("/stats")
	stats.Get("/", h.GetStats)
	stats.Post("/view", h.RecordView)
	stats.Post("/like", h.AddLike)
	stats.Delete("/like", h.RemoveLike)
	stats.Post("/watchlist", h.AddWatchlist)
	stats.Delete("/watchlist", h.RemoveWatchlist)
	stats.Post("/review", h.AddReview)
	stats.Delete("/review", h.RemoveReview)
}
