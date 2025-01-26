package models

import "time"

type Word struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Word        string    `json:"word" gorm:"not null"`
	Translation string    `json:"translation" gorm:"not null"`
	Example     string    `json:"example"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
