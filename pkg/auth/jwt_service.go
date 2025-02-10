package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	appleJWKSURL = "https://appleid.apple.com/auth/keys"
)

// AppleIDToken Apple ID Token 信息
type AppleIDToken struct {
	Subject string // 用户的唯一标识符
	Email   string // 用户的邮箱
	Name    string // 用户的名字
}

// JWKS Apple 的 JWKS 响应
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK JSON Web Key
type JWK struct {
	Kty string `json:"kty"` // Key Type
	Kid string `json:"kid"` // Key ID
	Use string `json:"use"` // Use (sig - signature)
	Alg string `json:"alg"` // Algorithm
	N   string `json:"n"`   // Modulus
	E   string `json:"e"`   // Exponent
}

// JWTServicer 定义了 JWT 服务的接口
type JWTServicer interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPassword, password string) bool
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (uint, error)
	GenerateRandomPassword() string
	GenerateRandomString(length int) string
	VerifyAppleIDToken(ctx context.Context, idToken string) (*AppleIDToken, error)
}

// JWTService 是 JWTServicer 接口的实现
type JWTService struct {
	secretKey     string
	appleClientID string
	appleKeys     map[string]*rsa.PublicKey
	appleKeysMu   sync.RWMutex
	httpClient    *http.Client
}

func NewJWTService(secretKey string, appleClientID string) *JWTService {
	return &JWTService{
		secretKey:     secretKey,
		appleClientID: appleClientID,
		appleKeys:     make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fetchApplePublicKeys 从 Apple 的 JWKS 端点获取公钥
func (s *JWTService) fetchApplePublicKeys() error {
	resp, err := s.httpClient.Get(appleJWKSURL)
	if err != nil {
		return fmt.Errorf("failed to fetch apple jwks: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read apple jwks response: %v", err)
	}

	var jwks JWKS
	if err := json.Unmarshal(body, &jwks); err != nil {
		return fmt.Errorf("failed to parse apple jwks: %v", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)
	for _, key := range jwks.Keys {
		if key.Use != "sig" || key.Kty != "RSA" {
			continue
		}

		// 解码模数和指数
		nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
		if err != nil {
			continue
		}
		eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
		if err != nil {
			continue
		}

		// 转换为大整数
		n := new(big.Int).SetBytes(nBytes)
		e := new(big.Int).SetBytes(eBytes)

		// 创建 RSA 公钥
		pubKey := &rsa.PublicKey{
			N: n,
			E: int(e.Int64()),
		}

		newKeys[key.Kid] = pubKey
	}

	// 更新公钥缓存
	s.appleKeysMu.Lock()
	s.appleKeys = newKeys
	s.appleKeysMu.Unlock()

	return nil
}

// getApplePublicKey 获取指定 kid 的公钥
func (s *JWTService) getApplePublicKey(kid string) (*rsa.PublicKey, error) {
	s.appleKeysMu.RLock()
	key, ok := s.appleKeys[kid]
	s.appleKeysMu.RUnlock()

	if !ok {
		// 如果没有找到对应的公钥，尝试重新获取
		if err := s.fetchApplePublicKeys(); err != nil {
			return nil, err
		}

		s.appleKeysMu.RLock()
		key, ok = s.appleKeys[kid]
		s.appleKeysMu.RUnlock()

		if !ok {
			return nil, fmt.Errorf("apple public key not found for kid: %s", kid)
		}
	}

	return key, nil
}

func (s *JWTService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedBytes), nil
}

func (s *JWTService) ComparePasswords(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func (s *JWTService) GenerateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return signedToken, nil
}

func (s *JWTService) ValidateToken(tokenString string) (uint, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id in token")
	}

	return uint(userID), nil
}

func (s *JWTService) GenerateRandomPassword() string {
	const length = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))[:length]
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func (s *JWTService) GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return base64.URLEncoding.EncodeToString([]byte(time.Now().String()))[:length]
	}

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}

func (s *JWTService) VerifyAppleIDToken(ctx context.Context, idToken string) (*AppleIDToken, error) {
	// 解析 ID Token
	token, err := jwt.Parse(idToken, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 获取 kid
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		// 获取对应的公钥
		return s.getApplePublicKey(kid)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse apple id token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// 验证 audience
	aud, ok := claims["aud"].(string)
	if !ok || aud != s.appleClientID {
		return nil, fmt.Errorf("invalid audience in apple id token")
	}

	// 获取必要信息
	sub, ok := claims["sub"].(string)
	if !ok {
		return nil, fmt.Errorf("missing subject in token")
	}

	email, _ := claims["email"].(string)

	return &AppleIDToken{
		Subject: sub,
		Email:   email,
		Name:    "", // Apple 不会在 ID Token 中返回用户名，需要在前端获取
	}, nil
}
