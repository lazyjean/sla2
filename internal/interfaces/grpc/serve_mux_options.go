package grpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpcmiddleware "github.com/lazyjean/sla2/internal/interfaces/grpc/middleware"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

// WithMetadata 从 HTTP 请求中提取 token 并设置到 gRPC metadata 中
func WithMetadata(ctx context.Context, req *http.Request) metadata.MD {
	log := logger.GetLogger(ctx)
	log.Info("gRPC-Gateway request",
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("query", req.URL.RawQuery),
		zap.String("remote_addr", req.RemoteAddr),
		zap.String("user_agent", req.UserAgent()),
	)

	md := metadata.New(nil)

	// 从 Authorization header 中提取 token
	if auth := req.Header.Get("Authorization"); auth != "" {
		// 检查是否是 Bearer token
		if len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
			// 提取 JWT token
			token := auth[7:]
			md = metadata.Join(md, metadata.Pairs(grpcmiddleware.MDHeaderAccessToken, token))
			log.Debug("Extracted token from Authorization header",
				zap.String("token", token[:10]+"..."), // 只记录 token 的前 10 个字符
			)
		}
	}

	// 从 cookie 中提取 token
	if cookie, err := req.Cookie(grpcmiddleware.HTTPCookieAccessTokenName); err == nil {
		md = metadata.Join(md, metadata.Pairs(grpcmiddleware.MDHeaderAccessToken, cookie.Value))
		log.Debug("Extracted token from cookie",
			zap.String("token", cookie.Value[:10]+"..."), // 只记录 token 的前 10 个字符
		)
	}

	return md
}

// WithForwardResponseOption 处理响应转发选项
func WithForwardResponseOption(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	log := logger.GetLogger(ctx)
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		log.Warn("No server metadata in context")
		return nil
	}

	// 处理 cookie 设置
	if vals := md.HeaderMD.Get(grpcmiddleware.MDHeaderAccessToken); len(vals) > 0 {
		grpcmiddleware.SetAccessTokenInHTTPResponseCookie(ctx, w, vals[0])
	}

	return nil
}

// WithErrorHandler 处理错误
func WithErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	log := logger.GetLogger(ctx)
	log.Error("gRPC-Gateway error",
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Error(err),
	)
	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}

// WithOutgoingHeaderMatcher 处理出站 header 匹配
func WithOutgoingHeaderMatcher(key string) (string, bool) {
	// 不需要转换任何 metadata 为 header，因为我们使用 SetCookie 直接设置
	return "", false
}
