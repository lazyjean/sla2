package grpc

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"google.golang.org/grpc"
)

// RegisterGateway registers the gRPC gateway.
func RegisterGateway(ctx context.Context, grpcEndpoint string, opts []grpc.DialOption) (http.Handler, error) {
	// 创建 HTTP 处理器
	mux := runtime.NewServeMux()

	// 注册服务处理器
	if err := pb.RegisterAIChatServiceHandlerFromEndpoint(ctx, mux, grpcEndpoint, opts); err != nil {
		return nil, fmt.Errorf("failed to register gateway: %w", err)
	}

	return mux, nil
}
