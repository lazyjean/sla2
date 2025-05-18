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
