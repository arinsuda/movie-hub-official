package follow_module

import (
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router fiber.Router,
	db *gorm.DB,
	achieve achievementsmodule.Service,
	notif *notification_module.Service,
) {
	svc := NewService(db, achieve, notif)
	h := NewHandler(svc)

	users := router.Group("/users")
	users.Post("/:userId/follow", h.Follow)
	users.Delete("/:userId/follow", h.Unfollow)
	users.Get("/:userId/followers", h.GetFollowers)
	users.Get("/:userId/following", h.GetFollowing)
	users.Get("/:userId/follow-status", h.GetFollowStatus)
	users.Get("/:userId/follow-stats", h.GetFollowStats)
	users.Get("/:userId/follow-requests", h.GetPendingRequests)
	users.Post("/:userId/follow-requests/:followerId/accept", h.AcceptRequest)
	users.Delete("/:userId/follow-requests/:followerId", h.RejectRequest)}
