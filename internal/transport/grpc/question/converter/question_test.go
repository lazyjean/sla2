package converter

import (
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

// ... existing tests ...

func TestQuestionConverter_ToPBTag(t *testing.T) {
	converter := NewQuestionConverter()

	tests := []struct {
		name     string
		tag      *entity.QuestionTag
		expected *pb.QuestionTag
	}{
		{
			name: "正常转换",
			tag: &entity.QuestionTag{
				ID:        "1",
				Name:      "test-tag",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &pb.QuestionTag{
				Name:   "test-tag",
				Weight: 0,
			},
		},
		{
			name:     "空标签",
			tag:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToPBTag(tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQuestionConverter_ToEntityTag(t *testing.T) {
	converter := NewQuestionConverter()

	tests := []struct {
		name     string
		tag      *pb.QuestionTag
		expected *entity.QuestionTag
	}{
		{
			name: "正常转换",
			tag: &pb.QuestionTag{
				Name:   "test-tag",
				Weight: 1,
			},
			expected: &entity.QuestionTag{
				Name: "test-tag",
			},
		},
		{
			name:     "空标签",
			tag:      nil,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToEntityTag(tt.tag)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQuestionConverter_ToPBTagList(t *testing.T) {
	converter := NewQuestionConverter()

	tests := []struct {
		name     string
		tags     []*entity.QuestionTag
		expected []*pb.QuestionTag
	}{
		{
			name: "正常转换列表",
			tags: []*entity.QuestionTag{
				{
					ID:   "1",
					Name: "tag1",
				},
				{
					ID:   "2",
					Name: "tag2",
				},
			},
			expected: []*pb.QuestionTag{
				{
					Name:   "tag1",
					Weight: 0,
				},
				{
					Name:   "tag2",
					Weight: 0,
				},
			},
		},
		{
			name:     "空列表",
			tags:     nil,
			expected: nil,
		},
		{
			name:     "空切片",
			tags:     []*entity.QuestionTag{},
			expected: []*pb.QuestionTag{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ToPBTagList(tt.tags)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQuestionConverter_ToEntityFromCreateRequest(t *testing.T) {
	converter := NewQuestionConverter()

	tests := []struct {
		name    string
		req     *pb.QuestionServiceCreateRequest
		wantErr bool
	}{
		{
			name: "正常转换",
			req: &pb.QuestionServiceCreateRequest{
				Title:          "测试问题",
				Content:        &pb.HyperTextTag{Type: pb.HyperTextTagType_HYPER_TEXT_TAG_TYPE_TEXT, Value: "测试内容"},
				SimpleQuestion: "简单问题",
				QuestionType:   pb.QuestionType_QUESTION_TYPE_SINGLE_CHOICE,
				Difficulty:     pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A1,
				Options: []*pb.QuestionOption{
					{Value: "选项1"},
					{Value: "选项2"},
				},
				OptionTuples: []*pb.QuestionOptionTuple{
					{
						Option1: &pb.QuestionOption{Value: "选项1"},
						Option2: &pb.QuestionOption{Value: "选项2"},
					},
				},
				Answers:     []string{"选项1"},
				Category:    pb.QuestionCategory_QUESTION_CATEGORY_EXERCISE,
				Labels:      []string{"标签1", "标签2"},
				Explanation: "解释",
				Attachments: []string{"附件1"},
				TimeLimit:   60,
			},
			wantErr: false,
		},
		{
			name:    "空请求",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ToEntityFromCreateRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityFromCreateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ToEntityFromCreateRequest() got = nil, want non-nil")
			}
		})
	}
}

func TestQuestionConverter_ToEntityFromUpdateRequest(t *testing.T) {
	converter := NewQuestionConverter()

	tests := []struct {
		name             string
		req              *pb.QuestionServiceUpdateRequest
		existingQuestion *entity.Question
		wantErr          bool
	}{
		{
			name: "正常转换",
			req: &pb.QuestionServiceUpdateRequest{
				Id:             1,
				Title:          "测试问题",
				Content:        &pb.HyperTextTag{Type: pb.HyperTextTagType_HYPER_TEXT_TAG_TYPE_TEXT, Value: "测试内容"},
				SimpleQuestion: "简单问题",
				Difficulty:     pb.QuestionDifficultyLevel_QUESTION_DIFFICULTY_LEVEL_CEFR_A1,
				Options: []*pb.QuestionOption{
					{Value: "选项1"},
					{Value: "选项2"},
				},
				OptionTuples: []*pb.QuestionOptionTuple{
					{
						Option1: &pb.QuestionOption{Value: "选项1"},
						Option2: &pb.QuestionOption{Value: "选项2"},
					},
				},
				Answers:     []string{"选项1"},
				Category:    pb.QuestionCategory_QUESTION_CATEGORY_EXERCISE,
				Labels:      []string{"标签1", "标签2"},
				Explanation: "解释",
				Attachments: []string{"附件1"},
				TimeLimit:   60,
			},
			existingQuestion: &entity.Question{
				Type: "single_choice",
			},
			wantErr: false,
		},
		{
			name:             "空请求",
			req:              nil,
			existingQuestion: &entity.Question{},
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := converter.ToEntityFromUpdateRequest(tt.req, tt.existingQuestion)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToEntityFromUpdateRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("ToEntityFromUpdateRequest() got = nil, want non-nil")
			}
		})
	}
}
