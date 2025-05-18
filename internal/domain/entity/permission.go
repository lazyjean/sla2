package entity

import (
	"time"
)

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeAPI  PermissionType = "api"  // API访问权限
	PermissionTypeMenu PermissionType = "menu" // 菜单权限
	PermissionTypeData PermissionType = "data" // 数据权限
)

// Permission 权限实体
type Permission struct {
	ID          UID            `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null"` // 权限名称
	Description string         `gorm:"type:varchar(255)"`                      // 权限描述
	Type        PermissionType `gorm:"type:varchar(20);not null"`              // 权限类型
	Object      string         `gorm:"type:varchar(100);not null"`             // 权限对象（资源）
	Action      string         `gorm:"type:varchar(100);not null"`             // 权限操作
	CreatedAt   time.Time      `gorm:"type:timestamptz;not null"`              // 创建时间
	UpdatedAt   time.Time      `gorm:"type:timestamptz;not null"`              // 更新时间
}

// NewPermission 创建新权限
func NewPermission(name, description string, permType PermissionType, object, action string) *Permission {
	now := time.Now()
	return &Permission{
		Name:        name,
		Description: description,
		Type:        permType,
		Object:      object,
		Action:      action,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update 更新权限信息
func (p *Permission) Update(name, description string, permType PermissionType, object, action string) {
	if p.Name != name {
		p.Name = name
	}
	if p.Description != description {
		p.Description = description
	}
	if p.Type != permType {
		p.Type = permType
	}
	if p.Object != object {
		p.Object = object
	}
	if p.Action != action {
		p.Action = action
	}
	p.UpdatedAt = time.Now()
}
