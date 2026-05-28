package follow_module

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB) {
	svc := NewService(db)
	h := NewHandler(svc)

	users := router.Group("/users")

	// Follow / Unfollow
	users.Post("/:userId/follow", h.Follow)
	users.Delete("/:userId/follow", h.Unfollow)

	// Follower / Following lists
	users.Get("/:userId/followers", h.GetFollowers)
	users.Get("/:userId/following", h.GetFollowing)

	// Pending requests (เจ้าของ account เท่านั้น)
	users.Get("/:userId/follow-requests", h.GetPendingRequests)
	users.Post("/:userId/follow-requests/:followerId/accept", h.AcceptRequest)
	users.Delete("/:userId/follow-requests/:followerId", h.RejectRequest)
}
