package database

import (
	"log"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	db, err := gorm.Open(postgres.Open(cfg.DB.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		log.Fatalf("❌ Cannot connect to DB: %v", err)
	}

	// Connection pool tuning
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❌ Cannot get sql.DB: %v", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(2 * time.Minute)

	DB = db
	log.Println("✅ Database connected")

	if err := autoMigrate(db); err != nil {
		log.Fatalf("❌ AutoMigrate failed: %v", err)
	}

	if err := seedRoles(db); err != nil {
		log.Fatalf("❌ Seed roles failed: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	log.Println("⏳ Running migrations...")
	err := db.AutoMigrate(
		&user_module.Role{},
		&user_module.User{},
		&user_module.EmailVerification{},
		&user_module.RefreshToken{},
	)
	if err != nil {
		return err
	}
	log.Println("✅ AutoMigrate completed")
	return nil
}

func seedRoles(db *gorm.DB) error {
	roles := []user_module.Role{
		{RoleName: user_module.RoleAdmin},
		{RoleName: user_module.RoleUser},
	}

	for _, role := range roles {
		result := db.FirstOrCreate(&role, user_module.Role{RoleName: role.RoleName})
		if result.Error != nil {
			return result.Error
		}
	}

	log.Println("✅ Roles seeded")
	return nil
}
