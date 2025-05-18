package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Admin 管理员实体
type Admin struct {
	ID            UID       `gorm:"primaryKey;autoIncrement"`                     // 修改为自增整型
	Username      string    `gorm:"uniqueIndex;not null"`                         // 用户名
	Password      string    `gorm:"not null"`                                     // 密码（加密存储）
	Nickname      string    `gorm:"not null"`                                     // 昵称
	Email         string    `gorm:"not null"`                                     // 邮箱
	EmailVerified bool      `gorm:"column:email_verified;not null;default:false"` // 邮箱是否已验证
	Roles         []string  `gorm:"type:jsonb;serializer:json"`                   // 权限列表
	CreatedAt     time.Time `gorm:"type:timestamptz;not null"`                    // 创建时间
	UpdatedAt     time.Time `gorm:"type:timestamptz;not null"`                    // 更新时间
}

// Value 实现 driver.Valuer 接口
func (a Admin) Value() (driver.Value, error) {
	if a.Roles == nil {
		return json.Marshal([]string{}) // 显式返回空数组的 JSON
	}
	return json.Marshal(a.Roles)
}

// Scan 实现 sql.Scanner 接口
func (a *Admin) Scan(value interface{}) error {
	if value == nil {
		a.Roles = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, &a.Roles)
}

// NewAdmin 创建新的管理员实体
func NewAdmin(username, password, nickname, email string) *Admin {
	now := time.Now()
	// 如果没有提供昵称，则使用用户名作为默认昵称
	if nickname == "" {
		nickname = username
	}
	return &Admin{
		Username:      username,
		Password:      password,
		Nickname:      nickname,
		Email:         email,
		EmailVerified: false,
		Roles:         []string{"admin"},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// HasRole 检查管理员是否拥有特定权限
func (a *Admin) HasRole(permission string) bool {
	for _, p := range a.Roles {
		if p == permission {
			return true
		}
	}
	return false
}

// AddRole 添加权限
func (a *Admin) AddRole(permission string) {
	if !a.HasRole(permission) {
		a.Roles = append(a.Roles, permission)
		a.UpdatedAt = time.Now()
	}
}

// RemoveRole 移除权限
func (a *Admin) RemoveRole(permission string) {
	for i, p := range a.Roles {
		if p == permission {
			a.Roles = append(a.Roles[:i], a.Roles[i+1:]...)
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
