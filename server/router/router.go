package router

import (
	"github.com/arinsuda/movie-hub/internal/auth_module"
	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config, m *mailer.Mailer) {
	app.Get("/", welcomeHandler)
	app.Get("/health", healthHandler)

	api := app.Group("/api")

	// ── Public routes ─────────────────────────────────────────────
	auth_module.RegisterRoutes(api, db, cfg, m)

	// ── Protected routes (ต้อง login ก่อนทุก route ด้านล่าง) ─────
	mw := auth_module.NewMiddleware(cfg)
	protected := api.Group("/", mw.RequireAuth)

	movie_module.RegisterRoutes(protected)
	user_module.RegisterRoutes(protected, db)
}

func welcomeHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "welcome to Movie Hub API"})
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
