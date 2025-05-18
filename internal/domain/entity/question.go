package entity

import (
	"time"
)

type QuestionID uint32

// 难度常量
const (
	// CEFR 等级
	DifficultyCefrA1 = "CEFR_A1" // 入门级
	DifficultyCefrA2 = "CEFR_A2" // 基础级
	DifficultyCefrB1 = "CEFR_B1" // 进阶级
	DifficultyCefrB2 = "CEFR_B2" // 高阶级
	DifficultyCefrC1 = "CEFR_C1" // 熟练级
	DifficultyCefrC2 = "CEFR_C2" // 精通级

	// HSK 等级
	DifficultyHsk1 = "HSK_1" // 入门级
	DifficultyHsk2 = "HSK_2" // 基础级
	DifficultyHsk3 = "HSK_3" // 进阶级
	DifficultyHsk4 = "HSK_4" // 中高级
	DifficultyHsk5 = "HSK_5" // 高级
	DifficultyHsk6 = "HSK_6" // 精通级

	// 预留其他难度等级
	/*
		// Math 等级
		DifficultyMath1 = "MATH_1" // 基础级
		DifficultyMath2 = "MATH_2" // 进阶级
		DifficultyMath3 = "MATH_3" // 中高级
		DifficultyMath4 = "MATH_4" // 高级
		DifficultyMath5 = "MATH_5" // 专家级

		// Music 等级
		DifficultyMusic1 = "MUSIC_1" // 入门级
		DifficultyMusic2 = "MUSIC_2" // 基础级
		DifficultyMusic3 = "MUSIC_3" // 进阶级
		DifficultyMusic4 = "MUSIC_4" // 中高级
		DifficultyMusic5 = "MUSIC_5" // 专业级
	*/
)

// Question 问题实体
type Question struct {
	ID             QuestionID `gorm:"primaryKey"`
	Title          string     `gorm:"type:varchar(255);not null"`
	Content        []byte     `gorm:"type:jsonb;not null"`                              // HyperText 对象
	SimpleQuestion string     `gorm:"type:text"`                                        // 简单文本内容
	Type           string     `gorm:"type:varchar(50);not null"`                        // 题目类型：单选、多选、填空等
	Difficulty     string     `gorm:"type:varchar(50);not null;default:'CEFR_A1'"`      // 难度等级：CEFR_A1, HSK_1 等
	Options        []byte     `gorm:"type:jsonb;not null;default:'[]'"`                 // 选项列表
	OptionTuples   []byte     `gorm:"type:jsonb;not null;default:'[]'"`                 // 选项双元组列表
	Answers        []string   `gorm:"type:jsonb;not null;default:'[]'"`                 // 答案列表
	Status         string     `gorm:"type:varchar(50);not null;default:'draft'"`        // 状态：draft-草稿，published-已发布
	Category       string     `gorm:"type:varchar(50)"`                                 // 题目分类
	Labels         []string   `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 标签列表
	Explanation    string     `gorm:"type:text"`                                        // 解析
	Attachments    []string   `gorm:"type:jsonb;serializer:json;not null;default:'[]'"` // 附件列表
	CorrectRate    float64    `gorm:"type:float8;not null;default:0"`                   // 正确率
	TimeLimit      uint32     `gorm:"type:int;not null;default:0"`                      // 时间限制，单位秒
	CreatedAt      time.Time  `gorm:"type:timestamptz;not null"`
	UpdatedAt      time.Time  `gorm:"type:timestamptz;not null"`
}

// TableName 指定表名
func (Question) TableName() string {
	return "questions"
}

// NewQuestion 创建新的问题实体
func NewQuestion(
	title string,
	content []byte,
	simpleQuestion string,
	questionType string,
	difficulty string, // 使用预定义的难度常量，如 DifficultyCefrA1, DifficultyHsk1 等
	options []byte,
	optionTuples []byte,
	answers []string,
	category string,
	labels []string,
	explanation string,
	attachments []string,
	timeLimit uint32,
) *Question {
	now := time.Now()
	return &Question{
		Title:          title,
		Content:        content,
		SimpleQuestion: simpleQuestion,
		Type:           questionType,
		Difficulty:     difficulty,
		Options:        options,
		OptionTuples:   optionTuples,
		Answers:        answers,
		Status:         "draft",
		Category:       category,
		Labels:         labels,
		Explanation:    explanation,
		Attachments:    attachments,
		TimeLimit:      timeLimit,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// Publish 发布问题
func (q *Question) Publish() {
	q.Status = "published"
	q.UpdatedAt = time.Now()
}

// Update 更新问题
func (q *Question) Update(
	title string,
	content []byte,
	simpleQuestion string,
	questionType string,
	difficulty string, // 使用预定义的难度常量，如 DifficultyCefrA1, DifficultyHsk1 等
	options []byte,
	optionTuples []byte,
	answers []string,
	category string,
	labels []string,
	explanation string,
	attachments []string,
	timeLimit uint32,
) {
	q.Title = title
	q.Content = content
	q.SimpleQuestion = simpleQuestion
	q.Type = questionType
	q.Difficulty = difficulty
	q.Options = options
	q.OptionTuples = optionTuples
	q.Answers = answers
	q.Category = category
	q.Labels = labels
	q.Explanation = explanation
	q.Attachments = attachments
	q.TimeLimit = timeLimit
	q.UpdatedAt = time.Now()
}

// Delete 删除问题
func (q *Question) Delete() {
	q.Status = "deleted"
	q.UpdatedAt = time.Now()
}
