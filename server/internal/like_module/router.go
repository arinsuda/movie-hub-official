package like_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	router.Get("/users/:userId/likes", h.GetLikes)
	router.Post("/stats/:mediaType/:mediaId/like", h.Like)
	router.Delete("/stats/:mediaType/:mediaId/like", h.Unlike)
}
