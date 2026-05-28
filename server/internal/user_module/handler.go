package user_module

import (
	"errors"
	"strconv"

	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetProfile(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)
	requesterID := claims.UserID

	profile, err := h.svc.GetProfile(targetID, requesterID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) UpdateProfile(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)
	requesterID := claims.UserID

	var req UpdateProfileRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	profile, err := h.svc.UpdateProfile(targetID, requesterID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func (h *Handler) DeleteUser(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	if err := h.svc.DeleteUser(targetID, claims.UserID, claims.Role); err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *Handler) UpdateFavoriteGenres(c fiber.Ctx) error {
	targetID, err := parseUserID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid user id"})
	}

	claims := mw.GetClaims(c)

	var req UpdateFavoriteGenresRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	profile, err := h.svc.UpdateFavoriteGenres(targetID, claims.UserID, req.FavoriteGenres)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"user": profile})
}

func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func handleError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, ErrForbidden):
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}
