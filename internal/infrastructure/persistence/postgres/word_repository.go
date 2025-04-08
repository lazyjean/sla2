package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) repository.WordRepository {
	return &WordRepository{db: db}
}

// List 获取单词列表
func (r *WordRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	var words []*entity.Word
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Word{})

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

// Create 创建单词
func (r *WordRepository) Create(ctx context.Context, word *entity.Word) error {
	// 先验证数据有效性
	if err := word.Validate(); err != nil {
		return err
	}

	// 检查单词是否已存在
	existing, err := r.GetByWord(ctx, word.Text)
	if err != nil && err != domainErrors.ErrWordNotFound {
		return err
	}
	if existing != nil {
		return domainErrors.ErrWordAlreadyExists
	}

	// 只插入必要的字段
	if err := r.db.WithContext(ctx).Create(word).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// Save 保存单词
func (r *WordRepository) Save(ctx context.Context, word *entity.Word) error {
	// 先验证数据有效性
	if err := word.Validate(); err != nil {
		return err
	}

	// 检查单词是否已存在
	existing, err := r.GetByWord(ctx, word.Text)
	if err != nil {
		return err
	}
	if existing != nil {
		return domainErrors.ErrWordAlreadyExists
	}

	// 只插入必要的字段
	if err := r.db.WithContext(ctx).Select(
		"Text",
		"Phonetic",
		"Definitions",
		"Examples",
		"Tags",
		"Difficulty",
		"CreatedAt",
		"UpdatedAt",
	).Create(word).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// FindByID 通过ID查找单词
func (r *WordRepository) FindByID(ctx context.Context, id uint) (*entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(ctx).
		First(&word, id).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &word, nil
}

// FindByUserAndText 通过用户ID和文本查找单词
func (r *WordRepository) FindByUserAndText(ctx context.Context, userID uint, text string) (*entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND text = ?", userID, text).
		First(&word).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &word, nil
}

func (r *WordRepository) ListByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.Word, int64, error) {
	var words []*entity.Word
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.Word{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	// 获取列表
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	return words, total, nil
}

// ListNeedReview 获取需要复习的单词列表
func (r *WordRepository) ListNeedReview(ctx context.Context, before time.Time, limit int) ([]*entity.Word, error) {
	var words []*entity.Word
	err := r.db.WithContext(ctx).
		Where("next_review_at <= ? AND mastery_level < ?", before, entity.MasteryLevelMastered).
		Order("next_review_at ASC").
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return words, nil
}

// Search 搜索单词
func (r *WordRepository) Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.Word{})

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

// WordQuery 定义查询参数
type WordQuery struct {
	Text          string
	Tags          []string
	MinDifficulty int
	MaxDifficulty int
	Offset        int
	Limit         int
}

// SearchWords 高级查询
func (r *WordRepository) SearchWords(ctx context.Context, query *WordQuery) ([]*entity.Word, error) {
	db := r.db.WithContext(ctx).Model(&entity.Word{})

	if query.Text != "" {
		db = db.Where("text ILIKE ?", "%"+query.Text+"%")
	}
	if len(query.Tags) > 0 {
		db = db.Joins("JOIN word_tags ON words.id = word_tags.word_id").
			Joins("JOIN tags ON word_tags.tag_id = tags.id").
			Where("tags.name IN ?", query.Tags)
	}
	if query.MinDifficulty > 0 {
		db = db.Where("difficulty >= ?", query.MinDifficulty)
	}
	if query.MaxDifficulty > 0 {
		db = db.Where("difficulty <= ?", query.MaxDifficulty)
	}

	var words []*entity.Word
	result := db.
		Preload("Examples").
		Preload("Tags").
		Offset(query.Offset).
		Limit(query.Limit).
		Find(&words)

	if result.Error != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return words, nil
}

// Delete 删除单词
func (r *WordRepository) Delete(ctx context.Context, id entity.WordID) error {
	result := r.db.WithContext(ctx).Delete(&entity.Word{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domainErrors.ErrWordNotFound
	}
	return nil
}

// Update 更新单词
func (r *WordRepository) Update(ctx context.Context, word *entity.Word) error {
	// 先验证数据有效性
	if err := word.Validate(); err != nil {
		return err
	}

	// 检查单词是否存在
	existing, err := r.GetByID(ctx, word.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return domainErrors.ErrWordNotFound
	}

	// 使用 Save 方法，GORM 会自动处理 JSON 类型的序列化
	if err := r.db.WithContext(ctx).Save(word).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// FindByText 通过文本查找单词
func (r *WordRepository) FindByText(ctx context.Context, text string) (*entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(ctx).
		Where("text = ?", text).
		First(&word).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &word, nil
}

// GetByWord 根据单词文本获取单词
func (r *WordRepository) GetByWord(ctx context.Context, text string) (*entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(ctx).
		Where("text = ?", text).
		First(&word).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &word, nil
}

// GetAllCategories 获取所有分类
func (r *WordRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).Model(&entity.Word{}).
		Distinct().
		Pluck("unnest(categories)", &categories).
		Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

// GetAllTags 获取所有标签
func (r *WordRepository) GetAllTags(ctx context.Context) ([]string, error) {
	var tags []string
	err := r.db.WithContext(ctx).Model(&entity.Word{}).
		Distinct().
		Pluck("unnest(tags)", &tags).
		Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetByID 根据ID获取单词
func (r *WordRepository) GetByID(ctx context.Context, id entity.WordID) (*entity.Word, error) {
	var word entity.Word
	err := r.db.WithContext(ctx).First(&word, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainErrors.ErrWordNotFound
		}
		return nil, err
	}
	return &word, nil
}

// ListByIDs 通过ID列表获取单词
func (r *WordRepository) ListByIDs(ctx context.Context, ids []entity.WordID) ([]*entity.Word, error) {
	var words []*entity.Word
	if len(ids) == 0 {
		return words, nil
	}

	query := r.db.WithContext(ctx).Model(&entity.Word{})
	if err := query.Where("id IN ?", ids).Find(&words).Error; err != nil {
		return nil, err
	}

	return words, nil
}

var _ repository.WordRepository = (*WordRepository)(nil)
