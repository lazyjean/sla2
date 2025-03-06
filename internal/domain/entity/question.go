package entity

import (
	"time"
)

type QuestionID uint32

// Question 问题实体
type Question struct {
	ID          QuestionID `gorm:"primaryKey"`
	Title       string     `gorm:"type:varchar(255);not null"`
	Content     string     `gorm:"type:text;not null"`
	Type        string     `gorm:"type:varchar(50);not null"`                        // 题目类型：单选、多选、填空等
	Difficulty  int        `gorm:"type:int;not null;default:1"`                      // 难度等级：1-5
	Options     []string   `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 选项列表
	Answer      string     `gorm:"type:text;not null"`                               // 答案
	Explanation string     `gorm:"type:text"`                                        // 解析
	Tags        []string   `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 标签列表
	Status      string     `gorm:"type:varchar(50);not null;default:'draft'"`        // 状态：draft-草稿，published-已发布
	CreatedAt   time.Time  `gorm:"not null"`
	UpdatedAt   time.Time  `gorm:"not null"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}

// NewQuestion 创建新的问题实体
func NewQuestion(title, content string, tags []string, creatorID string) *Question {
	now := time.Now()
	return &Question{
		Title:     title,
		Content:   content,
		Tags:      tags,
		Status:    "draft",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Publish 发布问题
func (q *Question) Publish() {
	q.Status = "published"
	q.UpdatedAt = time.Now()
}

// Update 更新问题
func (q *Question) Update(title, content string, tags []string) {
	q.Title = title
	q.Content = content
	q.Tags = tags
	q.UpdatedAt = time.Now()
}

// Delete 删除问题
func (q *Question) Delete() {
	q.Status = "deleted"
	q.UpdatedAt = time.Now()
}
