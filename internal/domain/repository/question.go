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
}
