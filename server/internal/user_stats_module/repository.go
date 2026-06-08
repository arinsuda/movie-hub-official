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

// statRow holds the joined result of user_stats + user_statuses.
type statRow struct {
	UserStat
	Level      int
	CurrentExp int
}

// GetByUserID fetches counts from the user_stats view and level/exp from user_statuses.
// Returns zero values when either record does not exist yet (new user, no activity).
func (r *repository) GetByUserID(userID uint) (*statRow, error) {
	// เปลี่ยนจาก First() → Find() เพราะ user_stats เป็น VIEW
	var stat UserStat
	if err := r.db.Where("user_id = ?", userID).Limit(1).Find(&stat).Error; err != nil {
		return nil, err
	}
	if stat.UserID == 0 {
		stat = UserStat{UserID: userID}
	}

	// UserStatus เป็น table ปกติ ใช้ First() ได้เหมือนเดิม
	var status UserStatus
	if err := r.db.Where("user_id = ?", userID).First(&status).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		status = UserStatus{UserID: userID, Level: 1, CurrentExp: 0}
	}

	return &statRow{
		UserStat:   stat,
		Level:      status.Level,
		CurrentExp: status.CurrentExp,
	}, nil
}

// UpsertStatus creates or updates the UserStatus row and applies level-up logic.
func (r *repository) UpsertStatus(userID uint, expDelta int) error {
	var status UserStatus
	err := r.db.Where("user_id = ?", userID).
		FirstOrCreate(&status, UserStatus{UserID: userID, Level: 1}).Error
	if err != nil {
		return err
	}

	status.CurrentExp += expDelta
	for status.CurrentExp >= ExpPerLevel {
		status.CurrentExp -= ExpPerLevel
		status.Level++
	}
	if status.CurrentExp < 0 {
		status.CurrentExp = 0
	}

	return r.db.Save(&status).Error
}
