package library_module

import (
	"errors"
	"strconv"

	"github.com/arinsuda/movie-hub/internal/movie_module"
	mw "github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) AddItem(c fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	var req AddItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	item, err := h.svc.AddItem(c.Context(), userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"item": item})
}

func (h *Handler) GetLibrary(c fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}
	claims := mw.GetClaims(c)

	var listType *movie_module.ListType
	if q := c.Query("list_type"); q != "" {
		lt := movie_module.ListType(q)
		listType = &lt
	}

	var mediaType *movie_module.MediaType
	if q := c.Query("media_type"); q != "" {
		mt := movie_module.MediaType(q)
		mediaType = &mt
	}

	items, err := h.svc.GetLibrary(c.Context(), userID, claims.UserID, listType, mediaType)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"items": items})
}

func (h *Handler) RemoveItem(c fiber.Ctx) error {
	userID, itemID, err := parseIDs(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	if err := h.svc.RemoveItem(c.Context(), itemID, userID); err != nil {
		return handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UpdateItem(c fiber.Ctx) error {
	userID, itemID, err := parseIDs(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	var req UpdateItemRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	item, err := h.svc.UpdateItem(c.Context(), itemID, userID, req)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(fiber.Map{"item": item})
}

func (h *Handler) GetMediaStatus(c fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	mediaID, err := strconv.Atoi(c.Query("media_id"))
	if err != nil || mediaID <= 0 {
		return badRequest(c, "invalid media_id")
	}

	mediaType := movie_module.MediaType(c.Query("media_type"))
	if mediaType != movie_module.MediaMovie && mediaType != movie_module.MediaSeries {
		return badRequest(c, "invalid media_type")
	}

	status, err := h.svc.GetMediaStatus(userID, mediaID, mediaType)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(status)
}

func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id")
	}
	return uint(id), nil
}

func parseItemID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("itemId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid item id")
	}
	return uint(id), nil
}

func parseIDs(c fiber.Ctx) (uint, uint, error) {
	userID, err := parseUserID(c)
	if err != nil {
		return 0, 0, err
	}
	itemID, err := parseItemID(c)
	return userID, itemID, err
}

func assertSelf(c fiber.Ctx, targetUserID uint) error {
	claims := mw.GetClaims(c)
	if claims == nil || claims.UserID != targetUserID {
		return ErrForbidden
	}
	return nil
}

func handleError(c fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, ErrItemNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "item not found"})
	case errors.Is(err, ErrForbidden):
		return forbidden(c)
	case errors.Is(err, ErrDuplicate):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "item already in list"})
	case errors.Is(err, ErrInvalidWatchedAt):
		return badRequest(c, "invalid watched_at format, use YYYY-MM-DD")
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal server error"})
	}
}

func badRequest(c fiber.Ctx, msg string) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": msg})
}

func forbidden(c fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
}
