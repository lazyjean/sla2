package postgres

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors"
	"github.com/lazyjean/sla2/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LearningRepository struct {
	db *gorm.DB
}

func NewLearningRepository(db *gorm.DB) repository.LearningRepository {
	return &LearningRepository{
		db: db,
	}
}

// SaveCourseProgress 保存课程进度
func (r *LearningRepository) SaveCourseProgress(ctx context.Context, progress *entity.CourseLearningProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetCourseProgress 获取课程进度
func (r *LearningRepository) GetCourseProgress(ctx context.Context, userID, courseID uint) (*entity.CourseLearningProgress, error) {
	var progress entity.CourseLearningProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListCourseProgress 列出用户的课程进度
func (r *LearningRepository) ListCourseProgress(ctx context.Context, userID uint, offset, limit int) ([]*entity.CourseLearningProgress, int64, error) {
	var progresses []*entity.CourseLearningProgress
	var total int64

	err := r.db.WithContext(ctx).Model(&entity.CourseLearningProgress{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	err = r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&progresses).Error
	if err != nil {
		return nil, 0, domainErrors.ErrFailedToQuery
	}

	return progresses, total, nil
}

// SaveSectionProgress 保存章节进度
func (r *LearningRepository) SaveSectionProgress(ctx context.Context, progress *entity.CourseSectionProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetSectionProgress 获取章节进度
func (r *LearningRepository) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.CourseSectionProgress, error) {
	var progress entity.CourseSectionProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND section_id = ?", userID, sectionID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListSectionProgress 列出章节进度
func (r *LearningRepository) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.CourseSectionProgress, error) {
	var progresses []*entity.CourseSectionProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("section_id ASC").
		Find(&progresses).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progresses, nil
}

// UpsertUnitProgress 保存或更新单元进度
func (r *LearningRepository) UpsertUnitProgress(ctx context.Context, progress *entity.CourseSectionUnitProgress) error {
	result := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "section_id"},
			{Name: "unit_id"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"status":         progress.Status,
			"complete_count": gorm.Expr("course_section_unit_progresses.complete_count + 1"),
			"updated_at":     time.Now(),
		}),
	}).Create(progress)

	if result.Error != nil {
		return domainErrors.ErrFailedToSave
	}

	return nil
}

// ListUnitProgress 获取章节的单元学习进度列表
func (r *LearningRepository) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.CourseSectionUnitProgress, error) {
	var progress []*entity.CourseSectionUnitProgress
	if err := r.db.WithContext(ctx).Where("user_id = ? AND section_id = ?", userID, sectionID).Find(&progress).Error; err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progress, nil
}
