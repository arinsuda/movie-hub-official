package achievementsmodule

import (
	users "github.com/arinsuda/movie-hub/internal/user_module"
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
	UserID        uint        `gorm:"uniqueIndex:idx_user_achievement;not null"`
	User          users.User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AchievementID uint        `gorm:"uniqueIndex:idx_user_achievement;not null"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CurrentCount  int         `gorm:"not null;default:0"`
	IsUnlocked    bool        `gorm:"not null;default:false;index"`
	UnlockedAt    *time.Time
}
