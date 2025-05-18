package entity

import (
	"time"
)

// Role 角色实体
type Role struct {
	ID          UID       `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null"` // 角色名称
	Description string    `gorm:"type:varchar(255)"`                     // 角色描述
	CreatedAt   time.Time `gorm:"type:timestamptz;not null"`             // 创建时间
	UpdatedAt   time.Time `gorm:"type:timestamptz;not null"`             // 更新时间
}

// NewRole 创建新角色
func NewRole(name, description string) *Role {
	now := time.Now()
	return &Role{
		Name:        name,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// Update 更新角色信息
func (r *Role) Update(name, description string) {
	if r.Name != name {
		r.Name = name
	}
	if r.Description != description {
		r.Description = description
	}
	r.UpdatedAt = time.Now()
}
