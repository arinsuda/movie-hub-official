package media_stats_module

import (
	"errors"
	"strings"
	"time"

	"github.com/arinsuda/movie-hub/internal/library_module"
	"github.com/arinsuda/movie-hub/internal/like_module"
	"github.com/arinsuda/movie-hub/internal/movie_module"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrStatNotFound = errors.New("stat not found")
	ErrAlreadyLiked = errors.New("already liked")
	ErrNotLiked     = errors.New("not liked")
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// ── View Counter ──────────────────────────────────────────────────

func (r *repository) UpsertStat(mediaID int, mediaType movie_module.MediaType) error {
	stat := MediaStat{MediaID: mediaID, MediaType: mediaType}
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&stat).Error
}

func (r *repository) IncrementView(mediaID int, mediaType movie_module.MediaType) error {
	result := r.db.Model(&MediaStat{}).
		Where("media_id = ? AND media_type = ?", mediaID, mediaType).
		UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStatNotFound
	}
	return nil
}

func (r *repository) GetViewCount(mediaID int, mediaType movie_module.MediaType) (int, error) {
	var stat MediaStat
	err := r.db.Select("view_count").
		Where("media_id = ? AND media_type = ?", mediaID, mediaType).
		First(&stat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return stat.ViewCount, err
}

// ── Like ──────────────────────────────────────────────────────────

func (r *repository) CreateLike(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	like := like_module.MediaLike{
		UserID:    userID,
		MediaID:   mediaID,
		MediaType: mediaType,
	}

	err := r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "media_id"},
			{Name: "media_type"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"deleted_at": nil,
			"created_at": time.Now(),
			"updated_at": time.Now(),
		}),
	}).Create(&like).Error

	if err != nil && isDuplicateError(err) {
		return ErrAlreadyLiked
	}
	return err
}

func (r *repository) DeleteLike(userID uint, mediaID int, mediaType movie_module.MediaType) error {
	result := r.db.Where(
		"user_id = ? AND media_id = ? AND media_type = ?",
		userID, mediaID, mediaType,
	).Delete(&like_module.MediaLike{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotLiked
	}
	return nil
}

func (r *repository) CountLikes(mediaID int, mediaType movie_module.MediaType) (int, error) {
	var count int64
	err := r.db.Model(&like_module.MediaLike{}).
		Where("media_id = ? AND media_type = ?", mediaID, mediaType).
		Count(&count).Error
	return int(count), err
}

func (r *repository) IsLiked(userID uint, mediaID int, mediaType movie_module.MediaType) (bool, error) {
	var count int64
	err := r.db.Model(&like_module.MediaLike{}).
		Where("user_id = ? AND media_id = ? AND media_type = ?", userID, mediaID, mediaType).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) GetLikedAt(userID uint, mediaID int, mediaType movie_module.MediaType) (*time.Time, error) {
	var like like_module.MediaLike

	// ดึงเฉพาะ created_at มาเหมือนเดิม
	tx := r.db.Select("created_at").
		Where("user_id = ? AND media_id = ? AND media_type = ?", userID, mediaID, mediaType).
		Limit(1).
		Find(&like)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 💡 เปลี่ยนมาเช็คจาก RowsAffected แทน! ถ้าเท่ากับ 0 แปลว่าหาไม่เจอจริงๆ
	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &like.CreatedAt, nil
}

// ── Review Count ──────────────────────────────────────────────────

func (r *repository) CountReviews(mediaID int, mediaType movie_module.MediaType) (int, error) {
	var count int64
	err := r.db.Table("reviews").
		Where("media_id = ? AND media_type = ? AND is_public = true AND deleted_at IS NULL",
			mediaID, string(mediaType)).
		Count(&count).Error
	return int(count), err
}

// ── Watchlist Count ───────────────────────────────────────────────
// นับเฉพาะ list_type = 'watchlist' เท่านั้น
// favorite และ watched ไม่นับรวม
func (r *repository) CountWatchlist(mediaID int, mediaType movie_module.MediaType) (int, error) {
	var count int64
	// 💡 เพิ่มการแปลงสถานะด้วย string(mediaType) หรือตรวจสอบให้แน่ใจว่าค่าที่ส่งมาจาก Handler เป็น "movie" หรือ "tv" อย่างถูกต้อง
	err := r.db.Model(&library_module.LibraryItem{}).
		Where("media_id = ? AND media_type = ? AND list_type = ? AND deleted_at IS NULL",
			mediaID, string(mediaType), string(movie_module.ListWatchlist)).
		Count(&count).Error
	return int(count), err
}

func (r *repository) GetWatchlistedAt(userID uint, mediaID int, mediaType movie_module.MediaType) (*time.Time, error) {
	var item library_module.LibraryItem

	tx := r.db.Select("created_at").
		Where("user_id = ? AND media_id = ? AND media_type = ? AND list_type = ? AND deleted_at IS NULL",
			userID, mediaID, mediaType, movie_module.ListWatchlist).
		Limit(1).
		Find(&item)

	if tx.Error != nil {
		return nil, tx.Error
	}

	// 💡 เช็คจาก RowsAffected เช่นเดียวกันครับ
	if tx.RowsAffected == 0 {
		return nil, nil
	}

	return &item.CreatedAt, nil
}

// ── Helpers ───────────────────────────────────────────────────────

func isDuplicateError(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "duplicate key") || strings.Contains(msg, "UNIQUE constraint")
}
