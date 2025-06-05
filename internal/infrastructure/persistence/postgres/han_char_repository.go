package postgres

import (
	"context"
	"errors"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// hanCharRepository PostgreSQL 汉字仓库实现
type hanCharRepository struct {
	*repository.GenericRepositoryImpl[*entity.HanChar, entity.HanCharID]
}

// NewHanCharRepository 创建汉字仓库实例
func NewHanCharRepository(db *gorm.DB) repository.HanCharRepository {
	return &hanCharRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.HanChar, entity.HanCharID](db),
	}
}

// Create 创建汉字
func (r *hanCharRepository) Create(ctx context.Context, hanChar *entity.HanChar) error {
	// 使用 GORM 的 Create 方法，GORM 会自动处理 JSONB 字段的序列化
	return r.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "character"}},
			DoNothing: true,
		}).
		Create(hanChar).Error
}

// GetByCharacter 根据字符获取汉字
func (r *hanCharRepository) GetByCharacter(ctx context.Context, character string) (*entity.HanChar, error) {
	var hanChar entity.HanChar
	err := r.DB.WithContext(ctx).Where("character = ?", character).First(&hanChar).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &hanChar, nil
}

// ListWithFilters 获取汉字列表（带过滤条件）
func (r *hanCharRepository) ListWithFilters(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error) {
	var hanChars []*entity.HanChar
	var total int64

	query := r.DB.WithContext(ctx).Model(&entity.HanChar{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level":
			if v, ok := value.(valueobject.WordDifficultyLevel); ok {
				query = query.Where("level = ?", v)
			}
		case "tags":
			if v, ok := value.([]string); ok && len(v) > 0 {
				query = query.Where("tags @> ?", v)
			}
		case "categories":
			if v, ok := value.([]string); ok && len(v) > 0 {
				query = query.Where("categories @> ?", v)
			}
		case "exclude_ids":
			if v, ok := value.([]uint); ok && len(v) > 0 {
				query = query.Where("id NOT IN ?", v)
			}
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err = query.Session(&gorm.Session{PrepareStmt: true}).Offset(offset).Limit(limit).Find(&hanChars).Error
	if err != nil {
		return nil, 0, err
	}
	return hanChars, total, nil
}

// SearchWithFilters 搜索汉字（带过滤条件）
func (r *hanCharRepository) SearchWithFilters(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error) {
	var hanChars []*entity.HanChar
	var total int64

	query := r.DB.WithContext(ctx).
		Where("character ILIKE ? OR pinyin ILIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "level":
			if v, ok := value.(valueobject.WordDifficultyLevel); ok {
				query = query.Where("level = ?", v)
			}
		case "tags":
			if v, ok := value.([]string); ok && len(v) > 0 {
				query = query.Where("tags @> ?", v)
			}
		case "categories":
			if v, ok := value.([]string); ok && len(v) > 0 {
				query = query.Where("categories @> ?", v)
			}
		}
	}

	if err := query.Model(&entity.HanChar{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Session(&gorm.Session{PrepareStmt: true}).Offset(offset).Limit(limit).Find(&hanChars).Error
	return hanChars, total, err
}

var _ repository.HanCharRepository = (*hanCharRepository)(nil)
