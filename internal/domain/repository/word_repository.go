package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// WordRepository 单词仓库接口
type WordRepository interface {
	// Create 创建单词
	Create(ctx context.Context, word *entity.Word) error
	// Update 更新单词
	Update(ctx context.Context, word *entity.Word) error
	// Delete 删除单词
	Delete(ctx context.Context, id entity.WordID) error
	// GetByID 根据ID获取单词
	GetByID(ctx context.Context, id entity.WordID) (*entity.Word, error)
	// GetByWord 根据单词获取
	GetByWord(ctx context.Context, word string) (*entity.Word, error)
	// List 获取单词列表
	List(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error)
	// Search 搜索单词
	Search(ctx context.Context, keyword string, offset, limit int, filters map[string]interface{}) ([]*entity.Word, int64, error)
	// GetAllTags 获取所有标签
	GetAllTags(ctx context.Context) ([]string, error)
	// GetAllCategories 获取所有分类
	GetAllCategories(ctx context.Context) ([]string, error)
}

type CachedWordRepository interface {
	WordRepository
}

// WordQuery 定义查询参数
type WordQuery struct {
	UserID        entity.UID // 用户ID
	Keyword       string     // 搜索关键词
	Tags          []string   // 标签列表
	Categories    []string   // 分类列表
	MinDifficulty int        // 最小难度
	MaxDifficulty int        // 最大难度
	MasteryLevel  *int       // 掌握程度
	ReviewBefore  time.Time  // 需要在此时间前复习
	CreatedAfter  time.Time  // 在此时间后创建
	OrderBy       string     // 排序字段
	OrderDesc     bool       // 是否降序
	Offset        int        // 分页偏移
	Limit         int        // 分页大小
}

// ToFilters 转换为过滤条件
func (q *WordQuery) ToFilters() map[string]interface{} {
	filters := make(map[string]interface{})
	if q.UserID > 0 {
		filters["user_id"] = q.UserID
	}
	if len(q.Tags) > 0 {
		filters["tags"] = q.Tags
	}
	if len(q.Categories) > 0 {
		filters["categories"] = q.Categories
	}
	if q.MinDifficulty > 0 {
		filters["min_difficulty"] = q.MinDifficulty
	}
	if q.MaxDifficulty > 0 {
		filters["max_difficulty"] = q.MaxDifficulty
	}
	if q.MasteryLevel != nil {
		filters["mastery_level"] = *q.MasteryLevel
	}
	if !q.ReviewBefore.IsZero() {
		filters["review_before"] = q.ReviewBefore
	}
	if !q.CreatedAfter.IsZero() {
		filters["created_after"] = q.CreatedAfter
	}
	return filters
}
