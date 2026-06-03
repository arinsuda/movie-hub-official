package media_stats_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	stats := router.Group("/stats")
	stats.Get("/", h.GetStats)           // GET  /stats?media_id=&media_type=
	stats.Post("/view", h.RecordView)    // POST /stats/view?media_id=&media_type=
	stats.Post("/like", h.LikeMedia)     // POST /stats/like?media_id=&media_type=
	stats.Delete("/like", h.UnlikeMedia) // DELETE /stats/like?media_id=&media_type=
}
