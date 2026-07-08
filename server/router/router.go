package router

import (
	"log"

	"github.com/arinsuda/movie-hub/config"
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/analytics_module"
	"github.com/arinsuda/movie-hub/internal/auth_module"
	"github.com/arinsuda/movie-hub/internal/follow_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/like_module"
	"github.com/arinsuda/movie-hub/internal/mailer"
	"github.com/arinsuda/movie-hub/internal/media_stats_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	notification_module "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/review_module"
	"github.com/arinsuda/movie-hub/internal/shared/storage"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"github.com/arinsuda/movie-hub/internal/user_stats_module"
	"github.com/arinsuda/movie-hub/middleware"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config, m *mailer.Mailer) *notification_module.Hub {

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
	passwordResetMailer := user_module.NewSMTPPasswordResetMailer()

	achieveModule := achievementsmodule.New(db)

	mw := middleware.NewAuthMiddleware(cfg)

	userNotifAdapter := user_module.NewNotificationUserAdapter(db)

	verifier := &jwtVerifier{secret: cfg.JWT.AccessSecret}

	notifHub := notification_module.NewHub(verifier, cfg.CORS.AllowedOrigin)
	notifSvc := notification_module.NewService(db, userNotifAdapter, notifHub)

	app.Get("/", welcomeHandler)
	app.Get("/health", healthHandler)

	api := app.Group("/api")

	authSvc := auth_module.RegisterRoutes(api, db, cfg, m, mc, notifSvc)

	protected := api.Group("/", mw.RequireAuth)

	userSvc := user_module.RegisterRoutes(protected, db, mc, statsSvc, authSvc, passwordResetMailer)
	authSvc.SetUserService(userSvc)

	achieveModule.RegisterRoutes(api, mw.RequireAuth)

	notification_module.RegisterRoutes(protected, notifSvc, notifHub)

	follow_module.RegisterRoutes(api, db, achieveModule.Service, notifSvc)
	review_module.RegisterRoutes(protected, db, mc, statsSvc, achieveModule.Service, notifSvc)
	library_module.RegisterRoutes(protected, db, statsSvc, achieveModule.Service, notifSvc)
	like_module.RegisterRoutes(protected, db, achieveModule.Service, notifSvc)

	analytics_module.RegisterRoutes(protected, db)
	movie_module.RegisterRoutes(protected)
	media_stats_module.RegisterRoutes(protected, db)
	user_stats_module.RegisterRoutes(protected, db)

	return notifHub
}

type jwtVerifier struct {
	secret string
}

func (v *jwtVerifier) VerifyToken(token string) (uint, error) {
	claims, err := middleware.ParseAccess(token, v.secret)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

func welcomeHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "welcome to Movie Hub API"})
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
