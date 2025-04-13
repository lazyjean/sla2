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
func (s *CourseService) CreateCourse(ctx context.Context, title, description, coverURL, level string, category entity.CourseCategory, tags []string, prompt string, resources []string, recommendedAge string, studyPlan string) (*entity.Course, error) {
	course := &entity.Course{
		Title:          title,
		Description:    description,
		CoverURL:       coverURL,
		Level:          level,
		Category:       category,
		Tags:           tags,
		Status:         "draft",
		Prompt:         prompt,
		Resources:      resources,
		RecommendedAge: recommendedAge,
		StudyPlan:      studyPlan,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.courseRepository.Create(ctx, course); err != nil {
		return nil, err
	}

	return course, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, id uint, title, description, coverURL, level string, category entity.CourseCategory, tags []string, status string, prompt string, resources []string, recommendedAge string, studyPlan string) (*entity.Course, error) {
	course, err := s.courseRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	course.Title = title
	course.Description = description
	course.CoverURL = coverURL
	course.Level = level
	course.Category = category
	course.Tags = tags
	course.Status = status
	course.Prompt = prompt
	course.Resources = resources
	course.RecommendedAge = recommendedAge
	course.StudyPlan = studyPlan
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
func (s *CourseService) ListCourses(ctx context.Context, page, pageSize int, level uint, category entity.CourseCategory, tags []string, status string) ([]*entity.Course, int64, error) {
	offset := (page - 1) * pageSize

	filters := make(map[string]interface{})
	if level != 0 {
		filters["level"] = level
	}
	if category != "" {
		filters["category"] = category
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
func (s *CourseService) SearchCourse(ctx context.Context, keyword string, page, pageSize int, level uint, category entity.CourseCategory, tags []string) ([]*entity.Course, int64, error) {
	offset := (page - 1) * pageSize
	filters := make(map[string]interface{})
	if level != 0 {
		filters["level"] = level
	}
	if category != "" {
		filters["category"] = category
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
func (s *CourseService) CreateUnit(ctx context.Context, sectionID entity.CourseSectionID, title, desc string, questionIds []uint32, tags []string, prompt string) (*entity.CourseSectionUnit, error) {
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
		Prompt:      prompt,
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
func (s *CourseService) UpdateUnit(ctx context.Context, id entity.CourseSectionUnitID, title, desc string, questionIds []uint32, tags []string, status int32, prompt string) (*entity.CourseSectionUnit, error) {
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
	unit.Prompt = prompt
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

// BatchCreateCourse 批量创建课程
func (s *CourseService) BatchCreateCourse(ctx context.Context, courses []struct {
	Title          string
	Description    string
	CoverURL       string
	Level          string
	Category       entity.CourseCategory
	Tags           []string
	Prompt         string   // AI 提示词
	Resources      []string // 推荐学习资源 URL 列表
	RecommendedAge string   // 推荐年龄范围
	StudyPlan      string   // 建议学习计划
	Sections       []struct {
		Title      string
		Desc       string
		OrderIndex int32
		Units      []struct {
			Title       string
			Desc        string
			QuestionIds []uint32
			OrderIndex  int32
			Tags        []string
			Prompt      string // AI 提示词
		}
	}
}) ([]uint32, error) {
	// 用于存储创建的课程ID列表
	var courseIds []uint32

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
		courseIds = append(courseIds, uint32(course.ID))

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

				if err := s.courseSectionRepository.CreateUnit(ctx, courseUnit); err != nil {
					return nil, err
				}
			}
		}
	}

	return courseIds, nil
}
