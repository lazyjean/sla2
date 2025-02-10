package entity

import (
	"time"
)

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusUnspecified UserStatus = iota
	UserStatusActive                 // 正常
	UserStatusInactive               // 未激活
	UserStatusSuspended              // 已停用
)

// User 用户实体
type User struct {
	ID        uint       `gorm:"primaryKey"`
	Username  string     `gorm:"type:varchar(50);not null;uniqueIndex"`
	Password  string     `gorm:"type:varchar(100);not null"`
	Email     string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	Nickname  string     `gorm:"type:varchar(50);not null"`
	Avatar    string     `gorm:"type:varchar(255)"`
	Status    UserStatus `gorm:"type:int;not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
