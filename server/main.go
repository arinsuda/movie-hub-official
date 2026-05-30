package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/database"
	"github.com/arinsuda/movie-hub/internal/mailer"
	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
	"github.com/arinsuda/movie-hub/router"
	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env (dev only — production ควรใช้ env จาก system โดยตรง)
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

	m := mailer.New(cfg.SMTP)
	router.Register(app, database.DB, cfg, m)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down server...")
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
