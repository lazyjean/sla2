package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// MetadataUnaryServerInterceptor 创建一个用于处理元数据的拦截器
func MetadataUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 获取传入的元数据
		var incomingMD metadata.MD
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			incomingMD = md
		}

		// 创建新的传出元数据
		outgoingMD := metadata.New(map[string]string{
			"content-type": "application/json",
		})

		// 合并传入和传出的元数据
		outgoingMD = metadata.Join(outgoingMD, incomingMD)

		// 创建新的 context，包含传出元数据
		newCtx := metadata.NewOutgoingContext(ctx, outgoingMD)

		// 调用处理器
		resp, err := handler(newCtx, req)

		// 如果没有错误，设置 header 和 trailer
		if err == nil {
			// 设置 header
			if err := grpc.SetHeader(ctx, outgoingMD); err != nil {
				return resp, err
			}
			// 设置 trailer
			if err := grpc.SetTrailer(ctx, outgoingMD); err != nil {
				return resp, err
			}
		}

		return resp, err
	}
}
