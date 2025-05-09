package service

import (
	"context"
	"testing"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockLearningRepository 是学习进度仓储的mock实现
type MockLearningRepository struct {
	mock.Mock
}

// MockMemoryService 是记忆服务的mock实现
type MockMemoryService struct {
	mock.Mock
}

func (m *MockMemoryService) RecordLearningResult(ctx context.Context, memoryUnitID uint32, result bool, responseTime uint32) error {
	args := m.Called(ctx, memoryUnitID, result, responseTime)
	return args.Error(0)
}

func (m *MockMemoryService) GetMemoryUnit(ctx context.Context, memoryUnitID uint32) (*entity.MemoryUnit, error) {
	args := m.Called(ctx, memoryUnitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.MemoryUnit), args.Error(1)
}

func (m *MockMemoryService) UpdateMemoryUnit(ctx context.Context, memoryUnit *entity.MemoryUnit) error {
	args := m.Called(ctx, memoryUnit)
	return args.Error(0)
}

func (m *MockMemoryService) GetLearningProgress(ctx context.Context) (*LearningProgress, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LearningProgress), args.Error(1)
}

func (m *MockMemoryService) GetNextReviewWords(ctx context.Context, limit int) ([]*entity.Word, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Word), args.Error(1)
}

func (m *MockMemoryService) GetWordStats(ctx context.Context, wordID entity.WordID) (*WordStats, error) {
	args := m.Called(ctx, wordID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*WordStats), args.Error(1)
}

func (m *MockMemoryService) ReviewWord(ctx context.Context, wordID entity.WordID, result bool, responseTime uint32) error {
	args := m.Called(ctx, wordID, result, responseTime)
	return args.Error(0)
}

func (m *MockMemoryService) ReviewHanChar(ctx context.Context, hanCharID uint32, result bool, responseTime uint32) error {
	args := m.Called(ctx, hanCharID, result, responseTime)
	return args.Error(0)
}

func (m *MockMemoryService) UpdateMemoryStatus(ctx context.Context, memoryUnitID uint32, masteryLevel entity.MasteryLevel, studyDuration uint32) error {
	args := m.Called(ctx, memoryUnitID, masteryLevel, studyDuration)
	return args.Error(0)
}

func (m *MockMemoryService) ListMemoriesForReview(ctx context.Context, page, pageSize uint32, types []entity.MemoryUnitType) ([]*entity.MemoryUnit, int, error) {
	args := m.Called(ctx, page, pageSize, types)
	units, _ := args.Get(0).([]*entity.MemoryUnit)
	total, _ := args.Get(1).(int)
	return units, total, args.Error(2)
}

func (m *MockMemoryService) GetMemoryStats(ctx context.Context, unitType *entity.MemoryUnitType) (*repository.MemoryUnitStats, error) {
	args := m.Called(ctx, unitType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.MemoryUnitStats), args.Error(1)
}

func (m *MockLearningRepository) SaveCourseProgress(ctx context.Context, progress *entity.CourseLearningProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockLearningRepository) GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error) {
	args := m.Called(ctx, userID, courseID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CourseLearningProgress), args.Error(1)
}

func (m *MockLearningRepository) ListCourseProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.CourseLearningProgress, int64, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).([]*entity.CourseLearningProgress), args.Get(1).(int64), args.Error(2)
}

func (m *MockLearningRepository) SaveSectionProgress(ctx context.Context, progress *entity.CourseSectionProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockLearningRepository) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.CourseSectionProgress, error) {
	args := m.Called(ctx, userID, sectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CourseSectionProgress), args.Error(1)
}

func (m *MockLearningRepository) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.CourseSectionProgress, error) {
	args := m.Called(ctx, userID, courseID)
	return args.Get(0).([]*entity.CourseSectionProgress), args.Error(1)
}

func (m *MockLearningRepository) SaveUnitProgress(ctx context.Context, progress *entity.CourseSectionUnitProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func (m *MockLearningRepository) GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.CourseSectionUnitProgress, error) {
	args := m.Called(ctx, userID, unitID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CourseSectionUnitProgress), args.Error(1)
}

func (m *MockLearningRepository) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.CourseSectionUnitProgress, error) {
	args := m.Called(ctx, userID, sectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.CourseSectionUnitProgress), args.Error(1)
}

func (m *MockLearningRepository) UpsertUnitProgress(ctx context.Context, progress *entity.CourseSectionUnitProgress) error {
	args := m.Called(ctx, progress)
	return args.Error(0)
}

func TestLearningService_SaveCourseProgress(t *testing.T) {
	mockRepo := new(MockLearningRepository)
	mockMemoryService := new(MockMemoryService)
	service := NewLearningService(mockRepo, mockMemoryService)
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    uint
		courseID  uint
		status    string
		score     int
		mockSetup func()
		wantErr   bool
	}{
		{
			name:     "正常保存课程进度",
			userID:   1,
			courseID: 100,
			status:   "in_progress",
			score:    80,
			mockSetup: func() {
				mockRepo.On("SaveCourseProgress", ctx, mock.MatchedBy(func(p *entity.CourseLearningProgress) bool {
					return p.UserID == 1 && p.CourseID == 100 && p.Status == "in_progress" && p.Score == 80
				})).Return(nil)
			},
			wantErr: false,
		},
		{
			name:     "完成课程进度",
			userID:   1,
			courseID: 100,
			status:   "completed",
			score:    100,
			mockSetup: func() {
				mockRepo.On("SaveCourseProgress", ctx, mock.Anything).Run(func(args mock.Arguments) {
					p := args.Get(1).(*entity.CourseLearningProgress)
					p.UpdatedAt = time.Now()
				}).Return(nil)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			progress, err := service.SaveCourseProgress(ctx, entity.UID(tt.userID), tt.courseID, tt.status, tt.score)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.NotNil(t, progress)
			assert.Equal(t, entity.UID(tt.userID), progress.UserID)
			assert.Equal(t, tt.courseID, progress.CourseID)
			assert.Equal(t, tt.status, progress.Status)
			assert.Equal(t, tt.score, progress.Score)
			if tt.status == "completed" {
				assert.True(t, progress.UpdatedAt.After(progress.CreatedAt))
			}
		})
	}
}

func TestLearningService_GetCourseProgress(t *testing.T) {
	mockRepo := new(MockLearningRepository)
	mockMemoryService := new(MockMemoryService)
	service := NewLearningService(mockRepo, mockMemoryService)
	ctx := WithUserID(context.Background(), entity.UID(1))

	mockProgress := &entity.CourseLearningProgress{
		ID:       1,
		UserID:   1,
		CourseID: 100,
		Status:   "in_progress",
		Score:    80,
	}

	mockRepo.On("GetCourseProgress", ctx, uint(1), uint(100)).Return(mockProgress, nil)
	mockRepo.On("GetCourseProgress", ctx, uint(1), uint(999)).Return(nil, nil)

	t.Run("获取存在的课程进度", func(t *testing.T) {
		progress, err := service.GetCourseProgress(ctx, uint(100))
		require.NoError(t, err)
		assert.NotNil(t, progress)
		assert.Equal(t, mockProgress.ID, progress.ID)
		assert.Equal(t, mockProgress.Status, progress.Status)
	})

	t.Run("获取不存在的课程进度", func(t *testing.T) {
		progress, err := service.GetCourseProgress(ctx, uint(999))
		require.NoError(t, err)
		assert.Nil(t, progress)
	})
}

func TestLearningService_ListCourseProgress(t *testing.T) {
	mockRepo := new(MockLearningRepository)
	mockMemoryService := new(MockMemoryService)
	service := NewLearningService(mockRepo, mockMemoryService)
	ctx := context.Background()

	mockProgresses := []*entity.CourseLearningProgress{
		{
			ID:       1,
			UserID:   1,
			CourseID: 100,
			Status:   "completed",
			Score:    90,
		},
		{
			ID:       2,
			UserID:   1,
			CourseID: 101,
			Status:   "in_progress",
			Score:    60,
		},
	}

	mockRepo.On("ListCourseProgress", ctx, uint(1), 0, 10).Return(mockProgresses, int64(2), nil)

	t.Run("列出课程进度", func(t *testing.T) {
		progresses, total, err := service.ListCourseProgress(ctx, 1, 1, 10)
		require.NoError(t, err)
		assert.Equal(t, int64(2), total)
		assert.Len(t, progresses, 2)
		assert.Equal(t, uint(100), progresses[0].CourseID)
		assert.Equal(t, uint(101), progresses[1].CourseID)
	})
}
