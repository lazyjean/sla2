package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	domainErrors "github.com/lazyjean/sla2/domain/errors"
	"gorm.io/gorm"
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

	// 设置时间戳
	now := time.Now()
	word.CreatedAt = now
	word.UpdatedAt = now

	// 保存主记录
	query := `
		INSERT INTO words (
			text, phonetic, translation, difficulty, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := tx.Raw(query,
		word.Text,
		word.Phonetic,
		word.Translation,
		word.Difficulty,
		word.CreatedAt,
		word.UpdatedAt,
	).Scan(&word.ID).Error

	if err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToSave
	}

	// 保存例句
	if len(word.Examples) > 0 {
		for i := range word.Examples {
			word.Examples[i].WordID = word.ID
			word.Examples[i].CreatedAt = now
		}
		if err := tx.Create(&word.Examples).Error; err != nil {
			tx.Rollback()
			return domainErrors.ErrFailedToSave
		}
	}

	// 保存标签
	if len(word.Tags) > 0 {
		for _, tag := range word.Tags {
			// 尝试查找或创建标签
			var existingTag entity.Tag
			err := tx.Where("name = ?", tag.Name).FirstOrCreate(&existingTag, &entity.Tag{
				Name:      tag.Name,
				CreatedAt: now,
			}).Error
			if err != nil {
				tx.Rollback()
				return domainErrors.ErrFailedToSave
			}

			// 关联标签和单词
			err = tx.Exec(`
				INSERT INTO word_tags (word_id, tag_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING`,
				word.ID, existingTag.ID,
			).Error
			if err != nil {
				tx.Rollback()
				return domainErrors.ErrFailedToSave
			}
		}
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

func (r *WordRepository) List(ctx context.Context, offset, limit int) ([]*entity.Word, error) {
	var words []*entity.Word

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
		return nil, domainErrors.ErrFailedToQuery
	}

	return words, nil
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
	result := r.db.WithContext(ctx).Delete(&entity.Word{}, "id = ?", id)
	if result.Error != nil {
		return domainErrors.ErrFailedToDelete
	}
	if result.RowsAffected == 0 {
		return domainErrors.ErrWordNotFound
	}
	return nil
}

// Update 实现更新方法
func (r *WordRepository) Update(ctx context.Context, word *entity.Word) error {
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return domainErrors.ErrFailedToUpdate
	}

	// 更新主记录
	result := tx.Model(&entity.Word{}).Where("id = ?", word.ID).Updates(map[string]interface{}{
		"text":        word.Text,
		"phonetic":    word.Phonetic,
		"translation": word.Translation,
		"difficulty":  word.Difficulty,
		"updated_at":  time.Now(),
	})

	if result.Error != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		return domainErrors.ErrWordNotFound
	}

	// 更新例句
	if err := tx.Model(&word).Association("Examples").Replace(word.Examples); err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	// 更新标签
	if err := tx.Model(&word).Association("Tags").Replace(word.Tags); err != nil {
		tx.Rollback()
		return domainErrors.ErrFailedToUpdate
	}

	if err := tx.Commit().Error; err != nil {
		return domainErrors.ErrFailedToUpdate
	}

	return nil
}
