package watch_log_module

import (
	"github.com/arinsuda/movie-hub/internal/feed_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router fiber.Router,
	db *gorm.DB,
	feedSvc feed_module.Service,
	libSvc *library_module.Service,
	policy privacy_policy.UserAccessPolicy,
) {
	svc := NewService(db, feedSvc, libSvc, policy)
	h := NewHandler(svc)

	router.Post("/media/:mediaType/:mediaId/watch-logs", h.CreateWatchLog)
	router.Get("/media/:mediaType/:mediaId/watch-logs/me", h.GetMyWatchLogs)
	router.Patch("/watch-logs/:watchLogId", h.UpdateWatchLog)
	router.Delete("/watch-logs/:watchLogId", h.DeleteWatchLog)
	router.Get("/users/:userId/watch-history", h.GetUserWatchHistory)
}
