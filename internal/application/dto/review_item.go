package dto

import "time"

// ReviewItem 表示一个复习项
type ReviewItem struct {
	MemoryUnitID   uint32
	Result         bool // true 表示正确，false 表示错误
	ReviewDuration time.Duration
}
