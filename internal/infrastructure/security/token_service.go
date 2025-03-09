package security

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

	"github.com/lazyjean/sla2/config"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/security"
)

// TokenError 定义 token 相关错误
type TokenError struct {
	Message string
	Err     error
}

func (e *TokenError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// JWTTokenService JWT 令牌服务实现
type JWTTokenService struct {
	tokenSecretKey    string
	refreshSecretKey  string
	tokenExpiration   time.Duration
	refreshExpiration time.Duration
}

// NewJWTTokenService 创建一个新的 JWT 令牌服务
func NewJWTTokenService(cfg *config.Config) security.TokenService {
	return &JWTTokenService{
		tokenSecretKey:    cfg.JWT.TokenSecretKey,
		refreshSecretKey:  cfg.JWT.RefreshSecretKey,
		tokenExpiration:   time.Duration(cfg.JWT.TokenExpiration) * time.Hour,
		refreshExpiration: time.Duration(cfg.JWT.RefreshExpiration) * time.Hour,
	}
}

// GenerateToken 生成访问令牌
func (s *JWTTokenService) GenerateToken(userID entity.UID, roles []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"exp":   now.Add(s.tokenExpiration).Unix(),
		"iat":   now.Unix(),
		"jti":   uuid.New().String(),
		"typ":   "access",
		"roles": roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.tokenSecretKey))
	if err != nil {
		return "", &TokenError{Message: "failed to sign token", Err: err}
	}
	return signedToken, nil
}

// ValidateToken 验证访问令牌
func (s *JWTTokenService) ValidateToken(tokenString string) (entity.UID, []string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &TokenError{Message: fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])}
		}
		return []byte(s.tokenSecretKey), nil
	})

	if err != nil {
		return 0, nil, &TokenError{Message: "failed to parse token", Err: err}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, nil, &TokenError{Message: "invalid token"}
	}

	// 验证 token 类型
	if typ, ok := claims["typ"].(string); !ok || typ != "access" {
		return 0, nil, &TokenError{Message: "invalid token type"}
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, nil, &TokenError{Message: "invalid subject claim"}
	}
	userID := entity.UID(sub)

	// 获取角色信息
	roles := make([]string, 0)
	if rolesInterface, exists := claims["roles"]; exists {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			for _, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	return userID, roles, nil
}

// ValidateTokenFromContext 从上下文中验证令牌
func (s *JWTTokenService) ValidateTokenFromContext(ctx context.Context) (entity.UID, []string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, nil, &TokenError{Message: "missing metadata"}
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return 0, nil, &TokenError{Message: "missing authorization header"}
	}

	auth := values[0]
	if !strings.HasPrefix(auth, "Bearer ") {
		return 0, nil, &TokenError{Message: "invalid authorization format"}
	}

	token := strings.TrimPrefix(auth, "Bearer ")
	return s.ValidateToken(token)
}

// GenerateRefreshToken 生成刷新令牌
func (s *JWTTokenService) GenerateRefreshToken(userID entity.UID, roles []string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub":   userID,
		"exp":   now.Add(s.refreshExpiration).Unix(),
		"iat":   now.Unix(),
		"jti":   uuid.New().String(),
		"typ":   "refresh",
		"roles": roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.refreshSecretKey))
	if err != nil {
		return "", &TokenError{Message: "failed to sign refresh token", Err: err}
	}
	return signedToken, nil
}

// ValidateRefreshToken 验证刷新令牌
func (s *JWTTokenService) ValidateRefreshToken(refreshToken string) (entity.UID, []string, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, &TokenError{Message: fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])}
		}
		return []byte(s.refreshSecretKey), nil
	})

	if err != nil {
		return 0, nil, &TokenError{Message: "failed to parse refresh token", Err: err}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, nil, &TokenError{Message: "invalid refresh token"}
	}

	// 验证 token 类型
	if typ, ok := claims["typ"].(string); !ok || typ != "refresh" {
		return 0, nil, &TokenError{Message: "invalid token type"}
	}

	sub, ok := claims["sub"].(float64)
	if !ok {
		return 0, nil, &TokenError{Message: "invalid subject claim"}
	}

	// 获取角色信息
	roles := make([]string, 0)
	if rolesInterface, exists := claims["roles"]; exists {
		if rolesArray, ok := rolesInterface.([]interface{}); ok {
			for _, role := range rolesArray {
				if roleStr, ok := role.(string); ok {
					roles = append(roles, roleStr)
				}
			}
		}
	}

	return entity.UID(sub), roles, nil
}

// ValidateTokenFromRequest 从HTTP请求中验证令牌
func (s *JWTTokenService) ValidateTokenFromRequest(r *http.Request) (entity.UID, []string, error) {
	var token string

	// 优先从请求头获取令牌
	auth := r.Header.Get("Authorization")
	if auth != "" {
		// 验证头部格式
		if !strings.HasPrefix(auth, "Bearer ") {
			return 0, nil, &TokenError{Message: "invalid authorization format"}
		}

		// 提取令牌
		token = strings.TrimPrefix(auth, "Bearer ")
	} else {
		// 如果请求头中没有令牌，尝试从URL查询参数获取
		token = r.URL.Query().Get("token")
		if token == "" {
			return 0, nil, &TokenError{Message: "token not found in header or query parameters"}
		}
	}

	// 验证令牌
	return s.ValidateToken(token)
}

var _ security.TokenService = (*JWTTokenService)(nil)
