package notification_module

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrNotificationNotFound = errors.New("notification not found")

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, n *Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *repository) CreateBatch(ctx context.Context, ns []Notification) error {
	if len(ns) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&ns).Error
}

func (r *repository) FindByUser(ctx context.Context, userID uint, q ListNotificationsQuery) ([]Notification, int64, error) {
	page, pageSize := normalizePagination(q.Page, q.PageSize)

	db := r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ?", userID)

	if q.Unread != nil && *q.Unread {
		db = db.Where("is_read = false")
	}
	if q.Type != "" {
		db = db.Where("type = ?", q.Type)
	}
	if q.Category != "" {
		db = db.Where("category = ?", q.Category) // NEW
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []Notification
	err := db.Order("created_at DESC").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&rows).Error

	return rows, total, err
}

// NEW: unread count แยกรายหมวด สำหรับ badge บนแท็บ
func (r *repository) CountUnreadByCategory(ctx context.Context, userID uint) ([]UnreadByCategoryResponse, error) {
	var rows []UnreadByCategoryResponse
	err := r.db.WithContext(ctx).Model(&Notification{}).
		Select("category, COUNT(*) as count").
		Where("user_id = ? AND is_read = false", userID).
		Group("category").
		Scan(&rows).Error
	return rows, err
}

func (r *repository) CountUnread(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

func (r *repository) MarkRead(ctx context.Context, userID uint, ids []uint) error {
	db := r.db.WithContext(ctx).Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID)

	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}

	return db.Updates(map[string]any{
		"is_read": true,
		"read_at": time.Now(),
	}).Error
}

func (r *repository) DeleteByUser(ctx context.Context, userID uint, ids []uint) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND id IN ?", userID, ids).
		Delete(&Notification{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrNotificationNotFound
	}
	return nil
}

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return page, pageSize
}
