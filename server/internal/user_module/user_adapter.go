package user_module

import "gorm.io/gorm"

// NotificationUserAdapter ทำให้ user_module.repository implement
// notification_module.UserProvider interface โดยไม่ต้อง import circular
type NotificationUserAdapter struct {
	repo *repository
}

func NewNotificationUserAdapter(db *gorm.DB) *NotificationUserAdapter {
	// สมมติว่าใน user_module มีการประกาศ type repository อยู่แล้ว
	return &NotificationUserAdapter{
		repo: &repository{db: db},
	}
}

// FindByID implements notification_module.UserProvider
func (a *NotificationUserAdapter) FindByID(id uint) (*User, int, int, int, error) {
	return a.repo.FindByID(id)
}

// FindFollowerIDs คืน user IDs ทั้งหมดที่ follow actorID อยู่
// ต้องมี follows table ในโปรเจกต์ (ปรับ table/column name ให้ตรง)
func (a *NotificationUserAdapter) FindFollowerIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := a.repo.db.
		Table("follows").
		Where("following_id = ?", userID).
		Pluck("follower_id", &ids).Error
	return ids, err
}
