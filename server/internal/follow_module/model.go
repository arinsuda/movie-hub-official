package follow_module

import "time"

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
	FolloweeID uint         `gorm:"not null;uniqueIndex:idx_follow_pair"`
	Status     FollowStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
