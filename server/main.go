package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var err error
	DB, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	defer DB.Close()

	if err := DB.Ping(context.Background()); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	log.Println("Database connected")

	app := fiber.New()

	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("PORT not set, defaulting to 8080")
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}
