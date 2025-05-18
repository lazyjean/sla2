package question

import (
	"context"
	"encoding/json"
	"strconv"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/transport/grpc/question/converter"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QuestionService struct {
	pb.UnimplementedQuestionServiceServer
	questionService *service.QuestionService
	converter       *converter.QuestionConverter
}

func NewQuestionService(questionService *service.QuestionService) *QuestionService {
	return &QuestionService{
		questionService: questionService,
		converter:       converter.NewQuestionConverter(),
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
		Questions: []*pb.Question{s.converter.ToPB(question)},
	}, nil
}

func (s *QuestionService) Create(ctx context.Context, req *pb.QuestionServiceCreateRequest) (*pb.QuestionServiceCreateResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("CreateQuestion called", zap.Any("req", req))

	// 将 HyperTextTag 对象转换为 JSON
	content, err := json.Marshal(req.GetContent())
	if err != nil {
		log.Error("Failed to marshal HyperTextTag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal HyperTextTag")
	}

	// 将选项列表转换为 JSON
	options, err := json.Marshal(req.GetOptions())
	if err != nil {
		log.Error("Failed to marshal options", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal options")
	}

	// 将选项双元组列表转换为 JSON
	optionTuples, err := json.Marshal(req.GetOptionTuples())
	if err != nil {
		log.Error("Failed to marshal option tuples", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal option tuples")
	}

	// 将请求转换为 DTO
	questionDto := dto.NewCreateQuestionDTO(
		req.GetTitle(),
		content,
		req.GetSimpleQuestion(),
		req.GetQuestionType().String(),
		s.converter.ToEntityDifficulty(req.GetDifficulty()),
		options,
		optionTuples,
		req.GetAnswers(),
		req.GetCategory().String(),
		req.GetLabels(),
		req.GetExplanation(),
		req.GetAttachments(),
		req.GetTimeLimit(),
	)

	// 创建问题
	question, err := s.questionService.Create(ctx, questionDto)
	if err != nil {
		log.Error("CreateQuestion failed", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question")
	}

	return &pb.QuestionServiceCreateResponse{
		Id: uint64(question.ID),
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
		pbQuestions = append(pbQuestions, s.converter.ToPB(q))
	}

	return &pb.QuestionServiceSearchResponse{
		Questions: pbQuestions,
		Total:     uint32(total),
	}, nil
}

func (s *QuestionService) Update(ctx context.Context, req *pb.QuestionServiceUpdateRequest) (*pb.QuestionServiceUpdateResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("UpdateQuestion called", zap.Any("req", req))

	// 将 HyperTextTag 对象转换为 JSON
	content, err := json.Marshal(req.GetContent())
	if err != nil {
		log.Error("Failed to marshal HyperTextTag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal HyperTextTag")
	}

	// 将选项列表转换为 JSON
	options, err := json.Marshal(req.GetOptions())
	if err != nil {
		log.Error("Failed to marshal options", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal options")
	}

	// 将选项双元组列表转换为 JSON
	optionTuples, err := json.Marshal(req.GetOptionTuples())
	if err != nil {
		log.Error("Failed to marshal option tuples", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to marshal option tuples")
	}

	// 获取当前问题以保持题目类型不变
	question, err := s.questionService.Get(ctx, strconv.FormatUint(uint64(req.GetId()), 10))
	if err != nil {
		log.Error("Failed to get question", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to get question")
	}

	// 创建更新 DTO
	updateDTO := dto.NewUpdateQuestionDTO(
		strconv.FormatUint(uint64(req.GetId()), 10),
		req.GetTitle(),
		content,
		req.GetSimpleQuestion(),
		question.Type, // 保持原有的题目类型
		s.converter.ToEntityDifficulty(req.GetDifficulty()),
		options,
		optionTuples,
		req.GetAnswers(),
		req.GetCategory().String(),
		req.GetLabels(),
		req.GetExplanation(),
		req.GetAttachments(),
		req.GetTimeLimit(),
	)

	// 更新问题
	_, err = s.questionService.Update(ctx, updateDTO)
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

// ListTag 实现标签列表查询接口
func (s *QuestionService) ListTag(ctx context.Context, req *pb.QuestionTagServiceListTagRequest) (*pb.QuestionTagServiceListTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("ListTag called", zap.Any("req", req))

	tags, err := s.questionService.FindAllTags(ctx, int(req.GetTopN()))
	if err != nil {
		log.Error("failed to list question tags", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list question tags")
	}

	return &pb.QuestionTagServiceListTagResponse{
		Tags: s.converter.ToPBTagList(tags),
	}, nil
}

// CreateTag 实现创建标签接口
func (s *QuestionService) CreateTag(ctx context.Context, req *pb.QuestionTagServiceCreateTagRequest) (*pb.QuestionTagServiceCreateTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("CreateTag called", zap.Any("req", req))

	createdTag, err := s.questionService.CreateTag(ctx, req.GetName())
	if err != nil {
		log.Error("failed to create question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to create question tag")
	}

	log.Info("created tag", zap.String("id", createdTag.ID), zap.String("name", createdTag.Name))
	return &pb.QuestionTagServiceCreateTagResponse{}, nil
}

// UpdateTag 实现更新标签接口
func (s *QuestionService) UpdateTag(ctx context.Context, req *pb.QuestionTagServiceUpdateTagRequest) (*pb.QuestionTagServiceUpdateTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("UpdateTag called", zap.Any("req", req))

	updatedTag, err := s.questionService.UpdateTag(ctx, req.GetName(), req.GetName())
	if err != nil {
		log.Error("failed to update question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to update question tag")
	}

	log.Info("updated tag", zap.String("id", updatedTag.ID), zap.String("name", updatedTag.Name))
	return &pb.QuestionTagServiceUpdateTagResponse{}, nil
}

// DeleteTag 实现删除标签接口
func (s *QuestionService) DeleteTag(ctx context.Context, req *pb.QuestionTagServiceDeleteTagRequest) (*pb.QuestionTagServiceDeleteTagResponse, error) {
	log := logger.GetLogger(ctx)
	log.Info("DeleteTag called", zap.Any("req", req))

	if err := s.questionService.DeleteTag(ctx, req.GetName()); err != nil {
		log.Error("failed to delete question tag", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to delete question tag")
	}

	return &pb.QuestionTagServiceDeleteTagResponse{}, nil
}
