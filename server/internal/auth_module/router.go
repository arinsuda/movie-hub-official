package auth_module

import (
	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config, m *mailer.Mailer, s *storage.MinIOClient) *Service {
	svc := NewService(db, cfg, m, s)
	h := NewHandler(svc, cfg)
	mw := NewMiddleware(cfg)

	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Get("/verify-email", h.VerifyEmail)
	auth.Post("/resend-verification", h.ResendVerification)
	auth.Post("/logout-all", mw.RequireAuth, h.LogoutAll)

	return svc
}
