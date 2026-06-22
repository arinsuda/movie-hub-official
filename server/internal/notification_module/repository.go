package notification_module

import (
	"errors"

	"gorm.io/gorm"
)

var ErrNotificationNotFound = errors.New("notification not found")

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

// Create บันทึก notification ใหม่
func (r *repository) Create(n *Notification) error {
	return r.db.Create(n).Error
}

// CreateBatch บันทึกหลาย notification พร้อมกัน (ใช้สำหรับ fan-out ไปหา followers)
func (r *repository) CreateBatch(ns []Notification) error {
	if len(ns) == 0 {
		return nil
	}
	return r.db.Create(&ns).Error
}

// FindByUser ดึง notification ของ user พร้อม filter และ pagination
func (r *repository) FindByUser(userID uint, q ListNotificationsQuery) ([]Notification, int64, error) {
	page, pageSize := normalizePagination(q.Page, q.PageSize)

	db := r.db.Model(&Notification{}).Where("user_id = ?", userID)

	if q.Unread != nil && *q.Unread {
		db = db.Where("is_read = false")
	}
	if q.Type != "" {
		db = db.Where("type = ?", q.Type)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []Notification
	err := db.
		Order("created_at DESC").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Find(&rows).Error

	return rows, total, err
}

// CountUnread นับ notification ที่ยังไม่อ่านของ user
func (r *repository) CountUnread(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error
	return count, err
}

// MarkRead mark notification ที่ระบุ (ids) ว่าอ่านแล้ว
// ถ้า ids ว่าง = mark ทั้งหมดของ user
func (r *repository) MarkRead(userID uint, ids []uint) error {
	db := r.db.Model(&Notification{}).
		Where("user_id = ? AND is_read = false", userID)

	if len(ids) > 0 {
		db = db.Where("id IN ?", ids)
	}

	return db.Updates(map[string]any{
		"is_read": true,
		"read_at": gorm.Expr("NOW()"),
	}).Error
}

// DeleteByUser ลบ notification ที่ระบุของ user (soft delete ผ่าน gorm.Model)
func (r *repository) DeleteByUser(userID uint, ids []uint) error {
	result := r.db.
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
