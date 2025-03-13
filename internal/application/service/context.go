package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// contextKey 是用于context的键类型
type contextKey string

const (
	// ContextKeyUserID 用户ID在context中的键
	ContextKeyUserID contextKey = "user_id"
	// ContextKeyRoles 用户角色在context中的键
	ContextKeyRoles contextKey = "user_roles"
)

var (
	// ErrNoUserInContext 当context中没有用户信息时返回此错误
	ErrNoUserInContext = errors.New("context中没有用户信息")
	// ErrInvalidUserIDType 当用户ID类型不正确时返回此错误
	ErrInvalidUserIDType = errors.New("context中的用户ID类型不正确")
	// ErrInvalidRolesType 当角色类型不正确时返回此错误
	ErrInvalidRolesType = errors.New("context中的角色类型不正确")
)

// GetUserID 从context中获取用户ID
func GetUserID(ctx context.Context) (entity.UID, error) {
	userID, exists := ctx.Value(ContextKeyUserID).(entity.UID)
	if !exists {
		// 尝试获取字符串类型并转换
		if userIDStr, ok := ctx.Value(ContextKeyUserID).(string); ok {
			// 直接使用 strconv 进行转换
			uid, err := strconv.ParseUint(userIDStr, 10, 64)
			if err != nil {
				return 0, fmt.Errorf("转换用户ID失败: %w", err)
			}
			return entity.UID(uid), nil
		}

		// 如果都获取不到，返回错误
		return 0, ErrNoUserInContext
	}

	// 如果用户ID为0，也视为未登录
	if userID == 0 {
		return 0, ErrNoUserInContext
	}

	return userID, nil
}

// GetUserRoles 从context中获取用户角色
func GetUserRoles(ctx context.Context) ([]string, error) {
	roles, ok := ctx.Value(ContextKeyRoles).([]string)
	if !ok {
		return nil, ErrInvalidRolesType
	}
	return roles, nil
}

// HasRole 检查用户是否拥有指定角色
func HasRole(ctx context.Context, role string) bool {
	roles, err := GetUserRoles(ctx)
	if err != nil {
		return false
	}

	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

// HasAnyRole 检查用户是否拥有任意一个指定角色
func HasAnyRole(ctx context.Context, roles ...string) bool {
	userRoles, err := GetUserRoles(ctx)
	if err != nil {
		return false
	}

	userRolesMap := make(map[string]struct{}, len(userRoles))
	for _, r := range userRoles {
		userRolesMap[r] = struct{}{}
	}

	for _, r := range roles {
		if _, exists := userRolesMap[r]; exists {
			return true
		}
	}

	return false
}

// WithUserID 将用户ID添加到context中
func WithUserID(ctx context.Context, userID entity.UID) context.Context {
	return context.WithValue(ctx, ContextKeyUserID, userID)
}

// WithRoles 将用户角色添加到context中
func WithRoles(ctx context.Context, roles []string) context.Context {
	return context.WithValue(ctx, ContextKeyRoles, roles)
}
