package entity

import (
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
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

// ToProto 将Question实体转换为protobuf消息
func (q *Question) ToProto() *pb.Question {
	return &pb.Question{
		Id:             uint32(q.ID),
		Title:          q.Title,
		SimpleQuestion: q.Content,
		QuestionType:   pb.QuestionType_QUESTION_TYPE_UNSPECIFIED, // 这里需要根据实际类型转换
		Tags:           q.Tags,
		Difficulty:     pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_MEDIUM, // 默认难度中等，可以根据实际值调整
		Status:         pb.QuestionStatus_QUESTION_STATUS_PUBLISHED,                 // 默认已发布，可以根据实际值调整
		Explanation:    q.Explanation,
		CreatedAt:      uint64(q.CreatedAt.Unix()),
		UpdatedAt:      uint64(q.UpdatedAt.Unix()),
	}
}
