package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// QuestionRepository 定义问题仓储接口
type QuestionRepository interface {
	GenericRepository[*entity.Question, entity.QuestionID]
	// Search 搜索问题
	Search(ctx context.Context, keyword string, tags []string, page, pageSize int) ([]*entity.Question, int64, error)
}
