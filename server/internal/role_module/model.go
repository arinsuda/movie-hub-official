package models

import "gorm.io/gorm"

type RoleName string

const (
	RoleAdmin RoleName = "admin"
	RoleUser  RoleName = "user"
)

type Role struct {
	gorm.Model
	RoleName RoleName `gorm:"type:varchar(50);uniqueIndex;not null"`
}
