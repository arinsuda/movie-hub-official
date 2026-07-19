package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
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

	if err := runSQLMigrations(db, cfg); err != nil {
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

func runMigrationWithHistory(db *gorm.DB, version string, sql string) error {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS migration_history (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	)`).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SELECT pg_advisory_xact_lock(1029384756)").Error; err != nil {
			return err
		}

		var count int64
		if err := tx.Table("migration_history").Where("version = ?", version).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			log.Printf("Migration already applied: %s", version)
			return nil
		}

		if err := tx.Exec(sql).Error; err != nil {
			return err
		}

		if err := tx.Exec("INSERT INTO migration_history (version) VALUES (?)", version).Error; err != nil {
			return err
		}

		log.Printf("Migration applied successfully: %s", version)
		return nil
	})
}

func runSQLMigrations(db *gorm.DB, cfg *config.Config) error {
	log.Println("Running SQL migrations...")

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
			log.Printf("Migration failed [%s]: %v", m.name, err)
			return err
		}
		log.Printf("Migration applied: %s", m.name)
	}

	privacyMigrationSQL := `
		ALTER TABLE activity_events ADD COLUMN IF NOT EXISTS visibility VARCHAR(20) DEFAULT 'default' NOT NULL;
		ALTER TABLE activity_events ADD COLUMN IF NOT EXISTS target_user_id INTEGER;

		ALTER TABLE activity_events DROP CONSTRAINT IF EXISTS fk_activity_events_target_user;
		ALTER TABLE activity_events ADD CONSTRAINT fk_activity_events_target_user 
			FOREIGN KEY (target_user_id) REFERENCES users(id) ON DELETE SET NULL;

		UPDATE activity_events SET visibility = 'default' WHERE visibility IS NULL OR visibility NOT IN ('default', 'public', 'followers', 'private');
		ALTER TABLE activity_events DROP CONSTRAINT IF EXISTS check_activity_events_visibility;
		ALTER TABLE activity_events ADD CONSTRAINT check_activity_events_visibility CHECK (visibility IN ('default', 'public', 'followers', 'private'));

		CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_media ON activity_events(actor_id, type, media_id, media_type) WHERE media_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_review ON activity_events(actor_id, type, review_id) WHERE review_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_comment ON activity_events(actor_id, type, comment_id) WHERE comment_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_achievement ON activity_events(actor_id, type, achievement_id) WHERE achievement_id IS NOT NULL AND deleted_at IS NULL;
		CREATE UNIQUE INDEX IF NOT EXISTS uq_feed_follow ON activity_events(actor_id, type, target_user_id) WHERE target_user_id IS NOT NULL AND deleted_at IS NULL;

		CREATE INDEX IF NOT EXISTS idx_activity_events_actor_created_at_id ON activity_events(actor_id, created_at DESC, id DESC);
		CREATE INDEX IF NOT EXISTS idx_activity_events_created_at_id ON activity_events(created_at DESC, id DESC);
	`

	if err := runMigrationWithHistory(db, "activity_feed_privacy_system_v1", privacyMigrationSQL); err != nil {
		return err
	}

	err := runMigrationWithHistoryTx(db, "review_unique_constraint_v1", cfg.DB.MigrationLockTimeoutMs, cfg.DB.MigrationStatementTimeoutMs, func(tx *gorm.DB) error {
		// 1. Strict validation of existing index
		var indexDef string
		err := tx.Raw(`SELECT indexdef FROM pg_indexes WHERE indexname = 'uq_active_user_media_review'`).Row().Scan(&indexDef)
		if err == nil {
			indexDefUpper := strings.ToUpper(indexDef)
			hasUnique := strings.Contains(indexDefUpper, "UNIQUE")
			collapsedDef := strings.ReplaceAll(indexDefUpper, " ", "")
			hasCols := strings.Contains(collapsedDef, "(USER_ID,MEDIA_ID,MEDIA_TYPE)")
			hasPredicate := strings.Contains(indexDefUpper, "WHERE (DELETED_AT IS NULL)") || 
							strings.Contains(collapsedDef, "WHERE(DELETED_ATISNULL)")
			
			if hasUnique && hasCols && hasPredicate {
				log.Println("Index 'uq_active_user_media_review' already exists with matching definition.")
				return nil
			} else {
				return fmt.Errorf("mismatched index 'uq_active_user_media_review' exists: %s", indexDef)
			}
		}

		// 2. Duplicate preflight checks
		type dupGroup struct {
			UserID      uint            `gorm:"column:user_id"`
			MediaID     int             `gorm:"column:media_id"`
			MediaType   string          `gorm:"column:media_type"`
			ReviewIDsJS json.RawMessage `gorm:"column:review_agg"`
		}

		rows, err := tx.Raw(`
			SELECT user_id, media_id, media_type, json_agg(id ORDER BY created_at DESC, id DESC) AS review_agg
			FROM reviews
			WHERE deleted_at IS NULL
			GROUP BY user_id, media_id, media_type
			HAVING COUNT(*) > 1`).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()

		var groups []DuplicateReviewGroup
		for rows.Next() {
			var dg dupGroup
			if err := rows.Scan(&dg.UserID, &dg.MediaID, &dg.MediaType, &dg.ReviewIDsJS); err != nil {
				return err
			}
			var ids []uint
			if err := json.Unmarshal(dg.ReviewIDsJS, &ids); err != nil {
				return err
			}
			groups = append(groups, DuplicateReviewGroup{
				UserID:    dg.UserID,
				MediaID:   dg.MediaID,
				MediaType: dg.MediaType,
				ReviewIDs: ids,
			})
		}
		if err = rows.Err(); err != nil {
			return err
		}

		if len(groups) > 0 {
			return &ErrActiveDuplicateReviews{Groups: groups}
		}

		// 3. Perform step-constraint verification on historical ratings
		var invalidCount int64
		err = tx.Table("reviews").Where("deleted_at IS NULL AND (rating < 0.5 OR rating > 5.0 OR (rating * 2.0) != ROUND(rating * 2.0))").Count(&invalidCount).Error
		if err != nil {
			return err
		}
		if invalidCount > 0 {
			return fmt.Errorf("active historical ratings violate step constraint: %d rows", invalidCount)
		}

		// 4. Create the partial unique index
		if err := tx.Exec(`
			CREATE UNIQUE INDEX uq_active_user_media_review
			ON reviews (user_id, media_id, media_type)
			WHERE deleted_at IS NULL`).Error; err != nil {
			return err
		}

		// 5. Add the database-level check constraint for rating step and bounds
		return tx.Exec(`
			ALTER TABLE reviews ADD CONSTRAINT chk_reviews_rating_step
			CHECK (rating >= 0.5 AND rating <= 5.0 AND (rating * 2.0) = ROUND(rating * 2.0))`).Error
	})

	if err != nil {
		var dupErr *ErrActiveDuplicateReviews
		if errors.As(err, &dupErr) {
			groupsJS, marshalErr := json.Marshal(dupErr.Groups)
			if marshalErr != nil {
				log.Fatalf("❌ Migration failed: active duplicate reviews exist (unserializable groups: %v)", marshalErr)
			}
			log.Fatalf("❌ Migration failed: active duplicate reviews exist: %s", string(groupsJS))
		}
		return err
	}

	log.Println("SQL migrations completed")
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
