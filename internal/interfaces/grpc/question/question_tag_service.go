package question

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// QuestionTagServiceGRPC 是标签服务的gRPC实现
type QuestionTagServiceGRPC struct {
	pb.UnimplementedQuestionTagServiceServer
	tagService *service.QuestionTagService // 依赖应用服务层而非直接依赖仓储
}

// NewQuestionTagServiceGRPC 创建标签服务的gRPC实现
func NewQuestionTagServiceGRPC(tagService *service.QuestionTagService) *QuestionTagServiceGRPC {
	return &QuestionTagServiceGRPC{
		tagService: tagService,
	}
}

// ListTag 实现标签列表查询接口
func (s *QuestionTagServiceGRPC) ListTag(ctx context.Context, req *pb.QuestionTagServiceListTagRequest) (*pb.QuestionTagServiceListTagResponse, error) {
	log := logger.GetLogger(ctx) // 从context获取logger而非保存在结构体中
	log.Info("ListTag called", zap.Any("req", req))

	// 调用应用服务层而非直接调用仓储
	tags, err := s.tagService.FindAll(ctx, int(req.GetTopN()))
	if err != nil {
		log.Error("failed to list question tags", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list question tags")
	}

	// 构建响应
	pbTags := make([]*pb.QuestionTag, 0, len(tags))
	for _, tag := range tags {
		// 标签权重设置为默认值1
		pbTags = append(pbTags, &pb.QuestionTag{
			Name:   tag.Name,
			Weight: 1,
		})
	}

	return &pb.QuestionTagServiceListTagResponse{
		Tags: pbTags,
	}, nil
}

// CreateTag 实现创建标签接口
func (s *QuestionTagServiceGRPC) CreateTag(ctx context.Context, req *pb.QuestionTagServiceCreateTagRequest) (*pb.QuestionTagServiceCreateTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("CreateTag called", zap.Any("req", req))

	// 调用应用服务层
	createdTag, err := s.tagService.Create(ctx, req.GetName())
	if err != nil {
		log.Error("failed to create question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question tag")
	}

	log.Info("created tag", zap.String("id", createdTag.ID), zap.String("name", createdTag.Name))
	return &pb.QuestionTagServiceCreateTagResponse{}, nil
}

// UpdateTag 实现更新标签接口
func (s *QuestionTagServiceGRPC) UpdateTag(ctx context.Context, req *pb.QuestionTagServiceUpdateTagRequest) (*pb.QuestionTagServiceUpdateTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("UpdateTag called", zap.Any("req", req))

	// 调用应用服务层
	updatedTag, err := s.tagService.Update(ctx, req.GetName(), req.GetName())
	if err != nil {
		log.Error("failed to update question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update question tag")
	}

	log.Info("updated tag", zap.String("id", updatedTag.ID), zap.String("name", updatedTag.Name))
	return &pb.QuestionTagServiceUpdateTagResponse{}, nil
}

// DeleteTag 实现删除标签接口
func (s *QuestionTagServiceGRPC) DeleteTag(ctx context.Context, req *pb.QuestionTagServiceDeleteTagRequest) (*pb.QuestionTagServiceDeleteTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("DeleteTag called", zap.Any("req", req))

	// 调用应用服务层
	if err := s.tagService.Delete(ctx, req.GetName()); err != nil {
		log.Error("failed to delete question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete question tag")
	}

	return &pb.QuestionTagServiceDeleteTagResponse{}, nil
}
