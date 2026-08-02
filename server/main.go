package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/database"
	"github.com/arinsuda/movie-hub/internal/mailer"
	tmdb "github.com/arinsuda/movie-hub/internal/tmdb_module"
	"github.com/arinsuda/movie-hub/router"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
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

	m := mailer.New(cfg.Brevo)
	notifHub := router.Register(app, database.DB, cfg, m)

	fiberHandler := adaptor.FiberApp(app)

	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/socket.io") {
			notifHub.Handler().ServeHTTP(w, r)
			return
		}
		fiberHandler.ServeHTTP(w, r)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mainHandler,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("🛑 Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("❌ Server shutdown error: %v", err)
		}
	}()

	log.Printf("🚀 Server running on port %s (HTTP API & Socket.io WebSockets)", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
