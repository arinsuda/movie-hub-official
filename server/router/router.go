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
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Register(app *fiber.App, db *gorm.DB, cfg *config.Config, m *mailer.Mailer) {
	// ── Infrastructure ────────────────────────────────────────────
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

	// ── Shared Services (สร้างก่อน — modules อื่นต้องการ) ────────
	//
	// ลำดับสำคัญ:
	//   1. achieveModule  — ไม่ขึ้นกับใคร
	//   2. notifSvc       — ไม่ขึ้นกับ achievement
	//   3. ทุก module อื่น — รับ achieveModule.Service + notifSvc

	// 1. Achievement module
	achieveModule := achievementsmodule.New(db)

	// 2. Notification module
	userNotifAdapter := user_module.NewNotificationUserAdapter(db)
	notifSvc := notification_module.NewService(db, userNotifAdapter)

	// ── Routes ────────────────────────────────────────────────────
	app.Get("/", welcomeHandler)
	app.Get("/health", healthHandler)

	api := app.Group("/api")

	// Auth (ไม่ต้องการ achieve/notif)
	authSvc := auth_module.RegisterRoutes(api, db, cfg, m, mc)

	mw := auth_module.NewMiddleware(cfg)
	protected := api.Group("/", mw.RequireAuth)

	// User
	userSvc := user_module.RegisterRoutes(protected, db, mc, statsSvc, authSvc, passwordResetMailer)
	authSvc.SetUserService(userSvc)

	// Achievement (public catalog + auth user achievements)
	achieveModule.RegisterRoutes(api, mw.RequireAuth)

	// Notification
	notification_module.RegisterRoutes(protected, db, userNotifAdapter)

	// Core modules — inject achieve + notif
	follow_module.RegisterRoutes(api, db, achieveModule.Service, notifSvc)
	review_module.RegisterRoutes(protected, db, mc, statsSvc, achieveModule.Service, notifSvc)
	library_module.RegisterRoutes(protected, db, statsSvc, achieveModule.Service, notifSvc)
	like_module.RegisterRoutes(protected, db, achieveModule.Service)

	// Modules ที่ไม่ต้องการ achieve/notif (ยังคงเหมือนเดิม)
	analytics_module.RegisterRoutes(protected, db)
	movie_module.RegisterRoutes(protected)
	media_stats_module.RegisterRoutes(protected, db)
	user_stats_module.RegisterRoutes(protected, db)
}

func welcomeHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "welcome to Movie Hub API"})
}

func healthHandler(c fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}
