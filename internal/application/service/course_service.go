package service

import (
	"context"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// CourseService 课程服务
type CourseService struct {
	courseRepo repository.CourseRepository
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseRepo repository.CourseRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
	}
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, title, description, coverURL, level string, duration int, tags []string) (*entity.Course, error) {
	course := &entity.Course{
		Title:       title,
		Description: description,
		CoverURL:    coverURL,
		Level:       level,
		Duration:    duration,
		Tags:        tags,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.courseRepo.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, id uint, title, description, coverURL, level string, duration int, tags []string, status string) (*entity.Course, error) {
	course, err := s.courseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	course.Title = title
	course.Description = description
	course.CoverURL = coverURL
	course.Level = level
	course.Duration = duration
	course.Tags = tags
	course.Status = status
	course.UpdatedAt = time.Now()

	if err := s.courseRepo.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// GetCourse 获取课程详情
func (s *CourseService) GetCourse(ctx context.Context, id uint) (*entity.Course, error) {
	return s.courseRepo.GetByID(ctx, uint(id))
}

// ListCourses 获取课程列表
func (s *CourseService) ListCourses(ctx context.Context, page, pageSize int, level uint, tags []string, status string) ([]*entity.Course, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != 0 {
		filters["level"] = level
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}
	if status != "" {
		filters["status"] = status
	}

	return s.courseRepo.List(ctx, offset, pageSize, filters)
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, id uint) error {
	return s.courseRepo.Delete(ctx, id)
}

// SearchCourse 搜索课程
func (s *CourseService) SearchCourse(ctx context.Context, keyword string, page, pageSize int, level uint, tags []string) ([]*entity.Course, int64, error) {
	offset := (page - 1) * pageSize
	filters := make(map[string]interface{})
	if level != 0 {
		filters["level"] = level
	}
	if len(tags) > 0 {
		filters["tags"] = tags
	}
	return s.courseRepo.Search(ctx, keyword, offset, pageSize, filters)
}
