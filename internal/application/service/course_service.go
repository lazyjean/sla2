package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// CourseService 课程服务
type CourseService struct {
	courseRepository        repository.CourseRepository
	courseSectionRepository repository.CourseSectionRepository
}

// NewCourseService 创建课程服务实例
func NewCourseService(
	courseRepository repository.CourseRepository,
	courseSectionRepository repository.CourseSectionRepository,
) *CourseService {
	return &CourseService{
		courseRepository:        courseRepository,
		courseSectionRepository: courseSectionRepository,
	}
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, title, description, coverURL, level string, tags []string) (*entity.Course, error) {
	course := &entity.Course{
		Title:       title,
		Description: description,
		CoverURL:    coverURL,
		Level:       level,
		Tags:        tags,
		Status:      "draft",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.courseRepository.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, id uint, title, description, coverURL, level string, tags []string, status string) (*entity.Course, error) {
	course, err := s.courseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	course.Title = title
	course.Description = description
	course.CoverURL = coverURL
	course.Level = level
	course.Tags = tags
	course.Status = status
	course.UpdatedAt = time.Now()

	if err := s.courseRepository.Update(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// GetCourse 获取课程详情
func (s *CourseService) GetCourse(ctx context.Context, id uint) (*entity.Course, error) {
	// 获取课程基本信息
	course, err := s.courseRepository.GetByID(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	// 获取课程章节
	sections, err := s.courseSectionRepository.ListByCourseID(ctx, course.ID)
	if err != nil {
		return nil, err
	}

	// 获取每个章节的单元信息
	for _, section := range sections {
		units, err := s.courseSectionRepository.ListUnitsBySectionID(ctx, section.ID)
		if err != nil {
			return nil, err
		}
		section.Units = units
	}

	course.Sections = sections
	return course, nil
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

	return s.courseRepository.List(ctx, offset, pageSize, filters)
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, id uint) error {
	return s.courseRepository.Delete(ctx, id)
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
	return s.courseRepository.Search(ctx, keyword, offset, pageSize, filters)
}

// CreateSection 创建课程章节
func (s *CourseService) CreateSection(ctx context.Context, courseID uint, title, desc string) (*entity.CourseSection, error) {
	// 检查课程是否存在
	course, err := s.courseRepository.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}

	// 获取当前课程的所有章节
	sections, err := s.courseSectionRepository.ListByCourseID(ctx, course.ID)
	if err != nil {
		return nil, err
	}

	// 计算新章节的顺序
	orderIndex := int32(0)
	if len(sections) > 0 {
		orderIndex = sections[len(sections)-1].OrderIndex + 1
	}

	// 创建新章节
	section := &entity.CourseSection{
		CourseID:   course.ID,
		Title:      title,
		Desc:       desc,
		OrderIndex: orderIndex,
		Status:     "enabled", // 默认启用
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.courseSectionRepository.Create(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

// UpdateSection 更新课程章节
func (s *CourseService) UpdateSection(ctx context.Context, id entity.CourseSectionID, title, desc string, orderIndex int32, status string) (*entity.CourseSection, error) {
	section, err := s.courseSectionRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	section.Title = title
	section.Desc = desc
	section.OrderIndex = orderIndex
	section.Status = status
	section.UpdatedAt = time.Now()

	if err := s.courseSectionRepository.Update(ctx, section); err != nil {
		return nil, err
	}

	return section, nil
}

// DeleteSection 删除课程章节
func (s *CourseService) DeleteSection(ctx context.Context, id entity.CourseSectionID) error {
	return s.courseSectionRepository.Delete(ctx, id)
}

// GetSection 获取课程章节
func (s *CourseService) GetSection(ctx context.Context, id entity.CourseSectionID) (*entity.CourseSection, error) {
	return s.courseSectionRepository.GetByID(ctx, id)
}

// ListSections 获取课程的所有章节
func (s *CourseService) ListSections(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error) {
	return s.courseSectionRepository.ListByCourseID(ctx, courseID)
}

// CreateUnit 创建课程单元
func (s *CourseService) CreateUnit(ctx context.Context, sectionID entity.CourseSectionID, title, desc string, questionIds []uint32, tags []string) (*entity.CourseSectionUnit, error) {
	// 将uint32数组转换为字符串数组
	questionIdStrs := make([]string, len(questionIds))
	for i, id := range questionIds {
		questionIdStrs[i] = strconv.FormatUint(uint64(id), 10)
	}

	unit := &entity.CourseSectionUnit{
		SectionID:   sectionID,
		Title:       title,
		Desc:        desc,
		QuestionIds: strings.Join(questionIdStrs, ","),
		Tags:        strings.Join(tags, ","),
		OrderIndex:  0, // 需要计算
		Status:      1, // 默认启用
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 获取当前章节的单元列表
	units, err := s.courseSectionRepository.ListUnitsBySectionID(ctx, sectionID)
	if err != nil {
		return nil, err
	}

	// 计算新单元的顺序
	if len(units) > 0 {
		unit.OrderIndex = units[len(units)-1].OrderIndex + 1
	}

	if err := s.courseSectionRepository.CreateUnit(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// UpdateUnit 更新课程单元
func (s *CourseService) UpdateUnit(ctx context.Context, id entity.CourseSectionUnitID, title, desc string, questionIds []uint32, tags []string, status int32) (*entity.CourseSectionUnit, error) {
	// 将uint32数组转换为字符串数组
	questionIdStrs := make([]string, len(questionIds))
	for i, id := range questionIds {
		questionIdStrs[i] = strconv.FormatUint(uint64(id), 10)
	}

	unit, err := s.courseSectionRepository.GetUnitByID(ctx, id)
	if err != nil {
		return nil, err
	}

	unit.Title = title
	unit.Desc = desc
	unit.QuestionIds = strings.Join(questionIdStrs, ",")
	unit.Tags = strings.Join(tags, ",")
	unit.Status = status
	unit.UpdatedAt = time.Now()

	if err := s.courseSectionRepository.UpdateUnit(ctx, unit); err != nil {
		return nil, err
	}

	return unit, nil
}

// DeleteUnit 删除课程单元
func (s *CourseService) DeleteUnit(ctx context.Context, id entity.CourseSectionUnitID) error {
	return s.courseSectionRepository.DeleteUnit(ctx, id)
}

// GetUnit 获取课程单元详情
func (s *CourseService) GetUnit(ctx context.Context, id entity.CourseSectionUnitID) (*entity.CourseSectionUnit, error) {
	return s.courseSectionRepository.GetUnitByID(ctx, id)
}
