package service

import (
	"context"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// domain course service
type CourseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(courseRepo repository.CourseRepository) *CourseService {
	return &CourseService{courseRepo: courseRepo}
}

func (s *CourseService) CreateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error) {
	return s.courseRepo.Create(ctx, course)
}

func (s *CourseService) UpdateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error) {
	return s.courseRepo.Update(ctx, course)
}

func (s *CourseService) GetCourse(ctx context.Context, id uint) (*entity.Course, error) {
	return s.courseRepo.GetByID(ctx, id)
}

func (s *CourseService) ListCourses(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]*entity.Course, int64, error) {
	return s.courseRepo.List(ctx, page, pageSize, filters)
}

func (s *CourseService) DeleteCourse(ctx context.Context, id uint) error {
	return s.courseRepo.Delete(ctx, id)
}
