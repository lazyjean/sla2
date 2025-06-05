package converter

import (
	"encoding/json"
	"errors"
	"fmt"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

// QuestionConverter 处理 Question 实体和 protobuf 消息之间的转换
type QuestionConverter struct{}

// NewQuestionConverter 创建新的 QuestionConverter 实例
func NewQuestionConverter() *QuestionConverter {
	return &QuestionConverter{}
}

// ToEntityFromCreateRequest 将创建请求转换为问题实体
func (c *QuestionConverter) ToEntityFromCreateRequest(req *pb.QuestionServiceCreateRequest) (*entity.Question, error) {
	// 将 HyperTextTag 对象转换为 JSON
	content, err := json.Marshal(req.GetContent())
	if err != nil {
		return nil, err
	}

	// 将选项列表转换为 JSON
	options, err := json.Marshal(req.GetOptions())
	if err != nil {
		return nil, err
	}

	// 将选项双元组列表转换为 JSON
	optionTuples, err := json.Marshal(req.GetOptionTuples())
	if err != nil {
		return nil, err
	}

	// 创建问题实体
	return entity.NewQuestion(
		req.GetTitle(),
		content,
		req.GetSimpleQuestion(),
		req.GetQuestionType().String(),
		c.ToEntityDifficulty(req.GetDifficulty()),
		options,
		optionTuples,
		req.GetAnswers(),
		req.GetCategory().String(),
		req.GetLabels(),
		req.GetExplanation(),
		req.GetAttachments(),
		req.GetTimeLimit(),
	), nil
}

// ToEntityFromUpdateRequest 将更新请求转换为问题实体
func (c *QuestionConverter) ToEntityFromUpdateRequest(req *pb.QuestionServiceUpdateRequest) (*entity.Question, error) {
	if req == nil {
		return nil, errors.New("request cannot be nil")
	}

	// 将 HyperTextTag 对象转换为 JSON
	content, err := json.Marshal(req.GetContent())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal HyperTextTag: %w", err)
	}

	// 将选项列表转换为 JSON
	options, err := json.Marshal(req.GetOptions())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal options: %w", err)
	}

	// 将选项双元组列表转换为 JSON
	optionTuples, err := json.Marshal(req.GetOptionTuples())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal option tuples: %w", err)
	}

	// todo: 这里没有设置题目类型, 理论上更新时, 只更新有变化的字段, 这个需要验证
	return &entity.Question{
		ID:             entity.QuestionID(req.GetId()),
		Title:          req.GetTitle(),
		Content:        content,
		SimpleQuestion: req.GetSimpleQuestion(),
		Difficulty:     c.ToEntityDifficulty(req.GetDifficulty()),
		Options:        options,
		OptionTuples:   optionTuples,
		Answers:        req.GetAnswers(),
		Category:       req.GetCategory().String(),
		Labels:         req.GetLabels(),
		Explanation:    req.GetExplanation(),
		Attachments:    req.GetAttachments(),
		TimeLimit:      req.GetTimeLimit(),
	}, nil
}

// ToPB 将 Question 实体转换为 protobuf 消息
func (c *QuestionConverter) ToPB(q *entity.Question) *pb.Question {
	// Map difficulty level based on the difficulty string
	var pbDifficulty pb.QuestionDifficultyLevel
	switch q.Difficulty {
	case entity.DifficultyCefrA1:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A1
	case entity.DifficultyCefrA2:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A2
	case entity.DifficultyCefrB1:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_B1
	case entity.DifficultyCefrB2:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_B2
	case entity.DifficultyCefrC1:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_C1
	case entity.DifficultyCefrC2:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_C2
	case entity.DifficultyHsk1:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_1
	case entity.DifficultyHsk2:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_2
	case entity.DifficultyHsk3:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_3
	case entity.DifficultyHsk4:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_4
	case entity.DifficultyHsk5:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_5
	case entity.DifficultyHsk6:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_HSK_6
	default:
		pbDifficulty = pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_UNSPECIFIED
	}

	// Convert content
	var content *pb.HyperTextTag
	if len(q.Content) > 0 {
		if err := json.Unmarshal(q.Content, &content); err != nil {
			// If JSON conversion fails, use simple text content
			content = &pb.HyperTextTag{
				Type:  pb.HyperTextTagType_HYPER_TEXT_TAG_TYPE_TEXT,
				Value: q.SimpleQuestion,
			}
		}
	}

	// Convert options
	var options []*pb.QuestionOption
	if len(q.Options) > 0 {
		if err := json.Unmarshal(q.Options, &options); err != nil {
			options = []*pb.QuestionOption{}
		}
	}

	// Convert option tuples
	var optionTuples []*pb.QuestionOptionTuple
	if len(q.OptionTuples) > 0 {
		if err := json.Unmarshal(q.OptionTuples, &optionTuples); err != nil {
			optionTuples = []*pb.QuestionOptionTuple{}
		}
	}

	// Convert status
	var status pb.QuestionStatus
	switch q.Status {
	case "DRAFT":
		status = pb.QuestionStatus_QUESTION_STATUS_DRAFT
	case "REVIEWING":
		status = pb.QuestionStatus_QUESTION_STATUS_REVIEWING
	case "PUBLISHED":
		status = pb.QuestionStatus_QUESTION_STATUS_PUBLISHED
	default:
		status = pb.QuestionStatus_QUESTION_STATUS_UNSPECIFIED
	}

	// Convert category
	var category pb.QuestionCategory
	switch q.Category {
	case "EXERCISE":
		category = pb.QuestionCategory_QUESTION_CATEGORY_EXERCISE
	case "TEST":
		category = pb.QuestionCategory_QUESTION_CATEGORY_TEST
	case "GRAMMAR":
		category = pb.QuestionCategory_QUESTION_CATEGORY_GRAMMAR
	default:
		category = pb.QuestionCategory_QUESTION_CATEGORY_UNSPECIFIED
	}

	return &pb.Question{
		Id:             uint64(q.ID),
		Title:          q.Title,
		Content:        content,
		SimpleQuestion: q.SimpleQuestion,
		Options:        options,
		OptionTuples:   optionTuples,
		Answers:        q.Answers,
		Difficulty:     pbDifficulty,
		Status:         status,
		Category:       category,
		Labels:         q.Labels,
		Explanation:    q.Explanation,
		Attachments:    q.Attachments,
		CorrectRate:    q.CorrectRate,
		CreatedAt:      uint64(q.CreatedAt.Unix()),
		UpdatedAt:      uint64(q.UpdatedAt.Unix()),
		TimeLimit:      q.TimeLimit,
	}
}

// ToEntityDifficulty 将 protobuf 难度等级转换为实体难度等级
func (c *QuestionConverter) ToEntityDifficulty(difficulty pb.QuestionDifficultyLevel) string {
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

// ToEntityOptions 将 protobuf 选项转换为实体选项
func (c *QuestionConverter) ToEntityOptions(options []*pb.QuestionOption) []string {
	result := make([]string, len(options))
	for i, opt := range options {
		result[i] = opt.GetValue()
	}
	return result
}
