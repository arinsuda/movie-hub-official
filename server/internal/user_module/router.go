package user_module

import (
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	stats "github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, mc *storage.MinIOClient, statsSvc *stats.Service) {
	mailer := NewSMTPMailer()                   
	svc := NewService(db, mc, statsSvc, mailer) 
	h := NewHandler(svc)

	users := router.Group("/users")
	users.Get("/:userId", h.GetProfile)
	users.Patch("/:userId", h.UpdateProfile)
	users.Delete("/:userId", h.DeleteUser)
	users.Patch("/:userId/genres", h.UpdateFavoriteGenres)

	// Email change flow (2 steps)
	users.Post("/:userId/email/request-change", h.RequestEmailChange)
	users.Post("/:userId/email/verify-change", h.VerifyEmailChange)
}
