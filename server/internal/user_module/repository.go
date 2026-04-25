package user_module

import (
	"errors"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Role").Where("id = ? AND is_active = true", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (r *repository) UpdateProfile(id uint, updates map[string]any) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *repository) DeleteUser(id uint) error {
	result := r.db.Model(&User{}).Where("id = ?", id).Update("is_active", false)
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return result.Error
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrForbidden    = errors.New("forbidden")
)
