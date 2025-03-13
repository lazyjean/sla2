package question

import (
	"context"
	"strconv"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QuestionService struct {
	pb.UnimplementedQuestionServiceServer
	questionService *service.QuestionService
}

func NewQuestionService(questionService *service.QuestionService) *QuestionService {
	return &QuestionService{
		questionService: questionService,
	}
}

func (s *QuestionService) Get(ctx context.Context, req *pb.QuestionServiceGetRequest) (*pb.QuestionServiceGetResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("GetQuestion called", zap.Any("req", req))

	// 应用层方法需要调整为接收单个ID
	// 将uint32 ID转换为字符串格式
	id := strconv.FormatUint(uint64(req.GetIds()[0]), 10)
	question, err := s.questionService.Get(ctx, id)
	if err != nil {
		log.Error("GetQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get question")
	}

	return &pb.QuestionServiceGetResponse{
		Questions: []*pb.Question{question.ToProto()},
	}, nil
}

func (s *QuestionService) Create(ctx context.Context, req *pb.QuestionServiceCreateRequest) (*pb.QuestionServiceCreateResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("CreateQuestion called", zap.Any("req", req))

	// 从请求中提取必要信息
	title := req.GetTitle()
	content := req.GetSimpleQuestion() // 使用SimpleQuestion作为content
	tags := req.GetTags()
	creatorID := "" // 可以从上下文获取，或者添加到请求参数中

	question, err := s.questionService.Create(ctx, title, content, tags, creatorID)
	if err != nil {
		log.Error("CreateQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question")
	}

	return &pb.QuestionServiceCreateResponse{
		Id: uint32(question.ID),
	}, nil
}

func (s *QuestionService) Search(ctx context.Context, req *pb.QuestionServiceSearchRequest) (*pb.QuestionServiceSearchResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("SearchQuestions called", zap.Any("req", req))

	keyword := req.GetKeyword()
	// SearchRequest中没有tags字段，传空切片
	var tags []string
	page := int(req.GetPage())
	pageSize := int(req.GetPageSize())

	questions, total, err := s.questionService.Search(ctx, keyword, tags, page, pageSize)
	if err != nil {
		log.Error("SearchQuestions failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to search questions")
	}

	var pbQuestions []*pb.Question
	for _, q := range questions {
		pbQuestions = append(pbQuestions, q.ToProto())
	}

	return &pb.QuestionServiceSearchResponse{
		Questions: pbQuestions,
		Total:     uint32(total),
	}, nil
}

func (s *QuestionService) Update(ctx context.Context, req *pb.QuestionServiceUpdateRequest) (*pb.QuestionServiceUpdateResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("UpdateQuestion called", zap.Any("req", req))

	id := strconv.FormatUint(uint64(req.GetId()), 10)
	title := req.GetTitle()
	content := req.GetSimpleQuestion() // 使用SimpleQuestion作为content
	// UpdateRequest中可能没有tags字段，检查proto定义
	var tags []string

	_, err := s.questionService.Update(ctx, id, title, content, tags)
	if err != nil {
		log.Error("UpdateQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update question")
	}

	return &pb.QuestionServiceUpdateResponse{}, nil
}

func (s *QuestionService) Delete(ctx context.Context, req *pb.QuestionServiceDeleteRequest) (*pb.QuestionServiceDeleteResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("DeleteQuestion called", zap.Any("req", req))

	id := strconv.FormatUint(uint64(req.GetId()), 10)
	err := s.questionService.Delete(ctx, id)
	if err != nil {
		log.Error("DeleteQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete question")
	}

	return &pb.QuestionServiceDeleteResponse{}, nil
}

func (s *QuestionService) Publish(ctx context.Context, req *pb.QuestionServicePublishRequest) (*pb.QuestionServicePublishResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("PublishQuestion called", zap.Any("req", req))

	id := strconv.FormatUint(uint64(req.GetId()), 10)
	_, err := s.questionService.Publish(ctx, id)
	if err != nil {
		log.Error("PublishQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to publish question")
	}

	return &pb.QuestionServicePublishResponse{}, nil
}
