package feed_module

import (
	"database/sql"
	"errors"
	"time"

	"gorm.io/gorm"
)

var (
	ErrActivityNotFound = errors.New("activity not found")
	ErrForbidden        = errors.New("forbidden")
	ErrUserNotFound     = errors.New("user not found")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// ── Activity events ──────────────────────────────────────────────

func (r *repository) Create(event *ActivityEvent) error {
	return r.db.Create(event).Error
}

func (r *repository) FindByID(id uint) (*ActivityEvent, error) {
	var e ActivityEvent
	err := r.db.First(&e, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrActivityNotFound
	}
	return &e, err
}

func (r *repository) UpdateVisibility(id uint, isVisible bool) error {
	result := r.db.Model(&ActivityEvent{}).Where("id = ?", id).Update("is_visible", isVisible)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrActivityNotFound
	}
	return nil
}

// feedRow คือ shape ของแถวที่ได้จาก join activity_events กับ users (และ user_follows ถ้าเป็น following feed)
type feedRow struct {
	ID   uint
	Type ActivityType

	ActorID          uint
	ActorUsername    string
	ActorDisplayName *string
	ActorAvatarURL   *string

	MediaID   *int
	MediaType *string

	ReviewID      *uint
	CommentID     *uint
	AchievementID *uint
	LibraryItemID *uint

	Message   string
	CreatedAt time.Time
}

const feedSelectColumns = `
	ae.id, ae.type,
	ae.actor_id, u.username AS actor_username, u.display_name AS actor_display_name, u.avatar_url AS actor_avatar_url,
	ae.media_id, ae.media_type,
	ae.review_id, ae.comment_id, ae.achievement_id, ae.library_item_id,
	ae.message, ae.created_at
`

// FindFeed ดึง activity ของคนที่ userID follow อยู่แบบ accepted เท่านั้น (following feed)
func (r *repository) FindFeed(userID uint, pq PaginationQuery) ([]feedRow, int64, error) {
	buildQuery := func() *gorm.DB {
		return r.db.Table("activity_events ae").
			Joins("JOIN user_follows uf ON uf.followee_id = ae.actor_id AND uf.status = 'accepted'").
			Joins("JOIN users u ON u.id = ae.actor_id AND u.is_active = true").
			Where("uf.follower_id = ? AND ae.is_visible = true", userID)
	}

	var total int64
	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []feedRow
	err := buildQuery().
		Select(feedSelectColumns).
		Order("ae.created_at DESC").
		Limit(pq.Limit).Offset(pq.Offset()).
		Scan(&rows).Error

	return rows, total, err
}

// FindUserActivities ดึง activity ของ user คนเดียว ใช้กับหน้า profile
func (r *repository) FindUserActivities(targetUserID uint, pq PaginationQuery) ([]feedRow, int64, error) {
	buildQuery := func() *gorm.DB {
		return r.db.Table("activity_events ae").
			Joins("JOIN users u ON u.id = ae.actor_id AND u.is_active = true").
			Where("ae.actor_id = ? AND ae.is_visible = true", targetUserID)
	}

	var total int64
	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []feedRow
	err := buildQuery().
		Select(feedSelectColumns).
		Order("ae.created_at DESC").
		Limit(pq.Limit).Offset(pq.Offset()).
		Scan(&rows).Error

	return rows, total, err
}

// canViewUserActivities เช็กสิทธิ์แบบเดียวกับ follow_module.canViewFollowList:
// ดูของตัวเองได้เสมอ / user ไม่ private ดูได้หมด / user private ต้อง follow กันแบบ accepted ก่อนถึงดูได้
func (r *repository) canViewUserActivities(requesterID, targetUserID uint) (bool, error) {
	if requesterID == targetUserID {
		return true, nil
	}

	var isPrivate bool
	err := r.db.
		Table("users").
		Select("is_private").
		Where("id = ? AND is_active = true", targetUserID).
		Row().
		Scan(&isPrivate)
	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrUserNotFound
	}
	if err != nil {
		return false, err
	}
	if !isPrivate {
		return true, nil
	}

	var count int64
	err = r.db.Table("user_follows").
		Where("follower_id = ? AND followee_id = ? AND status = 'accepted'", requesterID, targetUserID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// ── Privacy settings ─────────────────────────────────────────────

func (r *repository) FindSettings(userID uint) ([]ActivityPrivacySetting, error) {
	var rows []ActivityPrivacySetting
	err := r.db.Where("user_id = ?", userID).Find(&rows).Error
	return rows, err
}

// IsEnabled คืนค่าว่า user เปิดแชร์ activity type นี้ไหม
// ถ้ายังไม่เคยตั้งค่า (ไม่มี row) จะใช้ค่า default ตาม defaultEnabled
func (r *repository) IsEnabled(userID uint, t ActivityType) (bool, error) {
	var setting ActivityPrivacySetting
	err := r.db.Where("user_id = ? AND activity_type = ?", userID, t).First(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return defaultEnabled[t], nil
	}
	if err != nil {
		return false, err
	}
	return setting.Enabled, nil
}

func (r *repository) UpsertSetting(userID uint, t ActivityType, enabled bool) error {
	setting := ActivityPrivacySetting{UserID: userID, ActivityType: t, Enabled: enabled}
	return r.db.
		Where(ActivityPrivacySetting{UserID: userID, ActivityType: t}).
		Assign(ActivityPrivacySetting{Enabled: enabled}).
		FirstOrCreate(&setting).Error
}
