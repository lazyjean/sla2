package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	"github.com/lazyjean/sla2/domain/repository"
	"github.com/lazyjean/sla2/infrastructure/cache"
)

type CachedWordRepository struct {
	repo  repository.WordRepository
	cache cache.Cache
}

func NewCachedWordRepository(repo repository.WordRepository, cache cache.Cache) *CachedWordRepository {
	return &CachedWordRepository{
		repo:  repo,
		cache: cache,
	}
}

func (r *CachedWordRepository) Save(ctx context.Context, word *entity.Word) error {
	if err := r.repo.Save(ctx, word); err != nil {
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
	if err := r.repo.Update(ctx, word); err != nil {
		return err
	}

	// 更新缓存
	if wordJSON, err := json.Marshal(word); err == nil {
		cacheKey := fmt.Sprintf("word:%d", word.ID)
		r.cache.Set(ctx, cacheKey, string(wordJSON), 30*time.Minute)
	}

	return nil
}

func (r *CachedWordRepository) Delete(ctx context.Context, id uint) error {
	if err := r.repo.Delete(ctx, id); err != nil {
		return err
	}

	// 删除缓存
	cacheKey := fmt.Sprintf("word:%d", id)
	r.cache.Delete(ctx, cacheKey)

	return nil
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
	word, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if wordJSON, err := json.Marshal(word); err == nil {
		r.cache.Set(ctx, cacheKey, string(wordJSON), 30*time.Minute)
	}

	return word, nil
}

func (r *CachedWordRepository) FindByText(ctx context.Context, text string) (*entity.Word, error) {
	// 这个方法不使用缓存，因为是按文本搜索
	return r.repo.FindByText(ctx, text)
}

func (r *CachedWordRepository) List(ctx context.Context, offset, limit int) ([]*entity.Word, error) {
	// 列表查询暂不使用缓存，因为可能会频繁变化
	return r.repo.List(ctx, offset, limit)
}

// 其他方法类似，实现缓存逻辑...
