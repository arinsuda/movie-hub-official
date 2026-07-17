package achievementsmodule

import (
	"strconv"

	"github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type handler struct {
	svc Service
}

func newHandler(svc Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) listAllAchievements(c fiber.Ctx) error {
	var filter PaginationQuery
	if err := c.Bind().Query(&filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
	}

	res, err := h.svc.ListAllAchievements(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}

func (h *handler) getAchievement(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	filter := PaginationQuery{Page: 1, Limit: 1}
	res, err := h.svc.ListAllAchievements(filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	_ = id
	_ = res

	return fiber.NewError(fiber.StatusNotImplemented, "use service.GetByID")
}

func (h *handler) listUserAchievements(c fiber.Ctx) error {
	userID, err := strconv.ParseUint(c.Params("userID"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid userID")
	}

	var filter UserAchievementFilter
	if err := c.Bind().Query(&filter); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
	}

	if raw := c.Query("unlocked"); raw != "" {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "unlocked must be true or false")
		}
		filter.Unlocked = &b
	}

	var requesterID uint
	if claims := middleware.GetClaims(c); claims != nil {
		requesterID = claims.UserID
	}

	res, err := h.svc.ListUserAchievements(c.Context(), uint(userID), requesterID, filter)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(res)
}
