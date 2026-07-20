package library_module

import (
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/feed_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router fiber.Router,
	db *gorm.DB,
	statsSvc stats.ExpAdder,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
	feed feed_module.Service,
	policy privacy_policy.UserAccessPolicy,
) {
	svc := NewService(db, statsSvc, achieve, notif, feed, policy)
	h := NewHandler(svc)

	library := router.Group("/library")
	library.Post("", h.AddItem)
	library.Get("/user/:userId", h.GetLibrary)
	library.Get("/media/:mediaType/:mediaId", h.GetMediaStatus)
	library.Delete("/:itemId", h.RemoveItem)
	library.Patch("/:itemId", h.UpdateItem)
}
