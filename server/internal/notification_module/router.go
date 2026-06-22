package notification_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, userProvider UserProvider) {
	svc := NewService(db, userProvider)
	h := NewHandler(svc)

	notif := router.Group("/notifications")
	notif.Get("/", h.ListNotifications)
	notif.Get("/unread-count", h.GetUnreadCount)
	notif.Patch("/read", h.MarkRead)
	notif.Delete("/", h.DeleteNotifications)
}
