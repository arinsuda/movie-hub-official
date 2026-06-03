package database

import (
	"log"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/follow_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/like_module"
	"github.com/arinsuda/movie-hub/internal/media_stats_module"
	"github.com/arinsuda/movie-hub/internal/review_module"
	"github.com/arinsuda/movie-hub/internal/user_module"
	// "github.com/arinsuda/movie-hub/internal/user_stats_module"
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

	if err := seedRoles(db); err != nil {
		log.Fatalf("❌ Seed roles failed: %v", err)
	}
}

func autoMigrate(db *gorm.DB) error {
	log.Println("⏳ Running AutoMigrate...")
	err := db.AutoMigrate(
		// User
		&user_module.Role{},
		&user_module.User{},
		&user_module.EmailVerification{},
		&user_module.RefreshToken{},
		// Social
		&follow_module.UserFollow{},
		// Library (watchlist / favorite / watched)
		&library_module.LibraryItem{},
		// Review
		&review_module.Review{},
		&review_module.ReviewLike{},
		&review_module.ReviewComment{},
		// Media
		&like_module.MediaLike{}, // like ต่อ media — source of truth ของ like_count
		&media_stats_module.MediaStat{},
		// &user_stats_module.UserStat{}, // เก็บแค่ view_count
	)
	if err != nil {
		return err
	}
	log.Println("✅ AutoMigrate completed")
	return nil
}

// runSQLMigrations รัน SQL ที่ GORM AutoMigrate ทำไม่ได้ (VIEW, FUNCTION, INDEX พิเศษ)
// ใช้ CREATE OR REPLACE ทั้งหมด → idempotent รันซ้ำกี่ครั้งก็ได้
func runSQLMigrations(db *gorm.DB) error {
	log.Println("⏳ Running SQL migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{
            name: "user_stats view",
            sql: `
                -- เคลียร์ของเก่าออกให้หมดก่อนเพื่อความสะอาด ไม่ว่ามันจะเป็น Table หรือ View
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
