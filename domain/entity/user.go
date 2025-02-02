package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID        uint      `gorm:"primarykey"`
	Username  string    `gorm:"size:50;uniqueIndex"`
	Email     string    `gorm:"size:100;uniqueIndex"`
	Phone     string    `gorm:"size:20;uniqueIndex"`
	Password  string    `gorm:"size:100;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}
