package models

import (
	"gorm.io/gorm"
	"time"

	role "github.com/arinsuda/movie-hub/server/internal/role_module"
)

type User struct {
	gorm.Model
	Username        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password        string    `gorm:"not null"`
	Email           string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	RoleID          uint      `gorm:"not null;default:2"`
	Role            role.Role `gorm:"foreignKey:RoleID"`
	Age             int
	Gender          string `gorm:"type:varchar(20)"`
	VerifiedEmailAt *time.Time
}
