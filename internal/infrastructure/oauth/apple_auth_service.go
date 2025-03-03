package oauth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/lazyjean/sla2/config"
)

const (
	appleJWKSURL = "https://appleid.apple.com/auth/keys"
)

// AppleConfig Apple 认证配置
type AppleConfig struct {
	ClientID   string
	TeamID     string
	KeyID      string
	PrivateKey string
}

// NewAppleConfig 从配置创建 Apple 配置
func NewAppleConfig(cfg *config.Config) *AppleConfig {
	return &AppleConfig{
		ClientID:   cfg.Apple.ClientID,
		TeamID:     cfg.Apple.TeamID,
		KeyID:      cfg.Apple.KeyID,
		PrivateKey: cfg.Apple.PrivateKey,
	}
}

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

// AppleAuthService 处理 Apple 认证相关的功能
type AppleAuthService struct {
	clientID    string
	teamID      string
	keyID       string
	privateKey  string
	appleKeys   map[string]*rsa.PublicKey
	appleKeysMu sync.RWMutex
	httpClient  *http.Client
}

func NewAppleAuthService(cfg *AppleConfig) *AppleAuthService {
	return &AppleAuthService{
		clientID:   cfg.ClientID,
		teamID:     cfg.TeamID,
		keyID:      cfg.KeyID,
		privateKey: cfg.PrivateKey,
		appleKeys:  make(map[string]*rsa.PublicKey),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// fetchApplePublicKeys 从 Apple 的 JWKS 端点获取公钥
func (s *AppleAuthService) fetchApplePublicKeys() error {
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
func (s *AppleAuthService) getApplePublicKey(kid string) (*rsa.PublicKey, error) {
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

// getAppleClientSecret 获取 Apple 的 client secret
func (s *AppleAuthService) getAppleClientSecret() (string, error) {
	// 创建 JWT header
	token := jwt.New(jwt.SigningMethodES256)
	token.Header["kid"] = s.keyID // Apple Key ID
	token.Header["alg"] = "ES256"

	// 设置 claims
	now := time.Now()
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = s.teamID                       // 你的 Team ID
	claims["iat"] = now.Unix()                     // 发布时间
	claims["exp"] = now.Add(24 * time.Hour).Unix() // 过期时间（24小时）
	claims["aud"] = "https://appleid.apple.com"    // 固定值
	claims["sub"] = s.clientID                     // 你的 Client ID (Bundle ID)

	// 从配置的 Base64 字符串解码私钥
	privateKeyBytes, err := base64.StdEncoding.DecodeString(s.privateKey)
	if err != nil {
		return "", fmt.Errorf("解码私钥失败: %v", err)
	}

	// 解析私钥
	privateKey, err := jwt.ParseECPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("解析私钥失败: %v", err)
	}

	// 签名生成 token
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("生成 client secret 失败: %v", err)
	}

	return tokenString, nil
}

// verifyIDToken 验证 Apple ID Token
func (s *AppleAuthService) verifyIDToken(ctx context.Context, idToken string) (*AppleIDToken, error) {
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
	if !ok || aud != s.clientID {
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

// AuthCodeWithApple 使用 authorization code 获取 Apple ID Token
func (s *AppleAuthService) AuthCodeWithApple(ctx context.Context, authorizationCode string) (*AppleIDToken, error) {
	// 1. 构建请求 Apple 的 token 请求
	tokenURL := "https://appleid.apple.com/auth/token"
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", authorizationCode)
	data.Set("client_id", s.clientID)
	clientSecret, err := s.getAppleClientSecret()
	if err != nil {
		return nil, fmt.Errorf("获取 Apple 的 client secret 失败: %v", err)
	}
	data.Set("client_secret", clientSecret)

	// 2. 发送请求获取 token
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 Apple token 失败: %v", err)
	}
	defer resp.Body.Close()

	// 3. 解析响应
	var tokenResp struct {
		IDToken string `json:"id_token"`
		Error   string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	if tokenResp.Error != "" {
		return nil, fmt.Errorf("Apple 返回错误: %s", tokenResp.Error)
	}

	if tokenResp.IDToken == "" {
		return nil, fmt.Errorf("未收到 ID Token")
	}

	// 4. 验证 ID Token
	return s.verifyIDToken(ctx, tokenResp.IDToken)
}
