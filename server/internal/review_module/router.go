package review_module

import (
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/feed_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router fiber.Router,
	db *gorm.DB,
	mc *storage.MinIOClient,
	exp stats.ExpAdder,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
	feed feed_module.Service,
	policy privacy_policy.UserAccessPolicy,
) {
	svc := NewService(db, mc, exp, achieve, notif, feed, policy)
	h := NewHandler(svc)

	reviews := router.Group("/reviews")
	reviews.Post("", h.CreateReview)
	reviews.Get("/user/:userId", h.GetUserReviews)
	reviews.Get("/media/:mediaType/:mediaId", h.GetMediaReviews)
	reviews.Get("/media/:mediaType/:mediaId/rating", h.GetMediaRating)
	reviews.Patch("/:reviewId", h.UpdateReview)
	reviews.Delete("/:reviewId", h.DeleteReview)

	reviews.Post("/:reviewId/like", h.LikeReview)
	reviews.Delete("/:reviewId/like", h.UnlikeReview)

	reviews.Post("/:reviewId/helpful", h.MarkHelpful)
	reviews.Delete("/:reviewId/helpful", h.UnmarkHelpful)

	reviews.Post("/:reviewId/comments", h.CreateComment)
	reviews.Get("/:reviewId/comments", h.GetComments)
	reviews.Patch("/comments/:commentId", h.UpdateComment)
	reviews.Delete("/:reviewId/comments/:commentId", h.DeleteComment)
}
