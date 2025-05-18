package entity

import (
	"time"
)

// RolePermission 角色权限关联实体
type RolePermission struct {
	RoleID       UID         `gorm:"primaryKey"`                // 角色ID
	PermissionID UID         `gorm:"primaryKey"`                // 权限ID
	CreatedAt    time.Time   `gorm:"type:timestamptz;not null"` // 创建时间
	Role         *Role       `gorm:"foreignKey:RoleID"`         // 角色关联
	Permission   *Permission `gorm:"foreignKey:PermissionID"`   // 权限关联
}

// NewRolePermission 创建角色权限关联
func NewRolePermission(roleID, permissionID UID) *RolePermission {
	return &RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
		CreatedAt:    time.Now(),
	}
}
