package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	domainErrors "github.com/lazyjean/sla2/domain/errors"
	"github.com/lazyjean/sla2/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db: db}
}

func (r *WordRepository) Save(ctx context.Context, word *entity.Word) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domainErrors.ErrFailedToSave
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 设置时间戳
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	// 先处理所有的 tags
	for i, tag := range word.Tags {
		var existingTag entity.Tag
		err := tx.Where("name = ?", tag.Name).First(&existingTag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 标签不存在，创建新标签
				existingTag = entity.Tag{
					Name:      tag.Name,
					CreatedAt: now,
				}
				if err := tx.Create(&existingTag).Error; err != nil {
					tx.Rollback()
					return domainErrors.ErrFailedToSave
				}
			} else {
				tx.Rollback()
				return domainErrors.ErrFailedToSave
			}
		}
		word.Tags[i] = existingTag
	}

	// 保存主记录
	if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(word).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToSave
	}

	return tx.Commit().Error
}

func (r *WordRepository) FindByID(ctx context.Context, id uint) (*entity.Word, error) {
	var word entity.Word

	err := r.db.WithContext(ctx).
		Preload("Examples").
		Preload("Tags").
		Unscoped().
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&word).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}

	return &word, nil
}

func (r *WordRepository) FindByText(ctx context.Context, text string) (*entity.Word, error) {
	var word entity.Word

	err := r.db.WithContext(ctx).
		Preload("Examples").
		Preload("Tags").
		Unscoped().
		Where("text = ? AND deleted_at IS NULL", text).
		Take(&word).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domainErrors.ErrWordNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}

	return &word, nil
}

func (r *WordRepository) List(ctx context.Context, offset, limit int) ([]*entity.Word, int64, error) {
	var words []*entity.Word
	var total int64

	// 获取总数
	if err := r.db.WithContext(ctx).Model(&entity.Word{}).
		Where("deleted_at IS NULL").
		Count(&total).Error; err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	// 获取分页数据
	err := r.db.WithContext(ctx).
		Preload("Examples").
		Preload("Tags").
		Unscoped().
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&words).Error

	if err != nil {
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
func (r *WordRepository) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domainErrors.ErrFailedToDelete
	}

	// 先检查记录是否存在
	exists, err := r.exists(ctx, id)
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
			"translation": word.Translation,
			"difficulty":  word.Difficulty,
			"updated_at":  time.Now(),
		}).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	// 删除旧的例句
	if err := tx.Where("word_id = ?", word.ID).Delete(&entity.Example{}).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	// 添加新的例句
	if len(word.Examples) > 0 {
		for i := range word.Examples {
			word.Examples[i].WordID = word.ID
			word.Examples[i].CreatedAt = time.Now()
		}
		if err := tx.Create(&word.Examples).Error; err != nil {
			tx.Rollback()
			return domainErrors.ErrFailedToUpdate
		}
	}

	// 更新标签关联
	if err := tx.Where("word_id = ?", word.ID).Delete(&entity.WordTag{}).Error; err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	// 处理标签
	for _, tag := range word.Tags {
		var existingTag entity.Tag
		err := tx.Where("name = ?", tag.Name).First(&existingTag).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 标签不存在，创建新标签
				existingTag = entity.Tag{
					Name:      tag.Name,
					CreatedAt: time.Now(),
				}
				if err := tx.Create(&existingTag).Error; err != nil {
					tx.Rollback()
					return domainErrors.ErrFailedToUpdate
				}
			} else {
				tx.Rollback()
				return domainErrors.ErrFailedToUpdate
			}
		}

		// 创建标签关联
		if err := tx.Create(&entity.WordTag{
			WordID: word.ID,
			TagID:  existingTag.ID,
		}).Error; err != nil {
			tx.Rollback()
			return domainErrors.ErrFailedToUpdate
		}
	}

	return tx.Commit().Error
}

func (r *WordRepository) Search(ctx context.Context, query *repository.WordQuery) ([]*entity.Word, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.Word{})

	// 添加文本搜索条件
	if query.Text != "" {
		db = db.Where("text ILIKE ?", "%"+query.Text+"%")
	}

	// 添加标签过滤
	if len(query.Tags) > 0 {
		db = db.Joins("JOIN word_tags ON words.id = word_tags.word_id").
			Joins("JOIN tags ON word_tags.tag_id = tags.id").
			Where("tags.name IN ?", query.Tags)
	}

	// 添加难度范围过滤
	if query.MinDifficulty > 0 {
		db = db.Where("difficulty >= ?", query.MinDifficulty)
	}
	if query.MaxDifficulty > 0 {
		db = db.Where("difficulty <= ?", query.MaxDifficulty)
	}

	// 获取总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	var words []*entity.Word
	err := db.Preload("Examples").
		Preload("Tags").
		Offset(query.Offset).
		Limit(query.Limit).
		Order("created_at DESC").
		Find(&words).Error

	if err != nil {
		return nil, 0, err
	}

	return words, total, nil
}
