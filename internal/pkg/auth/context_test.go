package auth

import (
	"context"
	"testing"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name        string
		ctxSetup    func() context.Context
		expected    entity.UID
		expectError bool
	}{
		{
			name: "valid_uid",
			ctxSetup: func() context.Context {
				return context.WithValue(context.Background(), ContextKeyUserID, entity.UID(123))
			},
			expected:    entity.UID(123),
			expectError: false,
		},
		{
			name: "valid_string",
			ctxSetup: func() context.Context {
				return context.WithValue(context.Background(), ContextKeyUserID, "456")
			},
			expected:    entity.UID(456),
			expectError: false,
		},
		{
			name: "no_user_id",
			ctxSetup: func() context.Context {
				return context.Background()
			},
			expected:    entity.UID(0),
			expectError: true,
		},
		{
			name: "zero_user_id",
			ctxSetup: func() context.Context {
				return context.WithValue(context.Background(), ContextKeyUserID, entity.UID(0))
			},
			expected:    entity.UID(0),
			expectError: true,
		},
		{
			name: "invalid_string",
			ctxSetup: func() context.Context {
				return context.WithValue(context.Background(), ContextKeyUserID, "not-a-number")
			},
			expected:    entity.UID(0),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctxSetup()
			uid, err := GetUserID(ctx)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, uid)
			}
		})
	}
}

func TestGetUserIDString(t *testing.T) {
	ctx := context.WithValue(context.Background(), ContextKeyUserID, entity.UID(123))

	idStr, err := GetUserIDString(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "123", idStr)

	// 测试错误情况
	emptyCtx := context.Background()
	_, err = GetUserIDString(emptyCtx)
	assert.Error(t, err)
}

func TestGetUserRoles(t *testing.T) {
	roles := []string{"admin", "user"}
	ctx := context.WithValue(context.Background(), ContextKeyRoles, roles)

	fetchedRoles, err := GetUserRoles(ctx)
	assert.NoError(t, err)
	assert.Equal(t, roles, fetchedRoles)

	// 测试错误情况
	emptyCtx := context.Background()
	_, err = GetUserRoles(emptyCtx)
	assert.Error(t, err)
}

func TestHasRole(t *testing.T) {
	roles := []string{"admin", "user"}
	ctx := context.WithValue(context.Background(), ContextKeyRoles, roles)

	assert.True(t, HasRole(ctx, "admin"))
	assert.True(t, HasRole(ctx, "user"))
	assert.False(t, HasRole(ctx, "guest"))

	// 测试没有角色的情况
	emptyCtx := context.Background()
	assert.False(t, HasRole(emptyCtx, "admin"))
}

func TestHasAnyRole(t *testing.T) {
	roles := []string{"admin", "user"}
	ctx := context.WithValue(context.Background(), ContextKeyRoles, roles)

	assert.True(t, HasAnyRole(ctx, "admin", "guest"))
	assert.True(t, HasAnyRole(ctx, "user", "guest"))
	assert.False(t, HasAnyRole(ctx, "guest", "viewer"))

	// 测试没有角色的情况
	emptyCtx := context.Background()
	assert.False(t, HasAnyRole(emptyCtx, "admin", "user"))
}

func TestHasAllRoles(t *testing.T) {
	roles := []string{"admin", "user", "editor"}
	ctx := context.WithValue(context.Background(), ContextKeyRoles, roles)

	assert.True(t, HasAllRoles(ctx, "admin", "user"))
	assert.True(t, HasAllRoles(ctx, "admin", "editor"))
	assert.False(t, HasAllRoles(ctx, "admin", "guest"))

	// 测试没有角色的情况
	emptyCtx := context.Background()
	assert.False(t, HasAllRoles(emptyCtx, "admin", "user"))
}

func TestWithUserID(t *testing.T) {
	ctx := context.Background()
	uid := entity.UID(123)

	newCtx := WithUserID(ctx, uid)

	fetchedUID, err := GetUserID(newCtx)
	assert.NoError(t, err)
	assert.Equal(t, uid, fetchedUID)
}

func TestWithRoles(t *testing.T) {
	ctx := context.Background()
	roles := []string{"admin", "user"}

	newCtx := WithRoles(ctx, roles)

	fetchedRoles, err := GetUserRoles(newCtx)
	assert.NoError(t, err)
	assert.Equal(t, roles, fetchedRoles)
}
