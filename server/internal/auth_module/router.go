package auth_module

import (
	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func RegisterRoutes(router fiber.Router, db *gorm.DB, cfg *config.Config, m *mailer.Mailer) {
	svc := NewService(db, cfg, m)
	h := NewHandler(svc, cfg)
	mw := NewMiddleware(cfg)

	auth := router.Group("/auth")

	// ── Public ────────────────────────────────────────────────────
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/refresh", h.Refresh)
	auth.Post("/logout", h.Logout)
	auth.Get("/verify-email", h.VerifyEmail)
	auth.Post("/resend-verification", h.ResendVerification)

	// ── Protected ─────────────────────────────────────────────────
	auth.Post("/logout-all", mw.RequireAuth, h.LogoutAll)
	// /me ย้ายไปอยู่ที่ GET /api/users/:userId แล้ว
}
