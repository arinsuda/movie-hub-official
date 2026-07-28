package watch_log_module

import (
	"errors"
	"strconv"

	"github.com/arinsuda/movie-hub/internal/shared"
	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) CreateWatchLog(c fiber.Ctx) error {
	mediaType := c.Params("mediaType")
	mediaID, err := strconv.Atoi(c.Params("mediaId"))
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_MEDIA_ID", "Invalid media ID")
	}

	var req CreateWatchLogRequest
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	claims := mw.GetClaims(c)
	if claims == nil {
		return shared.WriteAPIError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	}

	resp, err := h.svc.CreateWatchLog(c.Context(), claims.UserID, mediaType, mediaID, req)
	if err != nil {
		if errors.Is(err, ErrInvalidDate) || errors.Is(err, ErrFutureDate) {
			return shared.WriteAPIError(c, fiber.StatusBadRequest, shared.ErrorCodeInvalidDate, err.Error())
		}
		if errors.Is(err, ErrInvalidVisibility) {
			return shared.WriteAPIError(c, fiber.StatusBadRequest, shared.ErrorCodeInvalidVisibility, err.Error())
		}
		return shared.WriteAPIError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *Handler) GetMyWatchLogs(c fiber.Ctx) error {
	mediaType := c.Params("mediaType")
	mediaID, err := strconv.Atoi(c.Params("mediaId"))
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_MEDIA_ID", "Invalid media ID")
	}

	claims := mw.GetClaims(c)
	if claims == nil {
		return shared.WriteAPIError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	}

	resp, err := h.svc.GetMyWatchLogs(c.Context(), claims.UserID, mediaType, mediaID)
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(resp)
}

func (h *Handler) UpdateWatchLog(c fiber.Ctx) error {
	watchLogID, err := strconv.Atoi(c.Params("watchLogId"))
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_ID", "Invalid watch log ID")
	}

	var req UpdateWatchLogRequest
	if err := c.Bind().JSON(&req); err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
	}

	claims := mw.GetClaims(c)
	if claims == nil {
		return shared.WriteAPIError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	}

	resp, err := h.svc.UpdateWatchLog(c.Context(), uint(watchLogID), claims.UserID, req)
	if err != nil {
		if errors.Is(err, ErrForbidden) || err.Error() == "forbidden" {
			return shared.WriteAPIError(c, fiber.StatusForbidden, "FORBIDDEN", err.Error())
		}
		if errors.Is(err, ErrWatchLogNotFound) || err.Error() == "watch log not found" {
			return shared.WriteAPIError(c, fiber.StatusNotFound, shared.ErrorCodeWatchLogNotFound, err.Error())
		}
		return shared.WriteAPIError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(resp)
}

func (h *Handler) DeleteWatchLog(c fiber.Ctx) error {
	watchLogID, err := strconv.Atoi(c.Params("watchLogId"))
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_ID", "Invalid watch log ID")
	}

	claims := mw.GetClaims(c)
	if claims == nil {
		return shared.WriteAPIError(c, fiber.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
	}

	err = h.svc.DeleteWatchLog(c.Context(), uint(watchLogID), claims.UserID)
	if err != nil {
		if errors.Is(err, ErrForbidden) || err.Error() == "forbidden" {
			return shared.WriteAPIError(c, fiber.StatusForbidden, "FORBIDDEN", err.Error())
		}
		if errors.Is(err, ErrWatchLogNotFound) || err.Error() == "watch log not found" {
			return shared.WriteAPIError(c, fiber.StatusNotFound, shared.ErrorCodeWatchLogNotFound, err.Error())
		}
		return shared.WriteAPIError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetUserWatchHistory(c fiber.Ctx) error {
	targetUserID, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_USER_ID", "Invalid user ID")
	}

	var pq PaginationQuery
	if err := c.Bind().Query(&pq); err != nil {
		return shared.WriteAPIError(c, fiber.StatusBadRequest, "INVALID_QUERY", "Invalid query parameters")
	}
	pq.Normalize()

	claims := mw.GetClaims(c)
	var requesterID uint
	if claims != nil {
		requesterID = claims.UserID
	}

	resp, err := h.svc.GetUserWatchHistory(c.Context(), uint(targetUserID), requesterID, pq)
	if err != nil {
		if errors.Is(err, ErrForbidden) || err.Error() == "forbidden" {
			return shared.WriteAPIError(c, fiber.StatusForbidden, "FORBIDDEN", err.Error())
		}
		return shared.WriteAPIError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(resp)
}
