package postgres

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

// PostgresQuestionTagRepository PostgreSQL 问题标签仓库实现
type PostgresQuestionTagRepository struct {
	db *gorm.DB
}

// NewQuestionTagRepository 创建问题标签仓库实例
func NewQuestionTagRepository(db *gorm.DB) repository.QuestionTagRepository {
	return &PostgresQuestionTagRepository{
		db: db,
	}
}

// Create 创建问题标签
func (r *PostgresQuestionTagRepository) Create(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error) {
	err := r.db.WithContext(ctx).Create(tag).Error
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// Get 根据ID获取问题标签
func (r *PostgresQuestionTagRepository) Get(ctx context.Context, id string) (*entity.QuestionTag, error) {
	var tag entity.QuestionTag
	err := r.db.WithContext(ctx).First(&tag, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Update 更新问题标签
func (r *PostgresQuestionTagRepository) Update(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error) {
	err := r.db.WithContext(ctx).Save(tag).Error
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// Delete 删除问题标签
func (r *PostgresQuestionTagRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.QuestionTag{}, "id = ?", id).Error
}

// FindAll 查询所有标签，可以限制返回数量
func (r *PostgresQuestionTagRepository) FindAll(ctx context.Context, limit int) ([]*entity.QuestionTag, error) {
	var tags []*entity.QuestionTag

	// 创建查询
	query := r.db.WithContext(ctx)

	// 如果指定了有效的limit，则限制返回数量
	if limit > 0 {
		query = query.Limit(limit)
	}

	// 执行查询并按名称升序排序
	err := query.Order("name ASC").Find(&tags).Error
	if err != nil {
		return nil, err
	}

	return tags, nil
}

var _ repository.QuestionTagRepository = (*PostgresQuestionTagRepository)(nil)
