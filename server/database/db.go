package database

import (
	"log"
	"time"

	"github.com/arinsuda/movie-hub/config"
	"github.com/arinsuda/movie-hub/internal/follow_module"
	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/review_module"
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

	// รัน SQL migrations (views, functions, indexes ที่ GORM AutoMigrate ทำไม่ได้)
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
		&user_module.Role{},
		&user_module.User{},
		&user_module.EmailVerification{},
		&user_module.RefreshToken{},
		&follow_module.UserFollow{},
		&library_module.LibraryItem{},
		&review_module.Review{},
		&review_module.ReviewLike{},
		&review_module.ReviewComment{},
	)
	if err != nil {
		return err
	}
	log.Println("✅ AutoMigrate completed")
	return nil
}

// runSQLMigrations รัน SQL ที่ GORM AutoMigrate ทำไม่ได้ เช่น VIEW, FUNCTION, INDEX แบบพิเศษ
// ใช้ CREATE OR REPLACE ทั้งหมด → รันซ้ำกี่ครั้งก็ได้ (idempotent)
func runSQLMigrations(db *gorm.DB) error {
	log.Println("⏳ Running SQL migrations...")

	migrations := []struct {
		name string
		sql  string
	}{
		{
			name: "user_stats view",
			sql: `
				CREATE OR REPLACE VIEW user_stats AS
				SELECT
					u.id AS user_id,

					-- reviews ที่ยังไม่ถูก soft-delete
					COUNT(DISTINCT r.id) AS review_count,

					-- คนอื่น follow user นี้ (follower)
					COUNT(DISTINCT f_in.id) AS follower_count,

					-- user นี้ follow คนอื่น (following)
					COUNT(DISTINCT f_out.id) AS following_count

				FROM users u

				LEFT JOIN reviews r
					ON r.user_id = u.id
					AND r.deleted_at IS NULL

				LEFT JOIN user_follows f_in
					ON f_in.followee_id = u.id

				LEFT JOIN user_follows f_out
					ON f_out.follower_id = u.id

				GROUP BY u.id
			`,
		}, {
			name: "user_stats view",
			sql: `
        CREATE OR REPLACE VIEW user_stats AS
        SELECT
            u.id AS user_id,
            COUNT(DISTINCT r.id)     AS review_count,
            COUNT(DISTINCT f_in.id)  AS follower_count,
            COUNT(DISTINCT f_out.id) AS following_count
        FROM users u
        LEFT JOIN reviews r
            ON r.user_id = u.id
            AND r.deleted_at IS NULL
        LEFT JOIN user_follows f_in
            ON f_in.followee_id = u.id
            AND f_in.status = 'accepted'
        LEFT JOIN user_follows f_out
            ON f_out.follower_id = u.id
            AND f_out.status = 'accepted'
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
