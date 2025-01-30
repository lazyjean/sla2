package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/lazyjean/sla2/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLearningRepository_CourseProgress(t *testing.T) {
	db := setupTestDB(t)
	repo := NewLearningRepository(db)
	ctx := context.Background()

	// 创建测试数据
	progress := &entity.CourseLearningProgress{
		UserID:    1,
		CourseID:  100,
		Status:    "in_progress",
		Score:     80,
		StartedAt: time.Now(),
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
	progress := &entity.SectionProgress{
		UserID:    1,
		CourseID:  100,
		SectionID: 1001,
		Status:    "completed",
		Progress:  100.0,
		StartedAt: time.Now(),
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
	progress := &entity.UnitProgress{
		UserID:    1,
		UnitID:    300,
		Status:    "in_progress",
		Progress:  75.0,
		StartedAt: time.Now(),
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
