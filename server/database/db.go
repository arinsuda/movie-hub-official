package database

import (
	"log"
	"os"
	"time"

	"github.com/arinsuda/movie-hub/config"
	achievementsmodule "github.com/arinsuda/movie-hub/internal/achievements_module"
	"github.com/arinsuda/movie-hub/internal/feed_module"
	"github.com/arinsuda/movie-hub/internal/follow_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/like_module"
	"github.com/arinsuda/movie-hub/internal/media_stats_module"
	noti "github.com/arinsuda/movie-hub/internal/notification_module"
	"github.com/arinsuda/movie-hub/internal/review_module"
	"github.com/arinsuda/movie-hub/internal/user_module"
	"github.com/arinsuda/movie-hub/internal/user_stats_module"
	"golang.org/x/crypto/bcrypt"

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

	if err := runSQLMigrations(db); err != nil {
		log.Fatalf("❌ SQL migrations failed: %v", err)
	}

	if err := seedInitialData(db); err != nil {
		log.Fatalf("❌ Seed roles failed: %v", err)
	}

	if err := achievementsmodule.SeedFromFile(db, "database/seeder/achievement.json"); err != nil {
		log.Fatalf("❌ Seed achievements failed: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	log.Println("⏳ Running AutoMigrate...")
	err := db.AutoMigrate(

		&user_module.Role{},
		&user_module.User{},
		&user_module.EmailVerification{},
		&user_module.RefreshToken{},
		&user_module.EmailChangeRequest{},
		&user_module.PasswordResetToken{},

		&follow_module.UserFollow{},
		&user_stats_module.UserStatus{},

		&library_module.LibraryItem{},

		&review_module.Review{},
		&review_module.ReviewLike{},
		&review_module.ReviewComment{},
		&review_module.ReviewHelpful{},

		&like_module.MediaLike{},
		&media_stats_module.MediaStat{},
		&noti.Notification{},

		&achievementsmodule.Achievement{},
		&achievementsmodule.UserAchievement{},

		// feed_module: activity feed + per-type privacy setting
		&feed_module.ActivityEvent{},
		&feed_module.ActivityPrivacySetting{},
	)
	if err != nil {
		return err
	}
	log.Println("✅ AutoMigrate completed")
	return nil
}

func runSQLMigrations(db *gorm.DB) error {
	log.Println("⏳ Running SQL migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "user_stats view",
			sql: `
                DROP VIEW IF EXISTS user_stats CASCADE;
                DROP TABLE IF EXISTS user_stats CASCADE;

                CREATE VIEW user_stats AS
                SELECT
                    u.id AS user_id,
                    COUNT(DISTINCT r.id)   AS review_count,
                    COUNT(DISTINCT ml.id)  AS like_count,
                    COUNT(DISTINCT CASE WHEN li_w.list_type  = 'watchlist' THEN li_w.id END) AS watchlist_count,
                    COUNT(DISTINCT CASE WHEN li_wd.list_type = 'watched'   THEN li_wd.id END) AS watched_count,
                    COUNT(DISTINCT f_in.id)  AS follower_count,
                    COUNT(DISTINCT f_out.id) AS following_count
                FROM users u
                LEFT JOIN reviews r ON r.user_id = u.id AND r.deleted_at IS NULL
                LEFT JOIN media_likes ml ON ml.user_id = u.id AND ml.deleted_at IS NULL
                LEFT JOIN library_items li_w ON li_w.user_id = u.id AND li_w.list_type = 'watchlist' AND li_w.deleted_at IS NULL
                LEFT JOIN library_items li_wd ON li_wd.user_id = u.id AND li_wd.list_type = 'watched' AND li_wd.deleted_at IS NULL
                LEFT JOIN user_follows f_in ON f_in.followee_id = u.id AND f_in.status = 'accepted'
                LEFT JOIN user_follows f_out ON f_out.follower_id = u.id AND f_out.status = 'accepted'
                GROUP BY u.id
            `,
		},
	}

	for _, m := range migrations {
		if err := db.Exec(m.sql).Error; err != nil {
			log.Printf("❌ Migration failed [%s]: %v", m.name, err)
			return err
		}
		log.Printf("✅ Migration applied: %s", m.name)
	}

	log.Println("✅ SQL migrations completed")
	return nil
}

func seedInitialData(db *gorm.DB) error {

	roles := []user_module.Role{
		{RoleName: user_module.RoleAdmin},
		{RoleName: user_module.RoleUser},
	}
	for _, role := range roles {
		if err := db.FirstOrCreate(&role, user_module.Role{RoleName: role.RoleName}).Error; err != nil {
			return err
		}
	}
	log.Println("✅ Roles seeded")

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminUsername := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminVerifiedEmailAt := time.Now().UTC()

	if adminEmail != "" && adminPassword != "" {
		var count int64

		db.Model(&user_module.User{}).Where("email = ?", adminEmail).Count(&count)

		if count == 0 {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
			if err != nil {
				return err
			}

			var roleAdmin user_module.Role
			db.Where("role_name = ?", user_module.RoleAdmin).First(&roleAdmin)

			adminUser := &user_module.User{
				Username:        adminUsername,
				Email:           adminEmail,
				Password:        string(hashedPassword),
				RoleID:          roleAdmin.ID,
				VerifiedEmailAt: &adminVerifiedEmailAt,
			}

			if err := db.Create(adminUser).Error; err != nil {
				return err
			}
			log.Println("✅ Admin account created from .env")
		}
	}
	return nil
}
