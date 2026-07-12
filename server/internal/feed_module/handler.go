package feed_module

import (
	"errors"
	"strconv"

	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type handler struct {
	svc Service
}

func newHandler(svc Service) *handler {
	return &handler{svc: svc}
}

// GetFeed handles GET /feed — timeline ของคนที่ผู้ใช้ login อยู่ follow
func (h *handler) GetFeed(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var pq PaginationQuery
	if err := c.Bind().Query(&pq); err != nil {
		return badRequest(c, "invalid query parameters")
	}

	res, err := h.svc.GetFeed(claims.UserID, pq)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

// GetUserActivities handles GET /users/:userId/activities — ใช้ในหน้า profile
func (h *handler) GetUserActivities(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	targetID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	var pq PaginationQuery
	if err := c.Bind().Query(&pq); err != nil {
		return badRequest(c, "invalid query parameters")
	}

	res, err := h.svc.GetUserActivities(targetID, claims.UserID, pq)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

// UpdateVisibility handles PATCH /activities/:activityId/visibility — hide/unhide activity รายอัน
func (h *handler) UpdateVisibility(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	activityID, err := parseActivityID(c)
	if err != nil {
		return badRequest(c, "invalid activity id")
	}

	var req UpdateVisibilityRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateVisibility(activityID, claims.UserID, req.IsVisible); err != nil {
		return handleErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// GetSettings handles GET /me/activity-settings
func (h *handler) GetSettings(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	res, err := h.svc.GetSettings(claims.UserID)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

// UpdateSettings handles PATCH /me/activity-settings
func (h *handler) UpdateSettings(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var req UpdateActivitySettingsRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	res, err := h.svc.UpdateSettings(claims.UserID, req)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id")
	}
	return uint(id), nil
}

func parseActivityID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("activityId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid activity id")
	}
	return uint(id), nil
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}

func handleErr(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrActivityNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "activity not found"})
	case errors.Is(err, ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
