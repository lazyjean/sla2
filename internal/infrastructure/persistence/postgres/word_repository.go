package postgres

import (
	"context"
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
	if userID, ok := filters["user_id"].(uint); ok {
		query = query.Where("user_id = ?", userID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	// 获取列表
	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&words).Error

	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	return words, total, nil
}

// Create 创建单词
func (r *WordRepository) Create(ctx context.Context, word *entity.Word) error {
	return r.Save(ctx, word)
}

// Save 保存单词
func (r *WordRepository) Save(ctx context.Context, word *entity.Word) error {
	// 先验证数据有效性
	if err := word.Validate(); err != nil {
		return err
	}

	// 检查单词是否已存在
	existing, err := r.FindByUserAndText(ctx, uint(word.UserID), word.Text)
	if err != nil && err != domainErrors.ErrWordNotFound {
		return err
	}
	if existing != nil {
		return domainErrors.ErrWordAlreadyExists
	}

	// 使用 Create 方法，GORM 会自动处理 JSON 类型的序列化
	if err := r.db.WithContext(ctx).Create(word).Error; err != nil {
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

func (r *WordRepository) ListNeedReview(ctx context.Context, userID uint, before time.Time) ([]*entity.Word, error) {
	var words []*entity.Word
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND next_review_at <= ? AND mastery_level < 5", userID, before).
		Order("next_review_at ASC").
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
	if userID, ok := filters["user_id"].(uint); ok {
		db = db.Where("user_id = ?", userID)
	}

	if keyword != "" {
		db = db.Where("text ILIKE ?", "%"+keyword+"%")
	}

	if tags, ok := filters["tags"].([]string); ok && len(tags) > 0 {
		db = db.Joins("JOIN word_tags ON words.id = word_tags.word_id").
			Joins("JOIN tags ON tags.id = word_tags.tag_id").
			Where("tags.name IN ?", tags)
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

	// 加载关联数据
	db = db.Preload("Examples").Preload("Tags")

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

// Delete 实现删除方法
func (r *WordRepository) Delete(ctx context.Context, id entity.WordID) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domainErrors.ErrFailedToDelete
	}

	// 先检查记录是否存在
	exists, err := r.exists(ctx, uint(id))
	if err != nil {
		tx.Rollback()
		return err
	}
	if !exists {
		tx.Rollback()
		return domainErrors.ErrWordNotFound
	}

	// 删除关联的例句和标签关系（由于设置了 ON DELETE CASCADE，这步可以省略）
	// 执行真实删除
	if err := tx.Unscoped().Delete(&entity.Word{}, id).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToDelete
	}

	return tx.Commit().Error
}

func (r *WordRepository) exists(ctx context.Context, id uint) (bool, error) {
	var exists bool
	err := r.db.WithContext(ctx).
		Model(&entity.Word{}).
		Select("1").
		Where("id = ?", id).
		Scan(&exists).Error

	if err != nil {
		return false, domainErrors.ErrFailedToQuery
	}
	return exists, nil
}

// Update 实现更新方法
func (r *WordRepository) Update(ctx context.Context, word *entity.Word) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domainErrors.ErrFailedToUpdate
	}

	// 先检查记录是否存在
	var exists bool
	err := tx.Model(&entity.Word{}).
		Select("1").
		Where("id = ?", word.ID).
		Scan(&exists).Error
	if err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToQuery
	}
	if !exists {
		tx.Rollback()
		return domainErrors.ErrWordNotFound
	}

	// 更新主记录
	if err := tx.Model(&entity.Word{}).
		Where("id = ?", word.ID).
		Updates(map[string]interface{}{
			"text":        word.Text,
			"phonetic":    word.Phonetic,
			"definitions": word.Definitions,
			"difficulty":  word.Difficulty,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	return tx.Commit().Error
}

// FindByText implements repository.WordRepository.
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

// GetAllCategories 获取所有分类
func (r *WordRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := r.db.WithContext(ctx).
		Model(&entity.Word{}).
		Distinct().
		Pluck("categories", &categories).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return categories, nil
}

// GetAllTags 获取所有标签
func (r *WordRepository) GetAllTags(ctx context.Context) ([]string, error) {
	var tags []string
	err := r.db.WithContext(ctx).
		Model(&entity.Word{}).
		Distinct().
		Pluck("tags", &tags).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return tags, nil
}

// GetByID 根据ID获取单词
func (r *WordRepository) GetByID(ctx context.Context, id entity.WordID) (*entity.Word, error) {
	return r.FindByID(ctx, uint(id))
}

// GetByWord 根据单词获取
func (r *WordRepository) GetByWord(ctx context.Context, word string) (*entity.Word, error) {
	return r.FindByText(ctx, word)
}

var _ repository.WordRepository = (*WordRepository)(nil)
