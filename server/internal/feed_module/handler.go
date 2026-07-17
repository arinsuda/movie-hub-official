package feed_module

import (
	"errors"
	"strconv"

	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type handler struct {
	svc Service
}

func newHandler(svc Service) *handler {
	return &handler{svc: svc}
}

func (h *handler) GetFeed(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var pq PaginationQuery
	if err := c.Bind().Query(&pq); err != nil {
		return badRequest(c, "invalid query parameters")
	}

	res, err := h.svc.GetFeed(c.Context(), claims.UserID, pq)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

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

	res, err := h.svc.GetUserActivities(c.Context(), targetID, claims.UserID, pq)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

func (h *handler) UpdateVisibility(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	activityID, err := parseActivityID(c)
	if err != nil {
		return badRequest(c, "invalid activity id")
	}

	var req UpdateActivityVisibilityRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateActivityVisibility(c.Context(), activityID, claims.UserID, req.Visibility); err != nil {
		return handleErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) DeleteActivity(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	activityID, err := parseActivityID(c)
	if err != nil {
		return badRequest(c, "invalid activity id")
	}

	if err := h.svc.DeleteActivity(c.Context(), activityID, claims.UserID); err != nil {
		return handleErr(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *handler) GetSettings(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	res, err := h.svc.GetSettings(c.Context(), claims.UserID)
	if err != nil {
		return handleErr(c, err)
	}
	return c.JSON(res)
}

func (h *handler) UpdateSettings(c fiber.Ctx) error {
	claims := mw.GetClaims(c)

	var req UpdateActivitySettingsRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	res, err := h.svc.UpdateSettings(c.Context(), claims.UserID, req)
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
	case errors.Is(err, ErrForbidden) || errors.Is(err, privacy_policy.ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
