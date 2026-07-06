package user_module

import "gorm.io/gorm"

type NotificationUserAdapter struct {
	repo *repository
}

func NewNotificationUserAdapter(db *gorm.DB) *NotificationUserAdapter {

	return &NotificationUserAdapter{
		repo: &repository{db: db},
	}
}

func (a *NotificationUserAdapter) FindByID(id uint) (*User, int, int, int, error) {
	return a.repo.FindByID(id)
}

func (a *NotificationUserAdapter) FindFollowerIDs(userID uint) ([]uint, error) {
	var ids []uint
	err := a.repo.db.
		Table("follows").
		Where("following_id = ?", userID).
		Pluck("follower_id", &ids).Error
	return ids, err
}
