package analytics_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	analytics := router.Group("/analytics")

	analytics.Get("/:mediaType/trending", h.GetTrending)
	analytics.Get("/:mediaType/:mediaId", h.GetMediaAnalytics)
}
