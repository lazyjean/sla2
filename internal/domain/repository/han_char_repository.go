package repository

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// HanCharRepository 汉字仓库接口
type HanCharRepository interface {
	// Create 创建汉字
	Create(ctx context.Context, hanChar *entity.HanChar) error
	// Update 更新汉字
	Update(ctx context.Context, hanChar *entity.HanChar) error
	// Delete 删除汉字
	Delete(ctx context.Context, id entity.HanCharID) error
	// GetByID 根据ID获取汉字
	GetByID(ctx context.Context, id entity.HanCharID) (*entity.HanChar, error)
	// GetByCharacter 根据字符获取汉字
	GetByCharacter(ctx context.Context, character string) (*entity.HanChar, error)
	// List 获取汉字列表
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error)
	// Search 搜索汉字
	Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error)
}
