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

func (h *Handler) GetUnreadCount(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	count, err := h.svc.GetUnreadCount(c.Context(), claims.UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{"unread_count": count})
}

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

func (h *Handler) DeleteNotifications(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	if c.Query("all") == "true" {
		if err := h.svc.DeleteAllNotifications(c.Context(), claims.UserID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
		}
		return c.Status(fiber.StatusNoContent).Send(nil)
	}

	var req DeleteNotificationsRequest
	if err := c.Bind().JSON(&req); err != nil || len(req.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ids required (or use ?all=true)"})
	}
	if err := h.svc.DeleteNotifications(c.Context(), claims.UserID, req.IDs); err != nil {
		if err == ErrNotificationNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "notification not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.Status(fiber.StatusNoContent).Send(nil)
}
