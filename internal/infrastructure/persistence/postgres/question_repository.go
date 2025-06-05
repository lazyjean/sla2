package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type questionRepository struct {
	*repository.GenericRepositoryImpl[*entity.Question, entity.QuestionID]
}

// NewQuestionRepository 创建问题仓储实例
func NewQuestionRepository(db *gorm.DB) repository.QuestionRepository {
	return &questionRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.Question, entity.QuestionID](db),
	}
}

// Search 搜索问题
func (r *questionRepository) Search(ctx context.Context, keyword string, tags []string, page int, pageSize int) ([]*entity.Question, int64, error) {
	db := r.DB.WithContext(ctx).Model(&entity.Question{})

	// 构建查询条件
	if keyword != "" {
		db = db.Where("title ILIKE ? OR simple_question ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if len(tags) > 0 {
		// 使用 JSONB 数组包含操作符
		for _, tag := range tags {
			db = db.Where("labels @> ?", []string{tag})
		}
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var questions []*entity.Question
	offset := (page - 1) * pageSize
	if err := db.Offset(offset).Limit(pageSize).Find(&questions).Error; err != nil {
		return nil, 0, err
	}

	return questions, total, nil
}

// Implement other repository methods as needed
var _ repository.QuestionRepository = (*questionRepository)(nil)
