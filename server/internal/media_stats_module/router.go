package media_stats_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	// กำหนด Base Path ให้มี Path Parameters สำหรับทุกสถิติ
	stats := router.Group("/stats/:media_type/:media_id")

	stats.Get("/", h.GetStats)        // GET    /stats/movie/123
	stats.Post("/view", h.RecordView) // POST   /stats/movie/123/view
}
