package entity

import (
	"time"

	"gorm.io/gorm"
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
	ID            UID            `gorm:"primaryKey;autoIncrement"`
	Username      string         `gorm:"type:varchar(50);uniqueIndex"`
	Password      string         `gorm:"type:varchar(100)"`
	Email         string         `gorm:"type:varchar(100);uniqueIndex"`
	Nickname      string         `gorm:"type:varchar(50)"`
	Avatar        string         `gorm:"type:varchar(255)"`
	Phone         *string        `gorm:"type:varchar(20);uniqueIndex"`  // 手机号，可以为 NULL
	AppleID       string         `gorm:"type:varchar(100);uniqueIndex"` // 苹果用户ID
	Status        UserStatus     `gorm:"type:int;not null;default:1"`
	EmailVerified bool           `gorm:"type:boolean;default:false"`
	CreatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time      `gorm:"type:timestamptz;not null;default:CURRENT_TIMESTAMP"`
	DeletedAt     gorm.DeletedAt `gorm:"index"` // 软删除
}

// GetID 获取用户ID
func (u *User) GetID() UID {
	return u.ID
}

// SetID 设置用户ID
func (u *User) SetID(id UID) {
	u.ID = id
}

// GetCreatedAt 获取创建时间
func (u *User) GetCreatedAt() time.Time {
	return u.CreatedAt
}

// SetCreatedAt 设置创建时间
func (u *User) SetCreatedAt(t time.Time) {
	u.CreatedAt = t
}

// GetUpdatedAt 获取更新时间
func (u *User) GetUpdatedAt() time.Time {
	return u.UpdatedAt
}

// SetUpdatedAt 设置更新时间
func (u *User) SetUpdatedAt(t time.Time) {
	u.UpdatedAt = t
}

// GetDeletedAt 获取删除时间
func (u *User) GetDeletedAt() gorm.DeletedAt {
	return u.DeletedAt
}

// SetDeletedAt 设置删除时间
func (u *User) SetDeletedAt(t gorm.DeletedAt) {
	u.DeletedAt = t
}
