package like_module

import (
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, achieve achievementsmodule.Service, notif *notification_module.Service) {
	svc := NewService(db, achieve, notif)
	h := NewHandler(svc)

	router.Get("/users/:userId/likes", h.GetLikes)
	router.Post("/stats/:mediaType/:mediaId/like", h.Like)
	router.Delete("/stats/:mediaType/:mediaId/like", h.Unlike)
}
