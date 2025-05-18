package test

import (
	"context"
	"net/http"

	"github.com/lazyjean/sla2/internal/domain/entity"
)

// MockTokenService 是一个用于测试的 token service
type MockTokenService struct{}

// NewMockTokenService 创建一个新的 mock token service
func NewMockTokenService() *MockTokenService {
	return &MockTokenService{}
}

// GenerateToken 生成一个固定的测试 token
func (s *MockTokenService) GenerateToken(userID entity.UID, roles []string) (string, error) {
	return "test_token", nil
}

// ValidateToken 验证 token，总是返回固定的用户 ID 和角色
func (s *MockTokenService) ValidateToken(tokenString string) (entity.UID, []string, error) {
	return 1, []string{"user"}, nil
}

// ValidateTokenFromContext 从上下文中验证 token，总是返回固定的用户 ID 和角色
func (s *MockTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	return 1, []string{"user"}, nil
}

// GenerateRefreshToken 生成一个固定的测试刷新 token
func (s *MockTokenService) GenerateRefreshToken(userID entity.UID, roles []string) (string, error) {
	return "test_refresh_token", nil
}

// ValidateRefreshToken 验证刷新 token，总是返回固定的用户 ID 和角色
func (s *MockTokenService) ValidateRefreshToken(refreshToken string) (entity.UID, []string, error) {
	return 1, []string{"user"}, nil
}

func (m *MockTokenService) ValidateTokenFromRequest(r *http.Request) (entity.UID, []string, error) {
	return 1, []string{"user"}, nil
}
