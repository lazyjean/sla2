package grpc

import (
	"context"
	"google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

// healthServer 健康检查服务实现
type healthServer struct {
	grpc_health_v1.UnimplementedHealthServer
}

func (s *healthServer) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (s *healthServer) Watch(req *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	// 创建一个定时器，每5秒发送一次健康状态
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// 定期发送健康状态
	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case <-ticker.C:
			if err := stream.Send(&grpc_health_v1.HealthCheckResponse{
				Status: grpc_health_v1.HealthCheckResponse_SERVING,
			}); err != nil {
				return err
			}
		}
	}
}
