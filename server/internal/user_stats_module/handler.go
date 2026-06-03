package user_stats_module

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// GetUserStats GET /users/:userId/stats
// public — ใครก็ดูได้ เพราะเป็น profile stats
func (h *Handler) GetUserStats(c fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Params("userId"))
	if err != nil || userID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	stats, err := h.svc.GetUserStats(uint(userID))
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(fiber.Map{"stats": stats})
}
