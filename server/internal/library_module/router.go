package library_module

import (
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, exp stats.ExpAdder) {
	svc := NewService(db, exp)
	h := NewHandler(svc)

	users := router.Group("/users/:userId")
	library := users.Group("/library")

	library.Post("/", h.AddItem)
	library.Get("/", h.GetLibrary)
	library.Get("/status", h.GetMediaStatus)
	library.Patch("/:itemId", h.UpdateItem)
	library.Delete("/:itemId", h.RemoveItem)
}
