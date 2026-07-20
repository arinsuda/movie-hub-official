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

	countQuery := r.db.Model(&Achievement{}).Session(&gorm.Session{})
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Model(&Achievement{}).
		Order("id ASC").
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

	buildQuery := func() *gorm.DB {
		q := r.db.Table("achievements").
			Joins(`LEFT JOIN user_achievements
				ON user_achievements.achievement_id = achievements.id
				AND user_achievements.user_id = ?
				AND user_achievements.deleted_at IS NULL`, userID).
			Where("achievements.deleted_at IS NULL")

		if filter.ActionType != "" {
			q = q.Where("achievements.action_type = ?", filter.ActionType)
		}

		if filter.Unlocked != nil {
			if *filter.Unlocked {
				q = q.Where("user_achievements.is_unlocked = true")
			} else {

				q = q.Where("COALESCE(user_achievements.is_unlocked, false) = false")
			}
		}

		return q
	}

	if err := buildQuery().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var achievements []Achievement
	err := buildQuery().
		Select("achievements.*").
		Order("achievements.id ASC").
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&achievements).Error
	if err != nil {
		return nil, 0, err
	}

	if len(achievements) == 0 {
		return list, total, nil
	}

	ids := make([]uint, len(achievements))
	for i, a := range achievements {
		ids[i] = a.ID
	}

	var uas []UserAchievement
	if err := r.db.
		Where("user_id = ? AND achievement_id IN ?", userID, ids).
		Find(&uas).Error; err != nil {
		return nil, 0, err
	}

	uaMap := make(map[uint]UserAchievement, len(uas))
	for _, ua := range uas {
		uaMap[ua.AchievementID] = ua
	}

	list = make([]UserAchievement, len(achievements))
	for i, a := range achievements {
		if ua, ok := uaMap[a.ID]; ok {
			ua.Achievement = a
			list[i] = ua
		} else {
			list[i] = UserAchievement{
				UserID:        userID,
				AchievementID: a.ID,
				Achievement:   a,
				CurrentCount:  0,
				IsUnlocked:    false,
			}
		}
	}

	return list, total, nil
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
