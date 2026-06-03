package user_stats_module

import (
	"errors"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetByUserID(userID uint) (*UserStat, error) {
	var stat UserStat
	err := r.db.Where("user_id = ?", userID).First(&stat).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// user ยังไม่มี activity เลย — คืน zero stats แทน 404
		return &UserStat{UserID: userID}, nil
	}
	return &stat, err
}
