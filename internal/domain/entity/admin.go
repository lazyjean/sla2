package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Admin 管理员实体
type Admin struct {
	ID          string    `gorm:"primaryKey"`           // 管理员ID
	Username    string    `gorm:"uniqueIndex;not null"` // 用户名
	Password    string    `gorm:"not null"`             // 密码（加密存储）
	Nickname    string    `gorm:"not null"`             // 昵称
	Permissions []string  `gorm:"type:jsonb"`           // 权限列表
	CreatedAt   time.Time `gorm:"not null"`             // 创建时间
	UpdatedAt   time.Time `gorm:"not null"`             // 更新时间
}

// Value 实现 driver.Valuer 接口
func (a Admin) Value() (driver.Value, error) {
	if a.Permissions == nil {
		return "[]", nil
	}
	return json.Marshal(a.Permissions)
}

// Scan 实现 sql.Scanner 接口
func (a *Admin) Scan(value interface{}) error {
	if value == nil {
		a.Permissions = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, &a.Permissions)
}

// NewAdmin 创建新的管理员实体
func NewAdmin(username, password, nickname string) *Admin {
	now := time.Now()
	return &Admin{
		Username:    username,
		Password:    password, // 注意：密码应该在应用层进行加密
		Nickname:    nickname,
		Permissions: []string{}, // 初始化空权限列表
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// HasPermission 检查管理员是否拥有特定权限
func (a *Admin) HasPermission(permission string) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// AddPermission 添加权限
func (a *Admin) AddPermission(permission string) {
	if !a.HasPermission(permission) {
		a.Permissions = append(a.Permissions, permission)
		a.UpdatedAt = time.Now()
	}
}

// RemovePermission 移除权限
func (a *Admin) RemovePermission(permission string) {
	for i, p := range a.Permissions {
		if p == permission {
			a.Permissions = append(a.Permissions[:i], a.Permissions[i+1:]...)
			a.UpdatedAt = time.Now()
			return
		}
	}
}

// UpdateNickname 更新昵称
func (a *Admin) UpdateNickname(nickname string) {
	if a.Nickname != nickname {
		a.Nickname = nickname
		a.UpdatedAt = time.Now()
	}
}

// UpdatePassword 更新密码
func (a *Admin) UpdatePassword(password string) {
	a.Password = password // 注意：密码应该在应用层进行加密
	a.UpdatedAt = time.Now()
}
