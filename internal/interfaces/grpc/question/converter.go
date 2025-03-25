package question

import (
	"encoding/json"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

// QuestionConverter 处理 Question 实体和 protobuf 消息之间的转换
type QuestionConverter struct{}

// NewQuestionConverter 创建新的 QuestionConverter 实例
func NewQuestionConverter() *QuestionConverter {
	return &QuestionConverter{}
}

// ToProto 将 Question 实体转换为 protobuf 消息
func (c *QuestionConverter) ToProto(q *entity.Question) *pb.Question {
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
