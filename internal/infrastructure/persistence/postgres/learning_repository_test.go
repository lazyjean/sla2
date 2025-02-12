package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLearningRepository_CourseProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseLearningProgress{
		UserID:   1,
		CourseID: 100,
		Status:   "in_progress",
		Score:    80,
	}

	// 测试保存
	err := repo.SaveCourseProgress(ctx, progress)
	require.NoError(t, err)
	assert.NotZero(t, progress.ID)

	// 测试查询
	found, err := repo.GetCourseProgress(ctx, progress.UserID, progress.CourseID)
	require.NoError(t, err)
	assert.Equal(t, progress.Status, found.Status)
	assert.Equal(t, progress.Score, found.Score)

	// 测试列表
	list, total, err := repo.ListCourseProgress(ctx, progress.UserID, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
}

func TestLearningRepository_SectionProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  100,
		SectionID: 1001,
		Status:    "completed",
		Progress:  100.0,
	}

	// 测试保存
	err := repo.SaveSectionProgress(ctx, progress)
	require.NoError(t, err)

	// 测试查询
	found, err := repo.GetSectionProgress(ctx, progress.UserID, progress.SectionID)
	require.NoError(t, err)
	assert.Equal(t, progress.Status, found.Status)
	assert.Equal(t, progress.Progress, found.Progress)

	// 测试列表
	list, err := repo.ListSectionProgress(ctx, progress.UserID, progress.CourseID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestLearningRepository_UnitProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseSectionUnitProgress{
		UserID:   1,
		UnitID:   300,
		Status:   "in_progress",
		Progress: 75.0,
	}

	// 测试保存
	err := repo.SaveUnitProgress(ctx, progress)
	require.NoError(t, err)
	assert.NotZero(t, progress.ID)

	// 测试查询
	found, err := repo.GetUnitProgress(ctx, progress.UserID, progress.UnitID)
	require.NoError(t, err)
	assert.Equal(t, progress.Status, found.Status)
	assert.Equal(t, progress.Progress, found.Progress)

	// 测试列表
	list, err := repo.ListUnitProgress(ctx, progress.UserID, progress.SectionID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestLearningRepository_SaveCourseProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseLearningProgress{
		UserID:   1,
		CourseID: 100,
		Status:   "in_progress",
		Score:    80,
	}

	// 测试保存
	err := repo.SaveCourseProgress(ctx, progress)
	require.NoError(t, err)
	assert.NotZero(t, progress.ID)

	// 测试查询
	found, err := repo.GetCourseProgress(ctx, progress.UserID, progress.CourseID)
	require.NoError(t, err)

	// 添加时间校验前的等待
	time.Sleep(1 * time.Millisecond)

	// 更新记录触发时间变更
	progress.Status = "completed"
	err = repo.SaveCourseProgress(ctx, progress)
	require.NoError(t, err)

	// 重新获取更新后的记录
	updated, err := repo.GetCourseProgress(ctx, progress.UserID, progress.CourseID)
	require.NoError(t, err)

	// 验证更新时间变化
	assert.True(t, updated.UpdatedAt.After(found.UpdatedAt))

	// 测试列表
	list, total, err := repo.ListCourseProgress(ctx, progress.UserID, 0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, list, 1)
}

func TestSaveSectionProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  100,
		SectionID: 1001,
		Status:    "completed",
		Progress:  100.0,
	}

	err := repo.SaveSectionProgress(ctx, progress)
	require.NoError(t, err)
}

func TestGetSectionProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 先保存数据
	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  100,
		SectionID: 1001,
		Status:    "completed",
		Progress:  100.0,
	}
	err := repo.SaveSectionProgress(ctx, progress)
	require.NoError(t, err)

	// 测试查询
	found, err := repo.GetSectionProgress(ctx, progress.UserID, progress.SectionID)
	require.NoError(t, err)
	assert.Equal(t, progress.Status, found.Status)
	assert.Equal(t, progress.Progress, found.Progress)
}

func TestListSectionProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 先保存数据
	progress := &entity.CourseSectionProgress{
		UserID:    1,
		CourseID:  100,
		SectionID: 1001,
		Status:    "completed",
		Progress:  100.0,
	}
	err := repo.SaveSectionProgress(ctx, progress)
	require.NoError(t, err)

	// 测试列表
	list, err := repo.ListSectionProgress(ctx, progress.UserID, progress.CourseID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}

func TestLearningRepository_SaveUnitProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseSectionUnitProgress{
		UserID:   1,
		UnitID:   300,
		Status:   "in_progress",
		Progress: 75.0,
	}

	// 测试保存
	err := repo.SaveUnitProgress(ctx, progress)
	require.NoError(t, err)
	assert.NotZero(t, progress.ID)

	// 测试查询
	found, err := repo.GetUnitProgress(ctx, progress.UserID, progress.UnitID)
	require.NoError(t, err)
	assert.Equal(t, progress.Status, found.Status)
	assert.Equal(t, progress.Progress, found.Progress)

	// 测试列表
	list, err := repo.ListUnitProgress(ctx, progress.UserID, progress.SectionID)
	require.NoError(t, err)
	assert.Len(t, list, 1)
}
