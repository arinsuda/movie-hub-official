package achievementsmodule

import (
	"gorm.io/gorm"
	"time"
)

type Achievement struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string `gorm:"type:text"`
	ExpReward   int    `gorm:"not null;default:0"`
	ActionType  string `gorm:"type:varchar(50);index;not null"`
	TargetCount int    `gorm:"not null;default:1"`
}

type UserAchievement struct {
	gorm.Model
	UserID        uint `gorm:"uniqueIndex:idx_user_achievement;not null"`
	AchievementID uint `gorm:"uniqueIndex:idx_user_achievement;not null"`
	CurrentCount  int  `gorm:"not null;default:0"`
	IsUnlocked    bool `gorm:"not null;default:false;index"`
	UnlockedAt    *time.Time
	Achievement   Achievement `gorm:"foreignKey:AchievementID"`
}
