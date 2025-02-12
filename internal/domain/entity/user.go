package entity

import (
	"time"
)

// UserStatus 用户状态
type UserStatus int
type UserID uint32

const (
	UserStatusUnspecified UserStatus = iota
	UserStatusActive                 // 正常
	UserStatusInactive               // 未激活
	UserStatusSuspended              // 已停用
)

// User 用户实体
type User struct {
	ID            UserID     `gorm:"primaryKey"`
	Username      string     `gorm:"type:varchar(50);not null;uniqueIndex"`
	Password      string     `gorm:"type:varchar(100);not null"`
	Email         string     `gorm:"type:varchar(100);not null;uniqueIndex"`
	Nickname      string     `gorm:"type:varchar(50)"`
	Avatar        string     `gorm:"type:varchar(255)"`
	Phone         string     `gorm:"type:varchar(20);uniqueIndex"`  // 手机号
	AppleID       string     `gorm:"type:varchar(100);uniqueIndex"` // 苹果用户ID
	Status        UserStatus `gorm:"type:int;default:1"`
	EmailVerified bool       `gorm:"type:boolean;default:false"`
	CreatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
