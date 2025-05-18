package test

import (
	"context"

	domainOauth "github.com/lazyjean/sla2/internal/domain/oauth"
)

// MockAppleAuthService 是一个用于测试的 Apple auth service
type MockAppleAuthService struct{}

// NewMockAppleAuthService 创建一个新的 mock Apple auth service
func NewMockAppleAuthService() domainOauth.AppleAuthService {
	return &MockAppleAuthService{}
}

// AuthCodeWithApple 验证 Apple 授权码，总是返回固定的用户信息
func (s *MockAppleAuthService) AuthCodeWithApple(ctx context.Context, authorizationCode string) (*domainOauth.AppleIDToken, error) {
	return &domainOauth.AppleIDToken{
		Subject: "test_user_id",
		Email:   "test@example.com",
		Name:    "Test User",
	}, nil
}
