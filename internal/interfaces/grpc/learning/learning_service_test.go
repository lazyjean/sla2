package learning

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	pg "github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/lazyjean/sla2/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	grpcService    *LearningService
	memoryUnitRepo repository.MemoryUnitRepository
	learningRepo   repository.LearningRepository
	db             *gorm.DB
)

// ensureLogger 初始化基础 logger (如果尚未初始化)
func ensureLogger() {
	if logger.Log == nil {
		logger.InitBaseLogger()
		logger.Log.Info("Initialized base logger for testing.")
	}
}

// SetupTestDB 设置测试数据库并初始化必要的 repo
func SetupTestDB(t *testing.T) error {
	ensureLogger() // Initialize logger here

	var cleanup func()
	db, cleanup = pg.SetupTestDB(t) // Assuming pg.SetupTestDB is the correct package and function
	t.Cleanup(cleanup)

	// Run migrations to ensure schema matches entities
	err := db.AutoMigrate(
		&entity.MemoryUnit{}, // Add other entities if needed for tests in this package
		// &entity.CourseLearningProgress{}, // Example
		// &entity.CourseSectionProgress{}, // Example
		// &entity.CourseSectionUnitProgress{}, // Example
	)
	if err != nil {
		logger.Log.Error("AutoMigrate failed", zap.Error(err))
		return fmt.Errorf("automigrate failed: %w", err)
	}
	logger.Log.Info("AutoMigrate completed for learning tests.")

	// 清理测试数据 (Truncate after migration)
	logger.Log.Info("Truncating memory_units table...")
	if err := db.Exec("TRUNCATE TABLE memory_units CASCADE").Error; err != nil {
		logger.Log.Error("Failed to truncate memory_units", zap.Error(err))
		return err
	}
	// ... add other TRUNCATE statements if needed ...

	// Initialize repositories used in tests
	learningRepo = pg.NewLearningRepository(db)
	memoryUnitRepo = pg.NewMemoryUnitRepository(db)

	// Initialize services needed for tests (can be done here or in TestMain/specific tests)
	memoryService := service.NewMemoryService(nil, memoryUnitRepo)             // wordRepo is nil, adjust if needed
	learningService := service.NewLearningService(learningRepo, memoryService) // Use initialized learningRepo
	grpcService = NewLearningService(learningService, memoryService)

	return nil
}

func TestMain(m *testing.M) {
	// TestMain can be simplified or removed if SetupTestDB handles all setup
	// Or, it can handle setup not specific to individual tests
	code := m.Run()
	os.Exit(code)
}

func TestGetCourseProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseLearningProgress{
		UserID:   1,
		CourseID: 1,
		Status:   "completed",
		Score:    100,
		Progress: 100,
	}
	if err := learningRepo.SaveCourseProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceGetCourseProgressRequest{
		CourseId: 1,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.GetCourseProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get course progress: %v", err)
	}

	// 验证结果
	assert.Equal(t, float32(100.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(0), resp.Progress.TotalItems)
}

func TestGetSectionProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  1,
		SectionID: 1,
		Status:    "completed",
		Progress:  100,
	}
	if err := learningRepo.SaveSectionProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceGetSectionProgressRequest{
		SectionId: 1,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.GetSectionProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to get section progress: %v", err)
	}

	// 验证结果
	assert.Equal(t, float32(100.0), resp.Progress.Progress)
	assert.Equal(t, uint32(0), resp.Progress.CompletedItems)
	assert.Equal(t, uint32(0), resp.Progress.TotalItems)
}

func TestUpdateUnitProgress(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	progress := &entity.CourseSectionUnitProgress{
		UserID:        1,
		SectionID:     1,
		UnitID:        1,
		Status:        "in_progress",
		CompleteCount: 0,
	}
	if err := learningRepo.UpsertUnitProgress(context.Background(), progress); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.LearningServiceUpdateUnitProgressRequest{
		UnitId:    1,
		SectionId: 1,
		Completed: true,
	}
	ctx := service.WithUserID(context.Background(), entity.UID(1))
	resp, err := grpcService.UpdateUnitProgress(ctx, req)
	if err != nil {
		t.Fatalf("Failed to update unit progress: %v", err)
	}

	// 验证结果
	assert.NotNil(t, resp)

	// 验证数据库中的更新
	unitProgresses, err := learningRepo.ListUnitProgress(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Failed to get updated progress: %v", err)
	}
	assert.Equal(t, 1, len(unitProgresses))
	assert.Equal(t, "completed", unitProgresses[0].Status)
	assert.Equal(t, uint(1), unitProgresses[0].CompleteCount)
}

func TestUpdateMemoryStatus(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	now := time.Now()
	unit := &entity.MemoryUnit{
		UserID:             1,
		Type:               entity.MemoryUnitTypeWord,
		ContentID:          1,
		CreatedAt:          now,
		UpdatedAt:          now,
		MasteryLevel:       entity.MasteryLevelBeginner,
		ReviewCount:        0,
		NextReviewAt:       now.Add(time.Hour),
		LastReviewAt:       now,
		StudyDuration:      0,
		RetentionRate:      0,
		ConsecutiveCorrect: 0,
		ConsecutiveWrong:   0,
	}
	if err := memoryUnitRepo.Create(context.Background(), unit); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}

	// 执行测试
	req := &pb.UpdateMemoryStatusRequest{
		MemoryUnitId:  unit.ID,
		MasteryLevel:  pb.MasteryLevel_MASTERY_LEVEL_FAMILIAR,
		StudyDuration: 100,
	}
	resp, err := grpcService.UpdateMemoryStatus(context.Background(), req)
	if err != nil {
		t.Fatalf("Failed to update memory status: %v", err)
	}

	// 验证结果
	assert.NotNil(t, resp)

	// 验证数据库中的更新
	updatedUnit, err := memoryUnitRepo.GetByID(context.Background(), unit.ID)
	if err != nil {
		t.Fatalf("Failed to get updated unit: %v", err)
	}
	assert.Equal(t, entity.MasteryLevelFamiliar, updatedUnit.MasteryLevel)
	assert.Equal(t, uint32(100), updatedUnit.StudyDuration)
}

func TestReviewWord(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ctx := context.Background()
	wordID := uint32(99999) // Use an ID unlikely to exist

	// 1. First review (Correct) - Should create the MemoryUnit
	req1 := &pb.ReviewWordRequest{
		WordId:       wordID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 1500, // 1.5 seconds
	}
	resp1, err := grpcService.ReviewWord(ctx, req1)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)

	// Verify creation and initial state in DB
	unit1, err := memoryUnitRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeWord, wordID)
	assert.NoError(t, err)
	assert.NotNil(t, unit1, "MemoryUnit should have been created")
	t.Logf("[Review 1] LastReviewAt: %s, NextReviewAt: %s", unit1.LastReviewAt, unit1.NextReviewAt)
	assert.True(t, unit1.NextReviewAt.After(unit1.LastReviewAt), "Next review time must be after last review")
	initialNextReviewAt := unit1.NextReviewAt

	// 2. Second review (Wrong)
	time.Sleep(10 * time.Millisecond)              // Ensure timestamp changes
	lastReviewAtBeforeUpdate := unit1.LastReviewAt // Record time before the update
	req2 := &pb.ReviewWordRequest{
		WordId:       wordID,
		Result:       pb.ReviewResult_REVIEW_RESULT_WRONG,
		ResponseTime: 2100, // 2.1 seconds
	}
	resp2, err := grpcService.ReviewWord(ctx, req2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	// Verify updated state in DB
	unit2, err := memoryUnitRepo.GetByID(ctx, unit1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, unit2)
	t.Logf("[Review 2] LastReviewAt: %s, NextReviewAt: %s", unit2.LastReviewAt, unit2.NextReviewAt)
	assert.True(t, unit2.LastReviewAt.After(lastReviewAtBeforeUpdate), "Last review time should be updated")
	assert.True(t, unit2.NextReviewAt.After(unit2.LastReviewAt), "Next review time must be after last review")
	assert.True(t, unit2.NextReviewAt.Before(initialNextReviewAt.Add(time.Hour)), "Next review time after wrong should be sooner than initial next review plus buffer")
}

func TestReviewHanChar(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	ctx := context.Background()
	hanCharID := uint32(88888) // Use an ID unlikely to exist

	// 1. First review (Correct) - Should create the MemoryUnit
	req1 := &pb.ReviewHanCharRequest{
		HanCharId:    hanCharID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 1800, // 1.8 seconds
	}
	resp1, err := grpcService.ReviewHanChar(ctx, req1)
	assert.NoError(t, err)
	assert.NotNil(t, resp1)

	// Verify creation and initial state in DB
	unit1, err := memoryUnitRepo.GetByTypeAndContentID(ctx, entity.MemoryUnitTypeHanChar, hanCharID)
	assert.NoError(t, err)
	assert.NotNil(t, unit1, "MemoryUnit should have been created")
	assert.Equal(t, uint32(1), unit1.ReviewCount)
	assert.Equal(t, entity.MasteryLevelBeginner, unit1.MasteryLevel)
	assert.Equal(t, uint32(1), unit1.ConsecutiveCorrect)
	assert.Equal(t, uint32(0), unit1.ConsecutiveWrong)
	assert.Equal(t, uint32(1), unit1.StudyDuration) // 1800ms / 1000 = 1s
	assert.True(t, unit1.NextReviewAt.After(unit1.LastReviewAt), "Next review time must be after last review")
	initialNextReviewAt2 := unit1.NextReviewAt

	// 2. Second review (Correct)
	time.Sleep(10 * time.Millisecond)
	lastReviewAtBeforeUpdate2 := unit1.LastReviewAt
	req2 := &pb.ReviewHanCharRequest{
		HanCharId:    hanCharID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 1200, // 1.2 seconds
	}
	resp2, err := grpcService.ReviewHanChar(ctx, req2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	// Verify updated state in DB
	unit2, err := memoryUnitRepo.GetByID(ctx, unit1.ID)
	assert.NoError(t, err)
	assert.NotNil(t, unit2)
	assert.Equal(t, uint32(2), unit2.ReviewCount)
	assert.Equal(t, entity.MasteryLevelBeginner, unit2.MasteryLevel) // 2 correct -> still Beginner (needs 3 for Familiar)
	assert.Equal(t, uint32(2), unit2.ConsecutiveCorrect)
	assert.Equal(t, uint32(0), unit2.ConsecutiveWrong)
	assert.Equal(t, uint32(1+1), unit2.StudyDuration) // 1s + 1s = 2s
	assert.True(t, unit2.LastReviewAt.After(lastReviewAtBeforeUpdate2), "Last review time should be updated")
	assert.True(t, unit2.NextReviewAt.After(unit2.LastReviewAt), "Next review time must be after last review")
	assert.True(t, unit2.NextReviewAt.After(initialNextReviewAt2), "Next review time should be further out after second correct")
}
