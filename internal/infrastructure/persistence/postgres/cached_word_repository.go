package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/lazyjean/sla2/internal/infrastructure/cache"
)

type CachedWordRepository struct {
	repo  repository.WordRepository
	cache cache.Cache
}

func NewCachedWordRepository(repo repository.WordRepository, cache cache.Cache) repository.CachedWordRepository {
	return &CachedWordRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedWordRepository) Create(ctx context.Context, word *entity.Word) error {
	if err := r.repo.Create(ctx, word); err != nil {
		return err
	}

	// 写入缓存
	if wordJSON, err := json.Marshal(word); err == nil {
		cacheKey := fmt.Sprintf("word:%d", word.ID)
		r.cache.Set(ctx, cacheKey, string(wordJSON), 30*time.Minute)
	}

	return nil
}

func (r *CachedWordRepository) Update(ctx context.Context, word *entity.Word) error {
	return r.repo.Update(ctx, word)
}

func (r *CachedWordRepository) Delete(ctx context.Context, id entity.WordID) error {
	return r.repo.Delete(ctx, id)
}

func (r *CachedWordRepository) GetByID(ctx context.Context, id entity.WordID) (*entity.Word, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *CachedWordRepository) GetByWord(ctx context.Context, word string) (*entity.Word, error) {
	return r.repo.GetByWord(ctx, word)
}

func (r *CachedWordRepository) List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	return r.repo.List(ctx, offset, limit, filters)
}

func (r *CachedWordRepository) Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error) {
	return r.repo.Search(ctx, keyword, offset, limit, filters)
}

func (r *CachedWordRepository) GetAllTags(ctx context.Context) ([]string, error) {
	return r.repo.GetAllTags(ctx)
}

func (r *CachedWordRepository) GetAllCategories(ctx context.Context) ([]string, error) {
	return r.repo.GetAllCategories(ctx)
}

// FindByText 根据文本查找单词
func (r *CachedWordRepository) FindByText(ctx context.Context, text string) (*entity.Word, error) {
	filters := make(map[string]interface{})
	filters["keyword"] = text

	words, _, err := r.Search(ctx, text, 0, 1, filters)
	if err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, errors.ErrWordNotFound
	}
	return words[0], nil
}

func (r *CachedWordRepository) FindByID(ctx context.Context, id uint) (*entity.Word, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf("word:%d", id)
	if cached, err := r.cache.Get(ctx, cacheKey); err == nil {
		var word entity.Word
		if err := json.Unmarshal([]byte(cached), &word); err == nil {
			return &word, nil
		}
	}

	// 从数据库获取
	word, err := r.repo.GetByID(ctx, entity.WordID(id))
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if wordJSON, err := json.Marshal(word); err == nil {
		r.cache.Set(ctx, cacheKey, string(wordJSON), 30*time.Minute)
	}

	return word, nil
}

// ListByUserID 根据用户ID获取单词列表
func (r *CachedWordRepository) ListByUserID(ctx context.Context, userID uint, offset, limit int) ([]*entity.Word, int64, error) {
	// 从缓存获取
	cacheKey := fmt.Sprintf("words:list:%d:%d:%d", userID, offset, limit)
	var words []*entity.Word
	var total int64

	// 尝试从缓存获取
	if cachedData, err := r.cache.Get(ctx, cacheKey); err == nil {
		if err := json.Unmarshal([]byte(cachedData), &words); err == nil {
			return words, total, nil
		}
	}

	filters := make(map[string]interface{})
	filters["user_id"] = userID

	// 从数据库获取
	words, total, err := r.List(ctx, offset, limit, filters)
	if err != nil {
		return nil, 0, err
	}

	// 设置缓存
	if wordsJSON, err := json.Marshal(words); err == nil {
		r.cache.Set(ctx, cacheKey, string(wordsJSON), 30*time.Minute)
	}

	return words, total, nil
}

var _ repository.CachedWordRepository = (*CachedWordRepository)(nil)
