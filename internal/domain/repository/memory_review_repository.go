package repository

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// MemoryReviewRepository 记忆复习记录仓储接口
type MemoryReviewRepository interface {
	// Create 创建复习记录
	Create(ctx context.Context, review *entity.MemoryReview) error
	// ListByUserIDAndTimeRange 获取用户指定时间范围内的复习记录
	ListByUserIDAndTimeRange(ctx context.Context, userID entity.UID, startTime, endTime time.Time) ([]*entity.MemoryReview, error)
}
