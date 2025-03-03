package utils

import (
	"context"
	"errors"
	"strconv"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

const (
	userIDKey = "user_id"
)

// GetUserIDFromContext 从上下文中获取用户ID
func GetUserIDFromContext(ctx context.Context) (entity.UID, error) {
	userID, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return 0, errors.New("user id not found in context")
	}

	// 将字符串转换为uint
	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id format")
	}

	return entity.UID(id), nil
}

// SetUserIDToContext 将用户ID设置到上下文中
func SetUserIDToContext(ctx context.Context, userID entity.UID) context.Context {
	return context.WithValue(ctx, userIDKey, strconv.FormatUint(uint64(userID), 10))
}
