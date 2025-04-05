package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/domain/valueobject"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// hanCharRepository PostgreSQL 汉字仓库实现
type hanCharRepository struct {
	db *gorm.DB
}

// NewHanCharRepository 创建汉字仓库实例
func NewHanCharRepository(db *gorm.DB) repository.HanCharRepository {
	return &hanCharRepository{
		db: db,
	}
}

// Create 创建汉字
func (r *hanCharRepository) Create(ctx context.Context, hanChar *entity.HanChar) (entity.HanCharID, error) {
	// 使用 GORM 的 Create 方法，GORM 会自动处理 JSONB 字段的序列化
	err := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "character"}},
			DoNothing: true,
		}).
		Create(hanChar).Error
	if err != nil {
		return 0, err
	}
	return hanChar.ID, nil
}

// Update 更新汉字
func (r *hanCharRepository) Update(ctx context.Context, hanChar *entity.HanChar) error {
	return r.db.WithContext(ctx).Save(hanChar).Error
}

// Delete 删除汉字
func (r *hanCharRepository) Delete(ctx context.Context, id entity.HanCharID) error {
	return r.db.WithContext(ctx).Delete(&entity.HanChar{}, id).Error
}

// GetByID 根据ID获取汉字
func (r *hanCharRepository) GetByID(ctx context.Context, id entity.HanCharID) (*entity.HanChar, error) {
	var hanChar entity.HanChar
	err := r.db.WithContext(ctx).First(&hanChar, id).Error
	if err != nil {
		return nil, err
	}
	return &hanChar, nil
}

// GetByCharacter 根据字符获取汉字
func (r *hanCharRepository) GetByCharacter(ctx context.Context, character string) (*entity.HanChar, error) {
	var hanChar entity.HanChar
	err := r.db.WithContext(ctx).Where("character = ?", character).First(&hanChar).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &hanChar, nil
}

// List 获取汉字列表
func (r *hanCharRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error) {
	var hanChars []*entity.HanChar
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.HanChar{})

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

// Search 搜索汉字
func (r *hanCharRepository) Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.HanChar, int64, error) {
	var hanChars []*entity.HanChar
	var total int64

	query := r.db.WithContext(ctx).
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
