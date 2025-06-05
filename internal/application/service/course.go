package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/domain/repository"
)

// CourseService 课程服务接口
type CourseService interface {
	// CreateCourse 创建课程
	CreateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error)
	// UpdateCourse 更新课程
	UpdateCourse(ctx context.Context, id entity.CourseID, course *entity.Course) error
	// GetCourse 获取课程详情
	GetCourse(ctx context.Context, id entity.CourseID) (*entity.Course, error)
	// ListCourses 获取课程列表
	ListCourses(ctx context.Context, input *dto.CourseListInput) ([]*entity.Course, int64, error)
	// DeleteCourse 删除课程
	DeleteCourse(ctx context.Context, id entity.CourseID) error
	// SearchCourse 搜索课程
	SearchCourse(ctx context.Context, input *dto.CourseSearchInput) ([]*entity.Course, int64, error)
	// CreateSection 创建课程章节
	CreateSection(ctx context.Context, courseID entity.CourseID, orderIndex int32, title, desc string) (entity.CourseSectionID, error)
	// UpdateSection 更新课程章节
	UpdateSection(ctx context.Context, id entity.CourseSectionID, title, desc string, orderIndex int32, status string) error
	// DeleteSection 删除课程章节
	DeleteSection(ctx context.Context, id entity.CourseSectionID) error
	// GetSection 获取课程章节
	GetSection(ctx context.Context, id entity.CourseSectionID) (*entity.CourseSection, error)
	// ListSections 获取课程的所有章节
	ListSections(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error)
	// CreateUnit 创建课程单元
	CreateUnit(ctx context.Context, sectionID entity.CourseSectionID, unit *entity.CourseSectionUnit) (*entity.CourseSectionUnit, error)
	// UpdateUnit 更新课程单元
	UpdateUnit(ctx context.Context, id entity.CourseSectionUnitID, unit *entity.CourseSectionUnit) error
	// DeleteUnit 删除课程单元
	DeleteUnit(ctx context.Context, id entity.CourseSectionUnitID) error
	// GetUnit 获取课程单元详情
	GetUnit(ctx context.Context, id entity.CourseSectionUnitID) (*entity.CourseSectionUnit, error)
	// BatchCreateCourse 批量创建课程
	BatchCreateCourse(ctx context.Context, courses []*entity.BatchCreateCourseInput) ([]entity.CourseID, error)
}

// courseService 课程服务实现
type courseService struct {
	courseRepository            repository.CourseRepository
	courseSectionRepository     repository.CourseSectionRepository
	courseSectionUnitRepository repository.CourseSectionUnitRepository
}

// NewCourseService 创建课程服务实例
func NewCourseService(
	courseRepository repository.CourseRepository,
	courseSectionRepository repository.CourseSectionRepository,
	courseSectionUnitRepository repository.CourseSectionUnitRepository,
) CourseService {
	return &courseService{
		courseRepository:            courseRepository,
		courseSectionRepository:     courseSectionRepository,
		courseSectionUnitRepository: courseSectionUnitRepository,
	}
}

// CreateCourse 创建课程
func (s *courseService) CreateCourse(ctx context.Context, course *entity.Course) (*entity.Course, error) {
	if err := s.courseRepository.Create(ctx, course); err != nil {
		return nil, err
	}
	return course, nil
}

// UpdateCourse 更新课程
func (s *courseService) UpdateCourse(ctx context.Context, id entity.CourseID, course *entity.Course) error {
	course.ID = id
	return s.courseRepository.Update(ctx, course)
}

// GetCourse 获取课程详情
func (s *courseService) GetCourse(ctx context.Context, id entity.CourseID) (*entity.Course, error) {
	course, err := s.courseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取课程章节
	sections, err := s.courseSectionRepository.ListByCourseID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 获取每个章节的单元信息
	for _, section := range sections {
		units, err := s.courseSectionUnitRepository.ListBySectionID(ctx, section.ID)
		if err != nil {
			return nil, err
		}
		section.Units = units
	}

	course.Sections = sections
	return course, nil
}

// ListCourses 获取课程列表
func (s *courseService) ListCourses(ctx context.Context, input *dto.CourseListInput) ([]*entity.Course, int64, error) {
	offset := (input.Page - 1) * input.PageSize

	filters := make(map[string]interface{})
	if input.Level != 0 {
		filters["level"] = input.Level
	}
	if input.Category != "" {
		filters["category"] = input.Category
	}
	if len(input.Tags) > 0 {
		filters["tags"] = input.Tags
	}
	if input.Status != "" {
		filters["status"] = input.Status
	}

	return s.courseRepository.ListWithFilters(ctx, offset, input.PageSize, filters)
}

// DeleteCourse 删除课程
func (s *courseService) DeleteCourse(ctx context.Context, id entity.CourseID) error {
	return s.courseRepository.Delete(ctx, id)
}

// SearchCourse 搜索课程
func (s *courseService) SearchCourse(ctx context.Context, input *dto.CourseSearchInput) ([]*entity.Course, int64, error) {
	offset := (input.Page - 1) * input.PageSize
	filters := make(map[string]interface{})
	if input.Level != 0 {
		filters["level"] = input.Level
	}
	if input.Category != "" {
		filters["category"] = input.Category
	}
	if len(input.Tags) > 0 {
		filters["tags"] = input.Tags
	}
	return s.courseRepository.Search(ctx, input.Keyword, offset, input.PageSize, filters)
}

// CreateSection 创建课程章节
func (s *courseService) CreateSection(ctx context.Context, courseID entity.CourseID, orderIndex int32, title, desc string) (entity.CourseSectionID, error) {
	// 创建新章节
	section := &entity.CourseSection{
		CourseID:   courseID,
		Title:      title,
		Desc:       desc,
		OrderIndex: orderIndex,
		Status:     "enabled", // 默认启用
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.courseSectionRepository.Create(ctx, section); err != nil {
		return 0, err
	}

	return section.ID, nil
}

// UpdateSection 更新课程章节
func (s *courseService) UpdateSection(ctx context.Context, id entity.CourseSectionID, title, desc string, orderIndex int32, status string) error {
	section := &entity.CourseSection{
		ID:         id,
		Title:      title,
		Desc:       desc,
		Status:     status,
		OrderIndex: orderIndex,
		UpdatedAt:  time.Now(),
	}
	if err := s.courseSectionRepository.Update(ctx, section); err != nil {
		return err
	}

	return nil
}

// DeleteSection 删除课程章节
func (s *courseService) DeleteSection(ctx context.Context, id entity.CourseSectionID) error {
	return s.courseSectionRepository.Delete(ctx, id)
}

// GetSection 获取课程章节
func (s *courseService) GetSection(ctx context.Context, id entity.CourseSectionID) (*entity.CourseSection, error) {
	return s.courseSectionRepository.GetByID(ctx, id)
}

// ListSections 获取课程的所有章节
func (s *courseService) ListSections(ctx context.Context, courseID entity.CourseID) ([]*entity.CourseSection, error) {
	return s.courseSectionRepository.ListByCourseID(ctx, courseID)
}

// CreateUnit 创建课程单元
func (s *courseService) CreateUnit(ctx context.Context, sectionID entity.CourseSectionID, unit *entity.CourseSectionUnit) (*entity.CourseSectionUnit, error) {
	if err := s.courseSectionUnitRepository.Create(ctx, unit); err != nil {
		return nil, err
	}
	return unit, nil
}

// UpdateUnit 更新课程单元
func (s *courseService) UpdateUnit(ctx context.Context, id entity.CourseSectionUnitID, unit *entity.CourseSectionUnit) error {
	return s.courseSectionUnitRepository.Update(ctx, unit)
}

// DeleteUnit 删除课程单元
func (s *courseService) DeleteUnit(ctx context.Context, id entity.CourseSectionUnitID) error {
	return s.courseSectionUnitRepository.Delete(ctx, id)
}

// GetUnit 获取课程单元详情
func (s *courseService) GetUnit(ctx context.Context, id entity.CourseSectionUnitID) (*entity.CourseSectionUnit, error) {
	return s.courseSectionUnitRepository.GetByID(ctx, id)
}

// BatchCreateCourse 批量创建课程
func (s *courseService) BatchCreateCourse(ctx context.Context, courses []*entity.BatchCreateCourseInput) ([]entity.CourseID, error) {
	// 用于存储创建的课程ID列表
	var courseIds []entity.CourseID

	// 遍历处理所有课程
	for _, courseData := range courses {
		// 1. 创建课程
		course := &entity.Course{
			Title:          courseData.Title,
			Description:    courseData.Description,
			CoverURL:       courseData.CoverURL,
			Level:          courseData.Level,
			Category:       courseData.Category,
			Tags:           courseData.Tags,
			Status:         "draft",
			Prompt:         courseData.Prompt,
			Resources:      courseData.Resources,
			RecommendedAge: courseData.RecommendedAge,
			StudyPlan:      courseData.StudyPlan,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.courseRepository.Create(ctx, course); err != nil {
			return nil, err
		}

		// 添加课程ID到结果列表
		courseIds = append(courseIds, course.ID)

		// 2. 创建章节和单元
		for _, section := range courseData.Sections {
			// 创建章节
			courseSection := &entity.CourseSection{
				CourseID:   course.ID,
				Title:      section.Title,
				Desc:       section.Desc,
				OrderIndex: section.OrderIndex,
				Status:     "enabled",
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := s.courseSectionRepository.Create(ctx, courseSection); err != nil {
				return nil, err
			}

			// 创建章节下的单元
			for _, unit := range section.Units {
				// 处理问题ID
				var questionIdsStr string
				if len(unit.QuestionIds) > 0 {
					// 将 uint32 数组转换为字符串
					strIds := make([]string, len(unit.QuestionIds))
					for i, id := range unit.QuestionIds {
						strIds[i] = strconv.FormatUint(uint64(id), 10)
					}
					questionIdsStr = strings.Join(strIds, ",")
				}

				// 处理标签
				tagsStr := strings.Join(unit.Tags, ",")

				// 创建单元
				courseUnit := &entity.CourseSectionUnit{
					SectionID:   courseSection.ID,
					Title:       unit.Title,
					Desc:        unit.Desc,
					QuestionIds: questionIdsStr,
					OrderIndex:  unit.OrderIndex,
					Status:      1, // 启用状态
					Tags:        tagsStr,
					Prompt:      unit.Prompt,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				if err := s.courseSectionUnitRepository.Create(ctx, courseUnit); err != nil {
					return nil, err
				}
			}
		}
	}
	return courseIds, nil
}
