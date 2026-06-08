package review_module

import (
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, exp stats.ExpAdder) {
	svc := NewService(db, exp)
	h := NewHandler(svc)

	// ── User reviews ─────────────────────────────────────────────
	users := router.Group("/users/:userId")
	users.Post("/reviews", h.CreateReview)
	users.Get("/reviews", h.GetUserReviews)
	users.Patch("/reviews/:reviewId", h.UpdateReview)
	users.Delete("/reviews/:reviewId", h.DeleteReview)

	// ── Media reviews & in-app rating ────────────────────────────
	// :mediaType = "movies" | "series"  (mapped to "movie" | "tv" in handler)
	router.Get("/:mediaType/:mediaId/reviews", h.GetMediaReviews)
	router.Get("/:mediaType/:mediaId/rating", h.GetMediaRating)

	// ── Like / Comment ────────────────────────────────────────────
	reviews := router.Group("/reviews/:reviewId")
	reviews.Post("/likes", h.LikeReview)
	reviews.Delete("/likes", h.UnlikeReview)
	reviews.Get("/comments", h.GetComments)
	reviews.Post("/comments", h.CreateComment)
	reviews.Patch("/comments/:commentId", h.UpdateComment)
	reviews.Delete("/comments/:commentId", h.DeleteComment)
}
