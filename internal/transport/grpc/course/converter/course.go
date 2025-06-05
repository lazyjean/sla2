package converter

import (
	"time"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CourseConverter 处理课程相关的转换
type CourseConverter struct{}

// NewCourseConverter 创建新的 CourseConverter 实例
func NewCourseConverter() *CourseConverter {
	return &CourseConverter{}
}

// ToPB 将课程实体转换为 PB 消息
func (c *CourseConverter) ToPB(course *entity.Course) *pb.Course {
	return &pb.Course{
		Id:             uint32(course.ID),
		Title:          course.Title,
		Desc:           course.Description,
		CoverUrl:       course.CoverURL,
		Level:          c.ToPBLevel(course.Level),
		Category:       c.ToPBCategory(course.Category),
		Tags:           course.Tags,
		Status:         c.ToPBStatus(course.Status),
		Prompt:         course.Prompt,
		Resources:      course.Resources,
		RecommendedAge: course.RecommendedAge,
		StudyPlan:      course.StudyPlan,
		CreatedAt:      timestamppb.New(course.CreatedAt),
		UpdatedAt:      timestamppb.New(course.UpdatedAt),
	}
}

// ToSimplePB 将课程实体转换为简化的 PB 消息
func (c *CourseConverter) ToSimplePB(course *entity.Course) *pb.SimpleCourse {
	return &pb.SimpleCourse{
		Id:             uint32(course.ID),
		Title:          course.Title,
		Desc:           course.Description,
		CoverUrl:       course.CoverURL,
		Level:          c.ToPBLevel(course.Level),
		Category:       c.ToPBCategory(course.Category),
		Tags:           course.Tags,
		Resources:      course.Resources,
		RecommendedAge: course.RecommendedAge,
		StudyPlan:      course.StudyPlan,
	}
}

// ToPBSections 将课程章节实体切片转换为 PB 消息切片
func (c *CourseConverter) ToPBSections(sections []*entity.CourseSection) []*pb.CourseSection {
	if sections == nil {
		return nil
	}

	pbSections := make([]*pb.CourseSection, len(sections))
	for i, section := range sections {
		pbSections[i] = &pb.CourseSection{
			Id:         int64(section.ID),
			Title:      section.Title,
			Desc:       section.Desc,
			OrderIndex: section.OrderIndex,
			Status:     c.ToPBSectionStatus(section.Status),
			Units:      c.ToPBUnits(section.Units),
			CreatedAt:  timestamppb.New(section.CreatedAt),
			UpdatedAt:  timestamppb.New(section.UpdatedAt),
		}
	}
	return pbSections
}

// ToPBUnits 将课程章节单元实体切片转换为 PB 消息切片
func (c *CourseConverter) ToPBUnits(units []*entity.CourseSectionUnit) []*pb.CourseSectionUnit {
	if units == nil {
		return nil
	}

	pbUnits := make([]*pb.CourseSectionUnit, len(units))
	for i, unit := range units {
		pbUnits[i] = &pb.CourseSectionUnit{
			Id:          int64(unit.ID),
			Title:       unit.Title,
			Desc:        unit.Desc,
			QuestionIds: unit.QuestionIds,
			OrderIndex:  unit.OrderIndex,
			Status:      int32(unit.Status),
			Tags:        unit.Tags,
			CreatedAt:   timestamppb.New(unit.CreatedAt),
			UpdatedAt:   timestamppb.New(unit.UpdatedAt),
		}
	}
	return pbUnits
}

// ToEntityLevel 将 PB 级别转换为实体级别
func (c *CourseConverter) ToEntityLevel(level pb.CourseLevel) string {
	switch level {
	case pb.CourseLevel_COURSE_LEVEL_A1:
		return "a1"
	case pb.CourseLevel_COURSE_LEVEL_A2:
		return "a2"
	case pb.CourseLevel_COURSE_LEVEL_B1:
		return "b1"
	case pb.CourseLevel_COURSE_LEVEL_B2:
		return "b2"
	case pb.CourseLevel_COURSE_LEVEL_C1:
		return "c1"
	case pb.CourseLevel_COURSE_LEVEL_C2:
		return "c2"
	default:
		return "a1"
	}
}

// ToPBLevel 将实体级别转换为 PB 级别
func (c *CourseConverter) ToPBLevel(level string) pb.CourseLevel {
	switch level {
	case "a1":
		return pb.CourseLevel_COURSE_LEVEL_A1
	case "a2":
		return pb.CourseLevel_COURSE_LEVEL_A2
	case "b1":
		return pb.CourseLevel_COURSE_LEVEL_B1
	case "b2":
		return pb.CourseLevel_COURSE_LEVEL_B2
	case "c1":
		return pb.CourseLevel_COURSE_LEVEL_C1
	case "c2":
		return pb.CourseLevel_COURSE_LEVEL_C2
	default:
		return pb.CourseLevel_COURSE_LEVEL_UNSPECIFIED
	}
}

// ToEntityStatus 将 PB 状态转换为实体状态
func (c *CourseConverter) ToEntityStatus(status pb.CourseStatus) string {
	switch status {
	case pb.CourseStatus_COURSE_STATUS_DRAFT:
		return "draft"
	case pb.CourseStatus_COURSE_STATUS_PUBLISHED:
		return "published"
	case pb.CourseStatus_COURSE_STATUS_ARCHIVED:
		return "archived"
	default:
		return "draft"
	}
}

// ToPBStatus 将实体状态转换为 PB 状态
func (c *CourseConverter) ToPBStatus(status string) pb.CourseStatus {
	switch status {
	case "draft":
		return pb.CourseStatus_COURSE_STATUS_DRAFT
	case "published":
		return pb.CourseStatus_COURSE_STATUS_PUBLISHED
	case "archived":
		return pb.CourseStatus_COURSE_STATUS_ARCHIVED
	default:
		return pb.CourseStatus_COURSE_STATUS_DRAFT
	}
}

// ToEntityCategory 将 PB 分类转换为实体分类
func (c *CourseConverter) ToEntityCategory(category pb.CourseCategory) entity.CourseCategory {
	switch category {
	case pb.CourseCategory_COURSE_CATEGORY_ENGLISH:
		return entity.CourseCategoryEnglish
	case pb.CourseCategory_COURSE_CATEGORY_CHINESE:
		return entity.CourseCategoryChinese
	case pb.CourseCategory_COURSE_CATEGORY_OTHER:
		return entity.CourseCategoryOther
	default:
		return entity.CourseCategoryUnspecified
	}
}

// ToPBCategory 将实体分类转换为 PB 分类
func (c *CourseConverter) ToPBCategory(category entity.CourseCategory) pb.CourseCategory {
	switch category {
	case entity.CourseCategoryEnglish:
		return pb.CourseCategory_COURSE_CATEGORY_ENGLISH
	case entity.CourseCategoryChinese:
		return pb.CourseCategory_COURSE_CATEGORY_CHINESE
	case entity.CourseCategoryOther:
		return pb.CourseCategory_COURSE_CATEGORY_OTHER
	default:
		return pb.CourseCategory_COURSE_CATEGORY_UNSPECIFIED
	}
}

// ToPBSectionStatus 将实体章节状态转换为 PB 状态
func (c *CourseConverter) ToPBSectionStatus(status string) pb.CourseSectionStatus {
	switch status {
	case "enabled":
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_ENABLED
	case "disabled":
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_DISABLED
	default:
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_UNSPECIFIED
	}
}

// ToEntitySectionStatus 将 PB 章节状态转换为实体状态
func (c *CourseConverter) ToEntitySectionStatus(status pb.CourseSectionStatus) string {
	switch status {
	case pb.CourseSectionStatus_COURSE_SECTION_STATUS_ENABLED:
		return "enabled"
	case pb.CourseSectionStatus_COURSE_SECTION_STATUS_DISABLED:
		return "disabled"
	default:
		return "enabled" // 默认启用
	}
}

// ToEntityCourse 将 CourseServiceCreateRequest 转换为领域实体 Course
func (c *CourseConverter) ToEntityCourse(req *pb.CourseServiceCreateRequest) *entity.Course {
	return &entity.Course{
		Title:          req.Title,
		Description:    req.Desc,
		CoverURL:       req.CoverUrl,
		Level:          c.ToEntityLevel(req.Level),
		Category:       c.ToEntityCategory(req.Category),
		Tags:           req.Tags,
		Status:         "draft",
		Prompt:         req.Prompt,
		Resources:      req.Resources,
		RecommendedAge: req.RecommendedAge,
		StudyPlan:      req.StudyPlan,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// ToEntityCourseID 将 uint32 转换为 CourseID
func (c *CourseConverter) ToEntityCourseID(id uint32) entity.CourseID {
	return entity.CourseID(id)
}

// ToEntityCourseListInput 将 CourseServiceListRequest 转换为 CourseListInput
func (c *CourseConverter) ToEntityCourseListInput(req *pb.CourseServiceListRequest) *dto.CourseListInput {
	return &dto.CourseListInput{
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Level:    uint(req.Level),
		Category: c.ToEntityCategory(req.Category),
		Tags:     req.Tags,
		Status:   c.ToEntityStatus(req.Status),
	}
}

// ToEntityCourseSearchInput 将 CourseServiceSearchRequest 转换为 CourseSearchInput
func (c *CourseConverter) ToEntityCourseSearchInput(req *pb.CourseServiceSearchRequest) *dto.CourseSearchInput {
	return &dto.CourseSearchInput{
		Keyword:  req.Keyword,
		Page:     int(req.Page),
		PageSize: int(req.PageSize),
		Level:    uint(req.Level),
		Category: c.ToEntityCategory(req.Category),
		Tags:     req.Tags,
	}
}

// ToEntityBatchCreateCourseInput 将 CourseServiceBatchCreateRequest 转换为 BatchCreateCourseInput
func (c *CourseConverter) ToEntityBatchCreateCourseInput(req *pb.CourseServiceBatchCreateRequest) []*entity.BatchCreateCourseInput {
	courses := make([]*entity.BatchCreateCourseInput, 0, len(req.Courses))
	for _, coursePb := range req.Courses {
		sections := make([]entity.BatchCreateSectionInput, 0, len(coursePb.Sections))
		for _, sectionPb := range coursePb.Sections {
			units := make([]entity.BatchCreateUnitInput, 0, len(sectionPb.Units))
			for _, unitPb := range sectionPb.Units {
				units = append(units, entity.BatchCreateUnitInput{
					Title:      unitPb.Title,
					Desc:       unitPb.Desc,
					OrderIndex: unitPb.OrderIndex,
					Tags:       unitPb.Labels,
					Prompt:     unitPb.Prompt,
				})
			}
			sections = append(sections, entity.BatchCreateSectionInput{
				Title:      sectionPb.Title,
				Desc:       sectionPb.Desc,
				OrderIndex: sectionPb.OrderIndex,
				Units:      units,
			})
		}
		courses = append(courses, &entity.BatchCreateCourseInput{
			Title:          coursePb.Title,
			Description:    coursePb.Desc,
			CoverURL:       coursePb.CoverUrl,
			Level:          c.ToEntityLevel(coursePb.Level),
			Category:       c.ToEntityCategory(coursePb.Category),
			Tags:           coursePb.Tags,
			Prompt:         coursePb.Prompt,
			Resources:      coursePb.Resources,
			RecommendedAge: coursePb.RecommendedAge,
			StudyPlan:      coursePb.StudyPlan,
			Sections:       sections,
		})
	}
	return courses
}

// ToPBCourseIDs 将 CourseID 切片转换为 uint32 切片
func (c *CourseConverter) ToPBCourseIDs(ids []entity.CourseID) []uint32 {
	result := make([]uint32, len(ids))
	for i, id := range ids {
		result[i] = uint32(id)
	}
	return result
}

// ToEntityCourseForUpdate 将 CourseServiceUpdateRequest 转换为领域实体 Course
func (c *CourseConverter) ToEntityCourseForUpdate(req *pb.CourseServiceUpdateRequest) *entity.Course {
	return &entity.Course{
		Title:          req.Title,
		Description:    req.Desc,
		CoverURL:       req.CoverUrl,
		Level:          c.ToEntityLevel(req.Level),
		Category:       c.ToEntityCategory(req.Category),
		Tags:           req.Tags,
		Status:         c.ToEntityStatus(req.Status),
		Prompt:         req.Prompt,
		Resources:      req.Resources,
		RecommendedAge: req.RecommendedAge,
		StudyPlan:      req.StudyPlan,
		UpdatedAt:      time.Now(),
	}
}
