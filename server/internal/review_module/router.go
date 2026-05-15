package review_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	// user reviews
	users := router.Group("/users/:userId")
	users.Post("/reviews", h.CreateReview)
	users.Get("/reviews", h.GetUserReviews)
	users.Patch("/reviews/:reviewId", h.UpdateReview)
	users.Delete("/reviews/:reviewId", h.DeleteReview)

	// media reviews — public อ่านได้ แต่ต้อง auth เพื่อเช็ค is_liked
	router.Get("/:mediaType/:mediaId/reviews", h.GetMediaReviews)

	// like/comment — ต้อง auth
	reviews := router.Group("/reviews/:reviewId")
	reviews.Post("/likes", h.LikeReview)
	reviews.Delete("/likes", h.UnlikeReview)
	reviews.Get("/comments", h.GetComments)
	reviews.Post("/comments", h.CreateComment)
	reviews.Patch("/comments/:commentId", h.UpdateComment)
	reviews.Delete("/comments/:commentId", h.DeleteComment)
}
