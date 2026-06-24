package follow_module

import (
	users "github.com/arinsuda/movie-hub/internal/user_module"
	"time"
)

type FollowStatus string

const (
	StatusPending  FollowStatus = "pending"
	StatusAccepted FollowStatus = "accepted"
)

// UserFollow เก็บความสัมพันธ์ follower → followee
//
// Logic:
//   - followee เป็น public  → status = accepted ทันที
//   - followee เป็น private → status = pending รอ approve
type UserFollow struct {
	ID         uint         `gorm:"primarykey;autoIncrement"`
	FollowerID uint         `gorm:"not null;uniqueIndex:idx_follow_pair"`
	Follower   users.User   `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE;"`
	FolloweeID uint         `gorm:"not null;uniqueIndex:idx_follow_pair"`
	Followee   users.User   `gorm:"foreignKey:FolloweeID;constraint:OnDelete:CASCADE;"`
	Status     FollowStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
