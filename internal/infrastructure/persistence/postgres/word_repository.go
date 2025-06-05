package postgres

import (
	"context"
	"errors"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type wordRepository struct {
	*repository.GenericRepositoryImpl[*entity.Word, entity.WordID]
}

func NewVocabularyRepository(db *gorm.DB) repository.WordRepository {
	return &wordRepository{
		GenericRepositoryImpl: repository.NewGenericRepository[*entity.Word, entity.WordID](db),
	}
}

// ListWithFilters 获取单词列表（带过滤条件）
func (r *wordRepository) ListWithFilters(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	var words []*entity.Word
	var total int64

	query := r.DB.WithContext(ctx).Model(&entity.Word{})

	// 应用过滤条件
	for key, value := range filters {
		switch key {
		case "tags":
			query = query.Where("tags @> ?", value)
		case "categories":
			query = query.Where("categories @> ?", value)
		case "level":
			// 将字符串转换为对应的数字值
			switch value {
			case "HSK1":
				query = query.Where("difficulty = ?", 1)
			case "HSK2":
				query = query.Where("difficulty = ?", 2)
			case "HSK3":
				query = query.Where("difficulty = ?", 3)
			case "HSK4":
				query = query.Where("difficulty = ?", 4)
			case "HSK5":
				query = query.Where("difficulty = ?", 5)
			case "HSK6":
				query = query.Where("difficulty = ?", 6)
			default:
				// 如果是不支持的级别，不添加过滤条件
			}
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Offset(offset).Limit(limit).Find(&words).Error; err != nil {
		return nil, 0, err
	}

	return words, total, nil
}

// GetByWord 根据单词文本获取单词
func (r *wordRepository) GetByWord(ctx context.Context, text string) (*entity.Word, error) {
	var word entity.Word
	err := r.DB.WithContext(ctx).
		Where("text = ?", text).
		First(&word).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &word, nil
}

// GetAllCategories 获取所有分类
func (r *wordRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.DB.WithContext(ctx).Model(&entity.Word{}).
		Distinct().
		Pluck("unnest(categories)", &categories).
		Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAllTags 获取所有标签
func (r *wordRepository) GetAllTags(ctx context.Context) ([]string, error) {
	var tags []string
	err := r.DB.WithContext(ctx).Model(&entity.Word{}).
		Distinct().
		Pluck("unnest(tags)", &tags).
		Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// ListByIDs 通过ID列表获取单词
func (r *wordRepository) ListByIDs(ctx context.Context, ids []entity.WordID) ([]*entity.Word, error) {
	var words []*entity.Word
	if len(ids) == 0 {
		return words, nil
	}

	query := r.DB.WithContext(ctx).Model(&entity.Word{})
	if err := query.Where("id IN ?", ids).Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}

// Search 搜索单词
func (r *wordRepository) Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	db := r.DB.WithContext(ctx).Model(&entity.Word{})

	// 构建查询条件
	if keyword != "" {
		db = db.Where("text ILIKE ?", "%"+keyword+"%")
	}

	if tags, ok := filters["tags"].([]string); ok && len(tags) > 0 {
		db = db.Where("tags @> ?", tags)
	}

	if level, ok := filters["level"].(string); ok && level != "" {
		db = db.Where("difficulty = ?", level)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	// 排序
	db = db.Order("created_at DESC")

	// 分页
	db = db.Offset(offset).Limit(limit)

	// 执行查询
	var words []*entity.Word
	if err := db.Find(&words).Error; err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	return words, total, nil
}

var _ repository.WordRepository = (*wordRepository)(nil)
