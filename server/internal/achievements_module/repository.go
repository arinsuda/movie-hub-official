package achievementsmodule

import (
	"gorm.io/gorm"
)

type Repository interface {
	ListAllAchievements(filter PaginationQuery) ([]Achievement, int64, error)
	FindAchievementByID(id uint) (*Achievement, error)
	FindAchievementsByActionType(actionType string) ([]Achievement, error)

	ListUserAchievements(userID uint, filter UserAchievementFilter) ([]UserAchievement, int64, error)
	FindUserAchievement(userID, achievementID uint) (*UserAchievement, error)
	UpsertUserAchievement(ua *UserAchievement) error

	BulkUpsertAchievements(achievements []Achievement) error
}

type repository struct {
	db *gorm.DB
}

func newRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) ListAllAchievements(filter PaginationQuery) ([]Achievement, int64, error) {
	var achievements []Achievement
	var total int64

	q := r.db.Model(&Achievement{})

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.
		Order("action_type ASC, target_count ASC").
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&achievements).Error

	return achievements, total, err
}

func (r *repository) FindAchievementByID(id uint) (*Achievement, error) {
	var a Achievement
	err := r.db.First(&a, id).Error
	return &a, err
}

func (r *repository) FindAchievementsByActionType(actionType string) ([]Achievement, error) {
	var achievements []Achievement
	err := r.db.
		Where("action_type = ?", actionType).
		Order("target_count ASC").
		Find(&achievements).Error
	return achievements, err
}

func (r *repository) ListUserAchievements(userID uint, filter UserAchievementFilter) ([]UserAchievement, int64, error) {
	var list []UserAchievement
	var total int64

	q := r.db.Model(&UserAchievement{}).
		Preload("Achievement").
		Where("user_id = ?", userID)

	if filter.Unlocked != nil {
		q = q.Where("is_unlocked = ?", *filter.Unlocked)
	}

	if filter.ActionType != "" {
		q = q.Joins("JOIN achievements ON achievements.id = user_achievements.achievement_id").
			Where("achievements.action_type = ?", filter.ActionType).
			Where("achievements.deleted_at IS NULL")
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := q.
		Order("is_unlocked DESC, current_count DESC").
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&list).Error

	return list, total, err
}

func (r *repository) FindUserAchievement(userID, achievementID uint) (*UserAchievement, error) {
	var ua UserAchievement
	err := r.db.
		Where("user_id = ? AND achievement_id = ?", userID, achievementID).
		Preload("Achievement").
		First(&ua).Error
	return &ua, err
}

func (r *repository) UpsertUserAchievement(ua *UserAchievement) error {
	return r.db.Save(ua).Error
}

func (r *repository) BulkUpsertAchievements(achievements []Achievement) error {
	return r.db.
		Clauses().
		Save(&achievements).Error
}
