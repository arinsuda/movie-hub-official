package database

import (
	"log"
	"os"

	"github.com/arinsuda/movie-hub/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	DB = db
	log.Println("Database connected")

	autoMigrate()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&models.Role{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	log.Println("AutoMigrate completed")

	seedRoles()
}

func seedRoles() {
	roles := []models.Role{
		{RoleName: models.RoleAdmin},
		{RoleName: models.RoleUser},
	}

	for _, role := range roles {
		DB.FirstOrCreate(&role, models.Role{RoleName: role.RoleName})
	}

	log.Println("Roles seeded")
}
