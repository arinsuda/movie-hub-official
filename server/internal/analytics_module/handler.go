package analytics_module

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

// GetMediaAnalytics handles:
//
//	GET /analytics/:mediaType/:mediaId
//
// :mediaType = "movies" | "series"
// ไม่ต้อง auth เพราะเป็น public aggregate
func (h *Handler) GetMediaAnalytics(c fiber.Ctx) error {
	mediaID, err := strconv.Atoi(c.Params("mediaId"))
	if err != nil || mediaID <= 0 {
		return badRequest(c, "invalid media id")
	}

	mt := routeToMediaType(c.Params("mediaType"))
	if mt == "" {
		return badRequest(c, "invalid media type: use 'movies' or 'series'")
	}

	result, err := h.svc.GetMediaAnalytics(mediaID, mt)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"analytics": result})
}

// GetTrending handles:
//
//	GET /analytics/:mediaType/trending
//
// คืน top 20 media ที่ trending ในหมวดนั้น
func (h *Handler) GetTrending(c fiber.Ctx) error {
	mt := routeToMediaType(c.Params("mediaType"))
	if mt == "" {
		return badRequest(c, "invalid media type: use 'movies' or 'series'")
	}

	items, err := h.svc.GetTrending(mt)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"trending": items})
}

// ── helpers ───────────────────────────────────────────────────────

// routeToMediaType แปลง URL segment → media_type value ใน DB
//
//	"movies" → "movie"
//	"series" → "tv"
func routeToMediaType(route string) string {
	switch route {
	case "movies":
		return "movie"
	case "series":
		return "tv"
	default:
		return ""
	}
}

func handleError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrInvalidMediaType):
		return badRequest(c, "media_type must be 'movie' or 'tv'")
	case errors.Is(err, ErrInvalidMediaID):
		return badRequest(c, "invalid media_id")
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}
