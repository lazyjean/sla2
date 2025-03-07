package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

// Create implements repository.QuestionRepository.
func (r *questionRepository) Create(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

// Delete implements repository.QuestionRepository.
func (r *questionRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Question{}, "id = ?", id).Error
}

// Get implements repository.QuestionRepository.
func (r *questionRepository) Get(ctx context.Context, id string) (*entity.Question, error) {
	var question entity.Question
	if err := r.db.WithContext(ctx).First(&question, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

// Search implements repository.QuestionRepository.
func (r *questionRepository) Search(ctx context.Context, keyword string, tags []string, page int, pageSize int) ([]*entity.Question, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.Question{})

	if keyword != "" {
		db = db.Where("title LIKE ? OR content LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if len(tags) > 0 {
		db = db.Joins("JOIN question_tags ON question_tags.question_id = questions.id").
			Where("question_tags.tag IN ?", tags)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var questions []*entity.Question
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&questions).Error
	return questions, total, err
}

// Update implements repository.QuestionRepository.
func (r *questionRepository) Update(ctx context.Context, question *entity.Question) error {
	return r.db.WithContext(ctx).Model(question).Updates(question).Error
}

func NewQuestionRepository(db *gorm.DB) repository.QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) FindByID(ctx context.Context, id entity.UID) (*entity.Question, error) {
	var question entity.Question
	if err := r.db.WithContext(ctx).First(&question, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &question, nil
}

// Implement other repository methods as needed
var _ repository.QuestionRepository = (*questionRepository)(nil)
