package router

import (
	"log"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/analytics_module"
	"github.com/arinsuda/movie-hub/internal/auth_module"
	"github.com/arinsuda/movie-hub/internal/follow_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/like_module"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/arinsuda/movie-hub/internal/media_stats_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/review_module"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config, m *mailer.Mailer) {
	mc, err := storage.NewMinIOClient(config.MinIOConfig{
		Endpoint:   cfg.MinIO.Endpoint,
		AccessKey:  cfg.MinIO.AccessKey,
		SecretKey:  cfg.MinIO.SecretKey,
		BucketName: cfg.MinIO.BucketName,
		UseSSL:     cfg.MinIO.UseSSL,
	})
	if err != nil {
		log.Fatalf("❌ Cannot connect to MinIO: %v", err)
	}
	log.Println("✅ MinIO connected")

	statsSvc := user_stats_module.NewService(db)

	app.Get("/", welcomeHandler)
	app.Get("/health", healthHandler)

	api := app.Group("/api")
	authSvc := auth_module.RegisterRoutes(api, db, cfg, m, mc)

	mw := auth_module.NewMiddleware(cfg)
	protected := api.Group("/", mw.RequireAuth)
	analytics_module.RegisterRoutes(protected, db)
	movie_module.RegisterRoutes(protected)
	library_module.RegisterRoutes(protected, db, statsSvc)
	user_module.RegisterRoutes(protected, db, mc, statsSvc, authSvc)
	follow_module.RegisterRoutes(api, db)
	review_module.RegisterRoutes(protected, db, mc, statsSvc)
	media_stats_module.RegisterRoutes(protected, db)
	like_module.RegisterRoutes(protected, db)
	user_stats_module.RegisterRoutes(protected, db)
	notification_module.RegisterRoutes(protected, db, user_module.NewNotificationUserAdapter(db))
}

func welcomeHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "welcome to Movie Hub API"})
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
