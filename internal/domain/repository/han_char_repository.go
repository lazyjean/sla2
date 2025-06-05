package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// HanCharRepository 汉字仓库接口
type HanCharRepository interface {
	GenericRepository[*entity.HanChar, entity.HanCharID]
	// GetByCharacter 根据字符获取汉字
	GetByCharacter(ctx context.Context, character string) (*entity.HanChar, error)
	// ListWithFilters 获取汉字列表（带过滤条件）
	ListWithFilters(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error)
	// SearchWithFilters 搜索汉字（带过滤条件）
	SearchWithFilters(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error)
}
