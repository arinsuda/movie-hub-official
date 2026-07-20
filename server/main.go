package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/database"
	"github.com/arinsuda/movie-hub/internal/mailer"
	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
	"github.com/arinsuda/movie-hub/router"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../.env"); err != nil {
		log.Println("⚠️  No .env file found, falling back to system env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	log.Println("HOST =", cfg.DB.Host)
	log.Println("USER =", cfg.DB.User)
	log.Println("PASS =", cfg.DB.Password)
	log.Println("DB =", cfg.DB.Name)
	log.Println("DSN =", cfg.DB.DSN())

	tmdb.Init(cfg)
	database.Connect(cfg)

	app := fiber.New(fiber.Config{
		AppName:      "movie-hub v1.0",
		ErrorHandler: customErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORS.AllowedOrigin, "http://localhost:3000", "http://localhost:5173", "https://movie-hub-official.pages.dev"},
		AllowCredentials: true,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}))

	m := mailer.New(cfg.SMTP)
	notifHub := router.Register(app, database.DB, cfg, m)

	socketAddr := ":8081"
	socketSrv := &http.Server{
		Addr: socketAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notifHub.Handler().ServeHTTP(w, r)
		}),
	}

	go func() {
		log.Printf("📡 Notification socket.io running on %s", socketAddr)
		if err := socketSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Socket.io server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := socketSrv.Shutdown(ctx); err != nil {
			log.Printf("❌ Socket.io shutdown error: %v", err)
		}

		if err := app.Shutdown(); err != nil {
			log.Printf("❌ Shutdown error: %v", err)
		}
	}()

	log.Printf("🚀 Server running on port %s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("❌ Server error: %v", err)
	}

	log.Println("👋 Server exited cleanly")
}

func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
