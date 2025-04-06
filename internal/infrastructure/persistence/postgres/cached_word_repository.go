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

// ListByIDs 通过ID列表获取单词
func (r *CachedWordRepository) ListByIDs(ctx context.Context, ids []entity.WordID) ([]*entity.Word, error) {
	// 1. 从缓存中获取
	var words []*entity.Word
	var missingIDs []entity.WordID
	for _, id := range ids {
		if wordStr, err := r.cache.Get(ctx, fmt.Sprintf("word:%d", id)); err == nil {
			var word entity.Word
			if err := json.Unmarshal([]byte(wordStr), &word); err == nil {
				words = append(words, &word)
			} else {
				missingIDs = append(missingIDs, id)
			}
		} else {
			missingIDs = append(missingIDs, id)
		}
	}

	// 2. 如果所有单词都在缓存中，直接返回
	if len(missingIDs) == 0 {
		return words, nil
	}

	// 3. 从数据库中获取缺失的单词
	dbWords, err := r.repo.ListByIDs(ctx, missingIDs)
	if err != nil {
		return nil, err
	}

	// 4. 将数据库中的单词存入缓存
	for _, word := range dbWords {
		wordBytes, err := json.Marshal(word)
		if err != nil {
			return nil, err
		}
		if err := r.cache.Set(ctx, fmt.Sprintf("word:%d", word.ID), string(wordBytes), 24*time.Hour); err != nil {
			return nil, err
		}
		words = append(words, word)
	}

	return words, nil
}

// ListNeedReview 获取需要复习的单词
func (r *CachedWordRepository) ListNeedReview(ctx context.Context, before time.Time, limit int) ([]*entity.Word, error) {
	// 1. 从数据库中获取需要复习的单词
	words, err := r.repo.ListNeedReview(ctx, before, limit)
	if err != nil {
		return nil, err
	}

	// 2. 将单词存入缓存
	for _, word := range words {
		wordBytes, err := json.Marshal(word)
		if err != nil {
			return nil, err
		}
		if err := r.cache.Set(ctx, fmt.Sprintf("word:%d", word.ID), string(wordBytes), 24*time.Hour); err != nil {
			return nil, err
		}
	}

	return words, nil
}

var _ repository.CachedWordRepository = (*CachedWordRepository)(nil)
