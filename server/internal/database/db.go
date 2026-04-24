package database

import (
	"fmt"
	"log"
	"os"

	"github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}

	DB = db
	log.Println("✅ Database connected")

	autoMigrate()
}

func autoMigrate() {
	err := DB.AutoMigrate(
		&user_module.Role{},
		&user_module.User{},
	)
	if err != nil {
		log.Fatal("AutoMigrate failed:", err)
	}
	log.Println("✅ AutoMigrate completed")

	seedRoles()
}

func seedRoles() {
	roles := []user_module.Role{
		{RoleName: user_module.RoleAdmin},
		{RoleName: user_module.RoleUser},
	}

	for _, role := range roles {
		DB.FirstOrCreate(&role, user_module.Role{RoleName: role.RoleName})
	}

	log.Println("✅ Roles seeded")
}
