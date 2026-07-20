package admin_module

import (
	"time"

	users "github.com/arinsuda/movie-hub/internal/user_module"
	"gorm.io/datatypes"
)

type AuditAction string

const (
	ActionUserDeactivated AuditAction = "USER_DEACTIVATED"
	ActionUserReactivated AuditAction = "USER_REACTIVATED"
	ActionUserRoleChanged AuditAction = "USER_ROLE_CHANGED"
	ActionReviewDeleted   AuditAction = "REVIEW_DELETED"
)

type AdminAuditLog struct {
	ID         uint           `gorm:"primarykey;autoIncrement" json:"id"`
	AdminID    uint           `gorm:"not null;index" json:"admin_id"`
	Admin      users.User     `gorm:"foreignKey:AdminID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"-"`
	Action     AuditAction    `gorm:"type:varchar(50);not null;index" json:"action"`
	TargetType string         `gorm:"type:varchar(50);not null;index" json:"target_type"`
	TargetID   uint           `gorm:"not null;index" json:"target_id"`
	Reason     *string        `gorm:"type:text" json:"reason"`
	MetaData   datatypes.JSON `gorm:"type:jsonb" json:"meta_data"`
	CreatedAt  time.Time      `gorm:"not null;index" json:"created_at"`
}

func (AdminAuditLog) TableName() string {
	return "admin_audit_logs"
}
