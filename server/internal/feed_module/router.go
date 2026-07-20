package feed_module

import (
	noti "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Module struct {
	Service Service
}

func New(db *gorm.DB, hub *noti.Hub, minio *storage.MinIOClient) *Module {
	return &Module{Service: NewService(db, hub, minio)}
}

func (m *Module) RegisterRoutes(router fiber.Router) {
	h := newHandler(m.Service)

	router.Get("/feed/new-count", h.GetNewCount)
	router.Get("/feed", h.GetFeed)
	router.Get("/users/:userId/activities", h.GetUserActivities)
	router.Patch("/activities/:activityId/visibility", h.UpdateVisibility)
	router.Delete("/activities/:activityId", h.DeleteActivity)
	router.Get("/me/activity-settings", h.GetSettings)
	router.Patch("/me/activity-settings", h.UpdateSettings)
}
