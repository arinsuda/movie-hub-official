package notification_module

import (
	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GET /notifications
// Query: ?unread=true&page=1&page_size=20&type=followed_you
func (h *Handler) ListNotifications(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var q ListNotificationsQuery
	if err := c.Bind().Query(&q); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid query params"})
	}

	result, err := h.svc.ListNotifications(c.Context(), claims.UserID, q)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(result)
}

// GET /notifications/unread-count
func (h *Handler) GetUnreadCount(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	count, err := h.svc.GetUnreadCount(c.Context(), claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{"unread_count": count})
}

// PATCH /notifications/read
// Body: { "ids": [1,2,3] }  — ถ้า ids ว่าง = mark ทั้งหมดว่าอ่านแล้ว
func (h *Handler) MarkRead(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var req MarkReadRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.svc.MarkRead(c.Context(), claims.UserID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// DELETE /notifications
// Body: { "ids": [1,2,3] }
func (h *Handler) DeleteNotifications(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var body struct {
		IDs []uint `json:"ids"`
	}
	if err := c.Bind().JSON(&body); err != nil || len(body.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ids required"})
	}

	if err := h.svc.DeleteNotifications(c.Context(), claims.UserID, body.IDs); err != nil {
		if err == ErrNotificationNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notification not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
