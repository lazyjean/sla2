package learning

import (
	"context"
	"testing"
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
	pg "github.com/lazyjean/sla2/internal/infrastructure/persistence/postgres"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db              *gorm.DB
	learningRepo    repository.LearningRepository
	memoryUnitRepo  repository.MemoryUnitRepository
	wordRepo        repository.WordRepository
	learningService *service.LearningService
	memoryService   service.MemoryService
	grpcService     *LearningService
)

func SetupTestDB(t *testing.T) error {
	var cleanup func()
	db, cleanup = pg.SetupTestDB(t)
	t.Cleanup(cleanup)

	// 清理测试数据
	if err := db.Exec("TRUNCATE TABLE course_learning_progresses CASCADE").Error; err != nil {
		return err
	}
	if err := db.Exec("TRUNCATE TABLE course_section_progresses CASCADE").Error; err != nil {
		return err
	}
	if err := db.Exec("TRUNCATE TABLE course_section_unit_progresses CASCADE").Error; err != nil {
		return err
	}
	if err := db.Exec("TRUNCATE TABLE memory_units CASCADE").Error; err != nil {
		return err
	}

	learningRepo = pg.NewLearningRepository(db)
	memoryUnitRepo = pg.NewMemoryUnitRepository(db)
	wordRepo = pg.NewWordRepository(db)
	learningService = service.NewLearningService(learningRepo, memoryService)
	memoryService = service.NewMemoryService(wordRepo, memoryUnitRepo)
	grpcService = NewLearningService(learningService, memoryService)

	return nil
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

func TestRecordLearningResult(t *testing.T) {
	if err := SetupTestDB(t); err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}

	// 准备测试数据
	unit := &entity.MemoryUnit{
		UserID:             1,
		Type:               entity.MemoryUnitTypeWord,
		ContentID:          1,
		MasteryLevel:       entity.MasteryLevelBeginner,
		ReviewCount:        0,
		StudyDuration:      0,
		RetentionRate:      0,
		ConsecutiveCorrect: 0,
		ConsecutiveWrong:   0,
	}
	if err := memoryUnitRepo.Create(context.Background(), unit); err != nil {
		t.Fatalf("Failed to save test data: %v", err)
	}
	t.Logf("Created memory unit: ID=%d", unit.ID)

	// 记录第一次学习结果（正确）
	req1 := &pb.RecordLearningResultRequest{
		MemoryUnitId: unit.ID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 100,
		UserNotes:    []string{"note1"},
	}
	resp1, err := grpcService.RecordLearningResult(context.Background(), req1)
	if err != nil {
		t.Fatalf("Failed to record first learning result: %v", err)
	}
	assert.NotNil(t, resp1)

	// 记录第二次学习结果（错误）
	req2 := &pb.RecordLearningResultRequest{
		MemoryUnitId: unit.ID,
		Result:       pb.ReviewResult_REVIEW_RESULT_WRONG,
		ResponseTime: 150,
		UserNotes:    []string{"note2"},
	}
	resp2, err := grpcService.RecordLearningResult(context.Background(), req2)
	if err != nil {
		t.Fatalf("Failed to record second learning result: %v", err)
	}
	assert.NotNil(t, resp2)

	// 记录第三次学习结果（正确）
	req3 := &pb.RecordLearningResultRequest{
		MemoryUnitId: unit.ID,
		Result:       pb.ReviewResult_REVIEW_RESULT_CORRECT,
		ResponseTime: 80,
		UserNotes:    []string{"note3"},
	}
	resp3, err := grpcService.RecordLearningResult(context.Background(), req3)
	if err != nil {
		t.Fatalf("Failed to record third learning result: %v", err)
	}
	assert.NotNil(t, resp3)

	// 获取最终状态
	getReq := &pb.GetMemoryStatusRequest{
		MemoryUnitId: unit.ID,
	}
	getResp, err := grpcService.GetMemoryStatus(context.Background(), getReq)
	if err != nil {
		t.Fatalf("Failed to get final status: %v", err)
	}
	t.Logf("Final status: ReviewCount=%d, ConsecutiveCorrect=%d, ConsecutiveWrong=%d, RetentionRate=%f",
		getResp.Status.ReviewCount, getResp.Status.ConsecutiveCorrect, getResp.Status.ConsecutiveWrong,
		getResp.Status.RetentionRate)
}
