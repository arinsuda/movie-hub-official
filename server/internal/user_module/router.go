package user_module

import (
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router fiber.Router,
	db *gorm.DB,
	mc *storage.MinIOClient,
	statsSvc StatsProvider,
	emailVerifier EmailVerificationSender,
	passwordResetMailer PasswordResetMailer,
) {
	mailer := NewSMTPMailer()
	svc := NewService(db, mc, statsSvc, mailer, emailVerifier, passwordResetMailer)
	h := NewHandler(svc)

	users := router.Group("/users")
	users.Get("/:userId", h.GetProfile)
	users.Patch("/:userId", h.UpdateProfile)
	users.Delete("/:userId", h.DeleteUser)
	users.Patch("/:userId/genres", h.UpdateFavoriteGenres)
	users.Post("/:userId/email", h.RequestEmailChange)
	users.Put("/:userId/email", h.VerifyEmailChange)
	users.Patch("/:userId/email", h.UpdateEmail)
	users.Patch("/:userId/password", h.ChangePassword)
}
