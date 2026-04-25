package user_module

import (
	"time"

	"gorm.io/gorm"
)

type RoleName string

const (
	RoleAdmin RoleName = "admin"
	RoleUser  RoleName = "user"
)

type Role struct {
	ID        uint `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	RoleName  RoleName `gorm:"type:varchar(50);uniqueIndex;not null"`
}

type GenderType string

const (
	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"
	GenderOther  GenderType = "other"
)

type User struct {
	gorm.Model
	Username        string `gorm:"type:varchar(50);uniqueIndex;not null"`
	Password        string `gorm:"not null"`
	Email           string `gorm:"type:varchar(100);uniqueIndex;not null"`
	VerifiedEmailAt *time.Time
	DisplayName     *string `gorm:"type:varchar(100)"`
	Bio             *string `gorm:"type:text"`
	AvatarURL       *string `gorm:"type:varchar(255)"`
	Age             int
	Gender          GenderType `gorm:"type:varchar(20)"`
	GenderOther     *string    `gorm:"type:varchar(100)"`
	RoleID          uint       `gorm:"not null;default:2"`
	Role            Role       `gorm:"foreignKey:RoleID"`
	FavoriteGenres  *string    `gorm:"type:text"`
	ReviewCount     int        `gorm:"default:0"`
	FollowerCount   int        `gorm:"default:0"`
	FollowingCount  int        `gorm:"default:0"`
	IsPrivate       bool       `gorm:"default:false"`
	IsActive        bool       `gorm:"default:true"`
}

type EmailVerification struct {
	ID        uint      `gorm:"primarykey;autoIncrement"`
	UserID    uint      `gorm:"not null;index"`
	Token     string    `gorm:"type:varchar(64);uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}

func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

type RefreshToken struct {
	ID          uint      `gorm:"primarykey;autoIncrement"`
	UserID      uint      `gorm:"not null;index"`
	HashedToken string    `gorm:"type:varchar(255);uniqueIndex;not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	CreatedAt   time.Time
	UserAgent string `gorm:"type:varchar(255)"`
	IPAddress string `gorm:"type:varchar(45)"`
}

func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}
