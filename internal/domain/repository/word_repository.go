package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// WordRepository 单词仓储接口
type WordRepository interface {
	Save(ctx context.Context, word *entity.Word) error
	Update(ctx context.Context, word *entity.Word) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint) (*entity.Word, error)
	FindByText(ctx context.Context, text string) (*entity.Word, error)
	List(ctx context.Context, userID uint, offset, limit int) ([]*entity.Word, int64, error)
	Search(ctx context.Context, query *WordQuery) ([]*entity.Word, int64, error)
}

type CachedWordRepository interface {
	WordRepository
}

// WordQuery 定义查询参数
type WordQuery struct {
	UserID        entity.UserID // 用户ID
	Text          string        // 单词文本
	Tags          []string      // 标签列表
	MinDifficulty int           // 最小难度
	MaxDifficulty int           // 最大难度
	MasteryLevel  *int          // 掌握程度
	ReviewBefore  time.Time     // 需要在此时间前复习
	CreatedAfter  time.Time     // 在此时间后创建
	OrderBy       string        // 排序字段
	OrderDesc     bool          // 是否降序
	Offset        int           // 分页偏移
	Limit         int           // 分页大小
}
