package statsmodule

import (
	"errors"
	"strconv"

	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetStats(c fiber.Ctx) error {
	mediaID, mediaType, err := parseMediaParams(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	stats, err := h.svc.GetStats(mediaID, mediaType)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"stats": stats})
}

func (h *Handler) RecordView(c fiber.Ctx) error {
	return h.increment(c, FieldView, 1)
}

func (h *Handler) AddLike(c fiber.Ctx) error {
	return h.increment(c, FieldLike, 1)
}

func (h *Handler) RemoveLike(c fiber.Ctx) error {
	return h.increment(c, FieldLike, -1)
}

func (h *Handler) AddWatchlist(c fiber.Ctx) error {
	return h.increment(c, FieldWatchlist, 1)
}

func (h *Handler) RemoveWatchlist(c fiber.Ctx) error {
	return h.increment(c, FieldWatchlist, -1)
}

func (h *Handler) AddReview(c fiber.Ctx) error {
	return h.increment(c, FieldReview, 1)
}

func (h *Handler) RemoveReview(c fiber.Ctx) error {
	return h.increment(c, FieldReview, -1)
}

func (h *Handler) increment(c fiber.Ctx, field IncrementField, delta int) error {
	mediaID, mediaType, err := parseMediaParams(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	stats, err := h.svc.IncrementStat(mediaID, mediaType, field, delta)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"stats": stats})
}

func parseMediaParams(c fiber.Ctx) (int, library_module.MediaType, error) {
	mediaID, err := strconv.Atoi(c.Query("media_id"))
	if err != nil || mediaID <= 0 {
		return 0, "", errors.New("invalid media_id")
	}

	mediaType := library_module.MediaType(c.Query("media_type"))
	if mediaType != library_module.MediaMovie && mediaType != library_module.MediaSeries {
		return 0, "", errors.New("invalid media_type, must be 'movie' or 'tv'")
	}

	return mediaID, mediaType, nil
}

func handleError(c fiber.Ctx, err error) error {
	if errors.Is(err, ErrStatNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "stats not found"})
	}
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}
