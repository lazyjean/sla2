package question

import (
	"context"
	"encoding/json"
	"strconv"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type QuestionService struct {
	pb.UnimplementedQuestionServiceServer
	questionService *service.QuestionService
	converter       *QuestionConverter
}

func NewQuestionService(questionService *service.QuestionService) *QuestionService {
	return &QuestionService{
		questionService: questionService,
		converter:       NewQuestionConverter(),
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
		Questions: []*pb.Question{s.converter.ToProto(question)},
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
		s.getDifficultyString(req.GetDifficulty()),
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

// getDifficultyString 将 QuestionDifficultyLevel 转换为对应的字符串常量
func (s *QuestionService) getDifficultyString(difficulty pb.QuestionDifficultyLevel) string {
	switch difficulty {
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A1:
		return entity.DifficultyCefrA1
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A2:
		return entity.DifficultyCefrA2
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_B1:
		return entity.DifficultyCefrB1
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_B2:
		return entity.DifficultyCefrB2
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_C1:
		return entity.DifficultyCefrC1
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_C2:
		return entity.DifficultyCefrC2
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_1:
		return entity.DifficultyHsk1
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_2:
		return entity.DifficultyHsk2
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_3:
		return entity.DifficultyHsk3
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_4:
		return entity.DifficultyHsk4
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_5:
		return entity.DifficultyHsk5
	case pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_6:
		return entity.DifficultyHsk6
	default:
		return entity.DifficultyCefrA1
	}
}

// convertOptionsToStrings 将 QuestionOption 转换为字符串数组
func convertOptionsToStrings(options []*pb.QuestionOption) []string {
	result := make([]string, len(options))
	for i, opt := range options {
		result[i] = opt.GetValue()
	}
	return result
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
		pbQuestions = append(pbQuestions, s.converter.ToProto(q))
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
		s.getDifficultyString(req.GetDifficulty()),
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
