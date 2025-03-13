package entity

import (
	"time"
)

// UserRole 用户角色关联实体
type UserRole struct {
	UserID    UID       `gorm:"primaryKey"`        // 用户ID
	RoleID    UID       `gorm:"primaryKey"`        // 角色ID
	CreatedAt time.Time `gorm:"not null"`          // 创建时间
	User      *User     `gorm:"foreignKey:UserID"` // 用户关联
	Role      *Role     `gorm:"foreignKey:RoleID"` // 角色关联
}

// NewUserRole 创建用户角色关联
func NewUserRole(userID, roleID UID) *UserRole {
	return &UserRole{
		UserID:    userID,
		RoleID:    roleID,
		CreatedAt: time.Now(),
	}
}
