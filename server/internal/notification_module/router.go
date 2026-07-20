package notification_module

import "github.com/gofiber/fiber/v3"

func RegisterRoutes(router fiber.Router, svc *Service, hub *Hub) fiber.Router {
	h := NewHandler(svc)

	notif := router.Group("/notifications")
	notif.Get("/", h.ListNotifications)
	notif.Get("/unread-count", h.GetUnreadCount)
	notif.Patch("/read", h.MarkRead)
	notif.Delete("/", h.DeleteNotifications)

	return router
}
