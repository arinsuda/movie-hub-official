package feed_module

import (
	"bufio"
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arinsuda/movie-hub/internal/privacy_policy"
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadEnv() {
	// Try loading from workspace root or current directory
	paths := []string{"../../../.env", "../../.env", "./.env", "../.env"}
	for _, p := range paths {
		file, err := os.Open(p)
		if err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				key := strings.TrimSpace(parts[0])
				val := strings.TrimSpace(parts[1])
				val = strings.Trim(val, `"'`)
				os.Setenv(key, val)
			}
			break
		}
	}
}

func getTestDB(t *testing.T) *gorm.DB {
	loadEnv()
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	if host == "" || user == "" || dbname == "" {
		t.Skip("Postgres environment not set, skipping integration test")
		return nil
	}

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skipf("skipping test as postgres connection failed: %v", err)
		return nil
	}
	return db
}

func TestMigrationAndConstraints(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		return
	}

	ctx := context.Background()

	// Run within a transaction that rolls back to keep database clean
	tx := db.Begin()
	defer tx.Rollback()

	// 1. Verify migration applied and stored in history
	var count int64
	err := tx.Table("migration_history").Where("version = ?", "activity_feed_privacy_system_v1").Count(&count).Error
	if err != nil {
		t.Fatalf("failed to check migration history: %v", err)
	}
	if count == 0 {
		t.Error("expected migration version 'activity_feed_privacy_system_v1' to exist in history")
	}

	// 2. CHECK constraint validation: invalid visibility rejected
	invalidEvent := ActivityEvent{
		ActorID:    1,
		Type:       ActivityMediaLiked,
		Visibility: "invalid_visibility",
	}
	err = tx.Create(&invalidEvent).Error
	if err == nil {
		t.Error("expected DB check constraint to reject invalid visibility value")
	}

	// 3. Unique index validation: duplicate active media liked rejected
	mID := 100
	mType := "movie"

	active1 := ActivityEvent{
		ActorID:    1,
		Type:       ActivityMediaLiked,
		MediaID:    &mID,
		MediaType:  &mType,
		Visibility: "default",
	}
	if err := tx.Create(&active1).Error; err != nil {
		t.Fatalf("failed to create first active event: %v", err)
	}

	active2 := ActivityEvent{
		ActorID:    1,
		Type:       ActivityMediaLiked,
		MediaID:    &mID,
		MediaType:  &mType,
		Visibility: "default",
	}
	err = tx.Create(&active2).Error
	if err == nil {
		t.Error("expected database unique index to reject duplicate active media liked event")
	}

	// 4. Soft-deleted activity does not block recreation (restore-or-create lifecycle)
	// First delete active1
	if err := tx.Delete(&active1).Error; err != nil {
		t.Fatalf("failed to soft delete event: %v", err)
	}

	// Restore-or-create repository instance
	repo := newRepository(tx)

	mediaID := 100
	mediaType := "movie"
	restoreEvent := &ActivityEvent{
		ActorID:    1,
		Type:       ActivityMediaLiked,
		MediaID:    &mediaID,
		MediaType:  &mediaType,
		Visibility: "followers",
		Message:    "Restored message",
	}

	err = repo.CreateOrRestore(ctx, restoreEvent)
	if err != nil {
		t.Fatalf("expected restore of soft-deleted activity to succeed: %v", err)
	}

	var checked ActivityEvent
	err = tx.Unscoped().First(&checked, active1.ID).Error
	if err != nil {
		t.Fatalf("failed to find restored event: %v", err)
	}
	if checked.DeletedAt.Valid {
		t.Error("expected restored event deleted_at to be NULL")
	}
	if checked.Visibility != "followers" {
		t.Errorf("expected visibility to be updated to 'followers', got %s", checked.Visibility)
	}

	// 5. ON DELETE SET NULL on target_user_id
	// Create test target user
	target := users.User{
		Username: "target_test_user",
		Email:    "target@test.com",
		Password: "password",
	}
	if err := tx.Create(&target).Error; err != nil {
		t.Fatalf("failed to create target user: %v", err)
	}

	followEvent := ActivityEvent{
		ActorID:      1,
		Type:         ActivityUserFollowed,
		TargetUserID: &target.ID,
		Visibility:   "default",
	}
	if err := tx.Create(&followEvent).Error; err != nil {
		t.Fatalf("failed to create follow activity: %v", err)
	}

	// Physically delete target user
	if err := tx.Unscoped().Delete(&target).Error; err != nil {
		t.Fatalf("failed to physically delete target user: %v", err)
	}

	var followChecked ActivityEvent
	if err := tx.First(&followChecked, followEvent.ID).Error; err != nil {
		t.Fatalf("failed to find follow activity: %v", err)
	}
	if followChecked.TargetUserID != nil {
		t.Errorf("expected target_user_id to be set to NULL after physical delete of target user, got %d", *followChecked.TargetUserID)
	}
}

func TestUserAccessPolicy(t *testing.T) {
	db := getTestDB(t)
	if db == nil {
		return
	}

	tx := db.Begin()
	defer tx.Rollback()

	// Set up users
	owner := users.User{Username: "owner_user", Email: "owner@test.com", Password: "pwd", IsPrivate: true, IsActive: true}
	viewer := users.User{Username: "viewer_user", Email: "viewer@test.com", Password: "pwd", IsActive: true}
	nonFollower := users.User{Username: "non_follower_user", Email: "non@test.com", Password: "pwd", IsActive: true}

	if err := tx.Create(&owner).Error; err != nil {
		t.Fatalf("failed to create owner: %v", err)
	}
	if err := tx.Create(&viewer).Error; err != nil {
		t.Fatalf("failed to create viewer: %v", err)
	}
	if err := tx.Create(&nonFollower).Error; err != nil {
		t.Fatalf("failed to create non-follower: %v", err)
	}

	// Create accepted follow relation: viewer follows owner
	follow := map[string]any{
		"follower_id": viewer.ID,
		"followee_id": owner.ID,
		"status":      "accepted",
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	if err := tx.Table("user_follows").Create(&follow).Error; err != nil {
		t.Fatalf("failed to create follow: %v", err)
	}

	policy := privacy_policy.NewUserAccessPolicy(tx)
	ctx := context.Background()

	// 1. Viewer (accepted follower) can view profile section of private owner
	allowed, err := policy.CanViewProfileSection(ctx, viewer.ID, owner.ID, privacy_policy.SectionProfile)
	if err != nil {
		t.Fatalf("error checking policy: %v", err)
	}
	if !allowed {
		t.Error("expected accepted follower to be allowed to view private profile section")
	}

	// 2. NonFollower cannot view profile section of private owner
	allowed, err = policy.CanViewProfileSection(ctx, nonFollower.ID, owner.ID, privacy_policy.SectionProfile)
	if err != nil {
		t.Fatalf("error checking policy: %v", err)
	}
	if allowed {
		t.Error("expected non-follower to be denied profile section of private owner")
	}

	// 3. CanViewActivity default override check for private account
	// An activity visibility "public" from a private account must restrict to followers only
	allowed, err = policy.CanViewActivity(ctx, nonFollower.ID, owner.ID, ActivityReviewCreated, privacy_policy.VisibilityPublic)
	if err != nil {
		t.Fatalf("error checking activity view policy: %v", err)
	}
	if allowed {
		t.Error("expected public activity from private account to be hidden from non-follower")
	}

	allowed, err = policy.CanViewActivity(ctx, viewer.ID, owner.ID, ActivityReviewCreated, privacy_policy.VisibilityPublic)
	if err != nil {
		t.Fatalf("error checking activity view policy: %v", err)
	}
	if !allowed {
		t.Error("expected public activity from private account to be visible to accepted follower")
	}
}
