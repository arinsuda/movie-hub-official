package review_module

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

// Review
func (h *Handler) CreateReview(c fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	var req CreateReviewRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	review, err := h.svc.CreateReview(userID, req)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"review": review})
}

func (h *Handler) GetUserReviews(c fiber.Ctx) error {
	userID, err := parseUserID(c)
	if err != nil {
		return badRequest(c, "invalid user id")
	}
	claims := mw.GetClaims(c)

	reviews, err := h.svc.GetUserReviews(userID, claims.UserID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"reviews": reviews})
}

func (h *Handler) GetMediaReviews(c fiber.Ctx) error {
	mediaID, err := strconv.Atoi(c.Params("mediaId"))
	if err != nil || mediaID <= 0 {
		return badRequest(c, "invalid media id")
	}

	mt := routeToMediaType(c.Params("mediaType"))
	if mt == "" {
		return badRequest(c, "invalid media type")
	}

	claims := mw.GetClaims(c)
	reviews, err := h.svc.GetMediaReviews(mediaID, mt, claims.UserID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"reviews": reviews})
}

func (h *Handler) UpdateReview(c fiber.Ctx) error {
	userID, reviewID, err := parseIDs(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	var req UpdateReviewRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	review, err := h.svc.UpdateReview(reviewID, userID, req)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"review": review})
}

func (h *Handler) DeleteReview(c fiber.Ctx) error {
	userID, reviewID, err := parseIDs(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	if err := assertSelf(c, userID); err != nil {
		return forbidden(c)
	}

	if err := h.svc.DeleteReview(reviewID, userID); err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Like
func (h *Handler) LikeReview(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	claims := mw.GetClaims(c)

	if err := h.svc.LikeReview(reviewID, claims.UserID); err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) UnlikeReview(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	claims := mw.GetClaims(c)

	if err := h.svc.UnlikeReview(reviewID, claims.UserID); err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Comment
func (h *Handler) CreateComment(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	claims := mw.GetClaims(c)

	var req CreateCommentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	comment, err := h.svc.CreateComment(reviewID, claims.UserID, req)
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"comment": comment})
}

func (h *Handler) GetComments(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	claims := mw.GetClaims(c)

	comments, err := h.svc.GetComments(reviewID, claims.UserID)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"comments": comments})
}

func (h *Handler) UpdateComment(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	commentID, err := parseCommentID(c)
	if err != nil {
		return badRequest(c, "invalid comment id")
	}
	claims := mw.GetClaims(c)

	var req UpdateCommentRequest
	if err := c.Bind().JSON(&req); err != nil {
		return badRequest(c, "invalid request body")
	}

	_ = reviewID // ใช้ใน DeleteComment เท่านั้น
	comment, err := h.svc.UpdateComment(commentID, claims.UserID, req)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(fiber.Map{"comment": comment})
}

func (h *Handler) DeleteComment(c fiber.Ctx) error {
	reviewID, err := parseReviewID(c)
	if err != nil {
		return badRequest(c, "invalid review id")
	}
	commentID, err := parseCommentID(c)
	if err != nil {
		return badRequest(c, "invalid comment id")
	}
	claims := mw.GetClaims(c)

	if err := h.svc.DeleteComment(commentID, reviewID, claims.UserID); err != nil {
		return handleError(c, err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// helpers
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

func parseUserID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("userId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid user id")
	}
	return uint(id), nil
}

func parseReviewID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("reviewId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid review id")
	}
	return uint(id), nil
}

func parseCommentID(c fiber.Ctx) (uint, error) {
	id, err := strconv.Atoi(c.Params("commentId"))
	if err != nil || id <= 0 {
		return 0, errors.New("invalid comment id")
	}
	return uint(id), nil
}

func parseIDs(c fiber.Ctx) (uint, uint, error) {
	userID, err := parseUserID(c)
	if err != nil {
		return 0, 0, err
	}
	reviewID, err := parseReviewID(c)
	return userID, reviewID, err
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
	case errors.Is(err, ErrReviewNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "review not found"})
	case errors.Is(err, ErrCommentNotFound):
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "comment not found"})
	case errors.Is(err, ErrForbidden):
		return forbidden(c)
	case errors.Is(err, ErrAlreadyLiked):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "already liked"})
	case errors.Is(err, ErrNotLiked):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "not liked"})
	case errors.Is(err, ErrInvalidWatchedAt):
		return badRequest(c, "invalid watched_at format, use YYYY-MM-DD")
	case errors.Is(err, ErrInvalidRating):
		return badRequest(c, "rating must be between 0 and 10")
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

func forbidden(c fiber.Ctx) error {
	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "forbidden"})
}
