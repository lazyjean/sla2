package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// QuestionRepository 定义问题仓储接口
type QuestionRepository interface {
	// Get 根据ID获取问题
	Get(ctx context.Context, id string) (*entity.Question, error)

	// Create 创建新问题
	Create(ctx context.Context, question *entity.Question) error

	// Update 更新问题
	Update(ctx context.Context, question *entity.Question) error

	// Delete 删除问题
	Delete(ctx context.Context, id string) error

	// Search 搜索问题
	Search(ctx context.Context, keyword string, tags []string, page, pageSize int) ([]*entity.Question, int64, error)

	// CreateTag 创建问题标签
	CreateTag(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error)

	// GetTag 根据ID获取问题标签
	GetTag(ctx context.Context, id string) (*entity.QuestionTag, error)

	// UpdateTag 更新问题标签
	UpdateTag(ctx context.Context, tag *entity.QuestionTag) (*entity.QuestionTag, error)

	// DeleteTag 删除问题标签
	DeleteTag(ctx context.Context, id string) error

	// FindAllTags 查询所有标签，可以限制返回数量
	FindAllTags(ctx context.Context, limit int) ([]*entity.QuestionTag, error)
}
