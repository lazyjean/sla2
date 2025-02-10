package postgres

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	domainErrors "github.com/lazyjean/sla2/internal/domain/errors"
	"gorm.io/gorm"
)

type LearningRepository struct {
	db *gorm.DB
}

func NewLearningRepository(db *gorm.DB) *LearningRepository {
	return &LearningRepository{db: db}
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
func (r *LearningRepository) SaveSectionProgress(ctx context.Context, progress *entity.SectionProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetSectionProgress 获取章节进度
func (r *LearningRepository) GetSectionProgress(ctx context.Context, userID, sectionID uint) (*entity.SectionProgress, error) {
	var progress entity.SectionProgress
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
func (r *LearningRepository) ListSectionProgress(ctx context.Context, userID, courseID uint) ([]*entity.SectionProgress, error) {
	var progresses []*entity.SectionProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Order("section_id ASC").
		Find(&progresses).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progresses, nil
}

// SaveUnitProgress 保存单元进度
func (r *LearningRepository) SaveUnitProgress(ctx context.Context, progress *entity.UnitProgress) error {
	if err := r.db.WithContext(ctx).Save(progress).Error; err != nil {
		return domainErrors.ErrFailedToSave
	}
	return nil
}

// GetUnitProgress 获取单元进度
func (r *LearningRepository) GetUnitProgress(ctx context.Context, userID, unitID uint) (*entity.UnitProgress, error) {
	var progress entity.UnitProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND unit_id = ?", userID, unitID).
		First(&progress).Error

	if err == gorm.ErrRecordNotFound {
		return nil, domainErrors.ErrProgressNotFound
	}
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return &progress, nil
}

// ListUnitProgress 列出单元进度
func (r *LearningRepository) ListUnitProgress(ctx context.Context, userID, sectionID uint) ([]*entity.UnitProgress, error) {
	var progresses []*entity.UnitProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND section_id = ?", userID, sectionID).
		Order("unit_id ASC").
		Find(&progresses).Error
	if err != nil {
		return nil, domainErrors.ErrFailedToQuery
	}
	return progresses, nil
}
