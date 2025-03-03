package valueobject

// AuthToken 认证令牌值对象
type AuthToken struct {
	AccessToken  string // 访问令牌
	RefreshToken string // 刷新令牌
} 