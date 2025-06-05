package question

import (
	"context"
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/transport/grpc/question/converter"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	pb.UnimplementedQuestionServiceServer
	questionService *service.QuestionService
	converter       *converter.QuestionConverter
}

func NewQuestionService(questionService *service.QuestionService, questionConverter *converter.QuestionConverter) *Service {
	return &Service{
		questionService: questionService,
		converter:       questionConverter,
	}
}

func (s *Service) Get(ctx context.Context, req *pb.QuestionServiceGetRequest) (*pb.QuestionServiceGetResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("GetQuestion called", zap.Any("req", req))

	// todo: 晚些再看这个接口是否合理, 先实现!!!, 这里的逻辑应该收到应用层, 或者领域内
	questions := make([]*pb.Question, 0)
	for id := range req.GetIds() {
		question, err := s.questionService.Get(ctx, entity.QuestionID(id))
		if err != nil {
			log.Error("GetQuestion failed", zap.Error(err))
			continue
		}
		questions = append(questions, s.converter.ToPB(question))
	}
	return &pb.QuestionServiceGetResponse{
		Questions: questions,
	}, nil
}

func (s *Service) Create(ctx context.Context, req *pb.QuestionServiceCreateRequest) (*pb.QuestionServiceCreateResponse, error) {
	log := logger.GetLogger(ctx)

	// 将请求转换为问题实体
	question, err := s.converter.ToEntityFromCreateRequest(req)
	if err != nil {
		log.Error("Failed to convert request to question entity", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to convert request to question entity")
	}

	// 创建问题
	id, err := s.questionService.Create(ctx, question)
	if err != nil {
		log.Error("Failed to create question", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question")
	}

	return &pb.QuestionServiceCreateResponse{
		Id: uint64(id),
	}, nil
}

func (s *Service) Search(ctx context.Context, req *pb.QuestionServiceSearchRequest) (*pb.QuestionServiceSearchResponse, error) {
	log := logger.GetLogger(ctx)

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
		pbQuestions = append(pbQuestions, s.converter.ToPB(q))
	}

	return &pb.QuestionServiceSearchResponse{
		Questions: pbQuestions,
		Total:     uint32(total),
	}, nil
}

func (s *Service) Update(ctx context.Context, req *pb.QuestionServiceUpdateRequest) (*pb.QuestionServiceUpdateResponse, error) {
	q, err := s.converter.ToEntityFromUpdateRequest(req)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to convert request to question entity")
	}
	if err := s.questionService.Update(ctx, q); err != nil {
		return nil, status.Error(codes.Internal, "failed to update question")
	}

	return &pb.QuestionServiceUpdateResponse{}, nil
}

func (s *Service) Delete(ctx context.Context, req *pb.QuestionServiceDeleteRequest) (*pb.QuestionServiceDeleteResponse, error) {
	err := s.questionService.Delete(ctx, entity.QuestionID(req.GetId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete question")
	}
	return &pb.QuestionServiceDeleteResponse{}, nil
}

func (s *Service) Publish(ctx context.Context, req *pb.QuestionServicePublishRequest) (*pb.QuestionServicePublishResponse, error) {
	if err := s.questionService.Publish(ctx, entity.QuestionID(req.GetId())); err != nil {
		return nil, status.Error(codes.Internal, "failed to publish question")
	}
	return &pb.QuestionServicePublishResponse{}, nil
}
