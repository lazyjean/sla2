package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/lazyjean/sla2/config"
)

const (
	TokenCookieName = "access_token"
)

// SetTokenCookie sets the token cookie in the response
func SetTokenCookie(ctx context.Context, w http.ResponseWriter, token string) {
	// todo: 这里使用了全局配置, 其它地方都是使用的注入方式
	// 获取环境配置
	cfg := config.GetConfig()
	isDev := cfg.Server.Mode == "debug" // 使用 Server.Mode 判断环境

	cookie := &http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		Domain:   "",           // 让浏览器自动设置为当前域名
		MaxAge:   24 * 60 * 60, // 24 hours
		HttpOnly: true,
		// 开发环境使用较宽松的设置，生产环境使用严格设置
		Secure: !isDev, // 开发环境 false，生产环境 true
		SameSite: func() http.SameSite { // 开发环境 Lax，生产环境 Strict
			if isDev {
				return http.SameSiteLaxMode
			}
			return http.SameSiteStrictMode
		}(),
	}

	if token == "" {
		cookie.Expires = time.Unix(0, 0) // 设置过期时间为过去，使 cookie 立即失效
		cookie.MaxAge = -1
	}

	http.SetCookie(w, cookie)
}
