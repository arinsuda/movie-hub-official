package admin_module

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

func (h *Handler) GetOverview(c fiber.Ctx) error {
	stats, err := h.svc.GetOverview()
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"overview": stats})
}

func (h *Handler) GetGrowth(c fiber.Ctx) error {
	growth, err := h.svc.GetGrowth()
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"growth": growth})
}

func (h *Handler) ListUsers(c fiber.Ctx) error {
	var filter UserFilter
	if err := c.Bind().Query(&filter); err != nil {
		return badRequest(c, "invalid query parameters")
	}
	res, err := h.svc.ListUsers(filter)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(res)
}

func (h *Handler) UpdateUserRole(c fiber.Ctx) error {
	claims := mw.GetClaims(c)
	if claims == nil || claims.UserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	targetUserID, err := parseIDParam(c, "userId")
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	var req UpdateRoleRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateUserRole(claims.UserID, targetUserID, req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message": "user role updated successfully",
		"user_id": targetUserID,
		"role":    req.Role,
	})
}

func (h *Handler) UpdateUserStatus(c fiber.Ctx) error {
	claims := mw.GetClaims(c)
	if claims == nil || claims.UserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	targetUserID, err := parseIDParam(c, "userId")
	if err != nil {
		return badRequest(c, "invalid user id")
	}

	var req UpdateStatusRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	if err := h.svc.UpdateUserStatus(claims.UserID, targetUserID, req); err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{
		"message":   "user status updated successfully",
		"user_id":   targetUserID,
		"is_active": req.IsActive,
	})
}

func (h *Handler) ListReviews(c fiber.Ctx) error {
	var filter ReviewFilter
	if err := c.Bind().Query(&filter); err != nil {
		return badRequest(c, "invalid query parameters")
	}
	res, err := h.svc.ListReviews(filter)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(res)
}

func (h *Handler) DeleteReview(c fiber.Ctx) error {
	claims := mw.GetClaims(c)
	if claims == nil || claims.UserID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	reviewID, err := parseIDParam(c, "reviewId")
	if err != nil {
		return badRequest(c, "invalid review id")
	}

	var req DeleteReviewRequest
	_ = c.Bind().JSON(&req)

	if err := h.svc.DeleteReview(claims.UserID, reviewID, req); err != nil {
		return handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ListAuditLogs(c fiber.Ctx) error {
	var filter AuditLogFilter
	if err := c.Bind().Query(&filter); err != nil {
		return badRequest(c, "invalid query parameters")
	}
	res, err := h.svc.ListAuditLogs(filter)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(res)
}

func parseIDParam(c fiber.Ctx, param string) (uint, error) {
	id, err := strconv.Atoi(c.Params(param))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return uint(id), nil
}

func handleError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
	case errors.Is(err, ErrReviewNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "review not found"})
	case errors.Is(err, ErrRoleNotFound):
		return badRequest(c, "invalid role: must be 'admin' or 'user'")
	case errors.Is(err, ErrSelfDeactivation):
		return badRequest(c, "cannot deactivate your own account")
	case errors.Is(err, ErrUserAlreadyInStatus):
		return badRequest(c, "user is already in requested status")
	case errors.Is(err, ErrFinalAdminProtection):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "cannot demote or deactivate final active admin"})
	case errors.Is(err, ErrInactiveUserRoleChange):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "cannot change role of inactive user"})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}
