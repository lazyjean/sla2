package course

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CourseService 课程服务
type CourseService struct {
	pb.UnimplementedCourseServiceServer
	courseService *service.CourseService
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseService *service.CourseService) *CourseService {
	return &CourseService{
		courseService: courseService,
	}
}

// convertLevelToString 将 protobuf 枚举类型转换为字符串
func convertLevelToString(level pb.CourseLevel) string {
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

// convertStatusToString 将 protobuf 枚举类型转换为字符串
func convertStatusToString(status pb.CourseStatus) string {
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

// convertStringToLevel 将字符串转换为 protobuf 枚举类型
func convertStringToLevel(level string) pb.CourseLevel {
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

// convertStringToStatus 将字符串转换为 protobuf 枚举类型
func convertStringToStatus(status string) pb.CourseStatus {
	switch status {
	case "draft":
		return pb.CourseStatus_COURSE_STATUS_DRAFT
	case "published":
		return pb.CourseStatus_COURSE_STATUS_PUBLISHED
	case "archived":
		return pb.CourseStatus_COURSE_STATUS_ARCHIVED
	default:
		return pb.CourseStatus_COURSE_STATUS_UNSPECIFIED
	}
}

// convertSectionStatusToString 将 protobuf 枚举类型转换为字符串
func convertSectionStatusToString(status pb.CourseSectionStatus) string {
	switch status {
	case pb.CourseSectionStatus_COURSE_SECTION_STATUS_ENABLED:
		return "enabled"
	case pb.CourseSectionStatus_COURSE_SECTION_STATUS_DISABLED:
		return "disabled"
	default:
		return "enabled"
	}
}

// convertStringToSectionStatus 将字符串转换为 protobuf 枚举类型
func convertStringToSectionStatus(status string) pb.CourseSectionStatus {
	switch status {
	case "enabled":
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_ENABLED
	case "disabled":
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_DISABLED
	default:
		return pb.CourseSectionStatus_COURSE_SECTION_STATUS_UNSPECIFIED
	}
}

// CreateCourse 创建课程
func (s *CourseService) Create(ctx context.Context, req *pb.CourseServiceCreateRequest) (*pb.CourseServiceCreateResponse, error) {
	course, err := s.courseService.CreateCourse(
		ctx,
		req.Title,
		req.Desc,
		req.CoverUrl,
		convertLevelToString(req.Level),
		convertCategoryToString(req.Category),
		req.Tags,
		req.Prompt,
		req.Resources,
		req.RecommendedAge,
		req.StudyPlan,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceCreateResponse{
		Id: uint32(course.ID),
	}, nil
}

// UpdateCourse 更新课程
func (s *CourseService) Update(ctx context.Context, req *pb.CourseServiceUpdateRequest) (*pb.CourseServiceUpdateResponse, error) {
	course, err := s.courseService.UpdateCourse(
		ctx,
		uint(req.Id),
		req.Title,
		req.Desc,
		req.CoverUrl,
		convertLevelToString(req.Level),
		convertCategoryToString(req.Category),
		req.Tags,
		convertStatusToString(req.Status),
		req.Prompt,
		req.Resources,
		req.RecommendedAge,
		req.StudyPlan,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceUpdateResponse{
		Id: uint32(course.ID),
	}, nil
}

// GetCourse 获取课程详情
func (s *CourseService) Get(ctx context.Context, req *pb.CourseServiceGetRequest) (*pb.CourseServiceGetResponse, error) {
	course, err := s.courseService.GetCourse(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceGetResponse{
		Course: convertToPbCourse(course),
	}, nil
}

// ListCourses 获取课程列表
func (s *CourseService) List(ctx context.Context, req *pb.CourseServiceListRequest) (*pb.CourseServiceListResponse, error) {
	courses, total, err := s.courseService.ListCourses(
		ctx,
		int(req.Page),
		int(req.PageSize),
		uint(req.Level),
		convertCategoryToString(req.Category),
		req.Tags,
		convertStatusToString(req.Status),
	)
	if err != nil {
		return nil, err
	}

	coursesPb := make([]*pb.SimpleCourse, 0, len(courses))
	for _, course := range courses {
		coursesPb = append(coursesPb, convertToSimpleCourseProto(course))
	}

	return &pb.CourseServiceListResponse{
		Courses: coursesPb,
		Total:   uint32(total),
	}, nil
}

// DeleteCourse 删除课程
func (s *CourseService) Delete(ctx context.Context, req *pb.CourseServiceDeleteRequest) (*pb.CourseServiceDeleteResponse, error) {
	id := req.Id

	err := s.courseService.DeleteCourse(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceDeleteResponse{}, nil
}

// convertToPbCourse 将实体转换为 protobuf 消息
func convertToPbCourse(course *entity.Course) *pb.Course {
	return &pb.Course{
		Id:             uint32(course.ID),
		Title:          course.Title,
		Desc:           course.Description,
		CoverUrl:       course.CoverURL,
		Level:          convertStringToLevel(course.Level),
		Category:       convertStringToCategory(course.Category),
		Tags:           course.Tags,
		Status:         convertStringToStatus(course.Status),
		Prompt:         course.Prompt,
		Resources:      course.Resources,
		RecommendedAge: course.RecommendedAge,
		StudyPlan:      course.StudyPlan,
		CreatedAt:      timestamppb.New(course.CreatedAt),
		UpdatedAt:      timestamppb.New(course.UpdatedAt),
	}
}

// convertToPbCourseSections 将课程章节实体转换为 protobuf 消息
func convertToPbCourseSections(sections []*entity.CourseSection) []*pb.CourseSection {
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
			Status:     convertStringToSectionStatus(section.Status),
			Units:      convertToPbCourseSectionUnits(section.Units),
			CreatedAt:  timestamppb.New(section.CreatedAt),
			UpdatedAt:  timestamppb.New(section.UpdatedAt),
		}
	}
	return pbSections
}

// convertToPbCourseSectionUnits 将课程章节单元实体转换为 protobuf 消息
func convertToPbCourseSectionUnits(units []*entity.CourseSectionUnit) []*pb.CourseSectionUnit {
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

// CreateSection 创建课程章节
func (s *CourseService) CreateSection(ctx context.Context, req *pb.CourseServiceCreateSectionRequest) (*pb.CourseServiceCreateSectionResponse, error) {
	section, err := s.courseService.CreateSection(ctx, uint(req.CourseId), req.Title, req.Desc)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceCreateSectionResponse{
		Id: int64(section.ID),
	}, nil
}

// UpdateSection 更新课程章节
func (s *CourseService) UpdateSection(ctx context.Context, req *pb.CourseServiceUpdateSectionRequest) (*pb.CourseServiceUpdateSectionResponse, error) {
	section, err := s.courseService.UpdateSection(ctx, entity.CourseSectionID(req.Id), req.Title, req.Desc, req.OrderIndex, convertSectionStatusToString(req.Status))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceUpdateSectionResponse{
		Id: int64(section.ID),
	}, nil
}

// DeleteSection 删除课程章节
func (s *CourseService) DeleteSection(ctx context.Context, req *pb.CourseServiceDeleteSectionRequest) (*pb.CourseServiceDeleteSectionResponse, error) {
	err := s.courseService.DeleteSection(ctx, entity.CourseSectionID(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceDeleteSectionResponse{}, nil
}

// DeleteUnit 删除章节单元
func (s *CourseService) DeleteUnit(ctx context.Context, req *pb.CourseServiceDeleteUnitRequest) (*pb.CourseServiceDeleteUnitResponse, error) {
	err := s.courseService.DeleteUnit(ctx, entity.CourseSectionUnitID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.CourseServiceDeleteUnitResponse{}, nil
}

// BatchCreate 批量创建课程
func (s *CourseService) BatchCreate(ctx context.Context, req *pb.CourseServiceBatchCreateRequest) (*pb.CourseServiceBatchCreateResponse, error) {
	courses := make([]struct {
		Title          string
		Description    string
		CoverURL       string
		Level          string
		Category       entity.CourseCategory
		Tags           []string
		Prompt         string
		Resources      []string
		RecommendedAge string
		StudyPlan      string
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
				Prompt      string
			}
		}
	}, 0, len(req.Courses))

	for _, coursePb := range req.Courses {
		sections := make([]struct {
			Title      string
			Desc       string
			OrderIndex int32
			Units      []struct {
				Title       string
				Desc        string
				QuestionIds []uint32
				OrderIndex  int32
				Tags        []string
				Prompt      string
			}
		}, 0, len(coursePb.Sections))

		for _, sectionPb := range coursePb.Sections {
			units := make([]struct {
				Title       string
				Desc        string
				QuestionIds []uint32
				OrderIndex  int32
				Tags        []string
				Prompt      string
			}, 0, len(sectionPb.Units))

			for _, unitPb := range sectionPb.Units {
				units = append(units, struct {
					Title       string
					Desc        string
					QuestionIds []uint32
					OrderIndex  int32
					Tags        []string
					Prompt      string
				}{
					Title:      unitPb.Title,
					Desc:       unitPb.Desc,
					OrderIndex: unitPb.OrderIndex,
					Tags:       unitPb.Labels,
					Prompt:     unitPb.Prompt,
				})
			}

			sections = append(sections, struct {
				Title      string
				Desc       string
				OrderIndex int32
				Units      []struct {
					Title       string
					Desc        string
					QuestionIds []uint32
					OrderIndex  int32
					Tags        []string
					Prompt      string
				}
			}{
				Title:      sectionPb.Title,
				Desc:       sectionPb.Desc,
				OrderIndex: sectionPb.OrderIndex,
				Units:      units,
			})
		}

		courses = append(courses, struct {
			Title          string
			Description    string
			CoverURL       string
			Level          string
			Category       entity.CourseCategory
			Tags           []string
			Prompt         string
			Resources      []string
			RecommendedAge string
			StudyPlan      string
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
					Prompt      string
				}
			}
		}{
			Title:          coursePb.Title,
			Description:    coursePb.Desc,
			CoverURL:       coursePb.CoverUrl,
			Level:          convertLevelToString(coursePb.Level),
			Category:       convertCategoryToString(coursePb.Category),
			Tags:           coursePb.Tags,
			Prompt:         coursePb.Prompt,
			Resources:      coursePb.Resources,
			RecommendedAge: coursePb.RecommendedAge,
			StudyPlan:      coursePb.StudyPlan,
			Sections:       sections,
		})
	}

	ids, err := s.courseService.BatchCreateCourse(ctx, courses)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceBatchCreateResponse{
		Ids: ids,
	}, nil
}

// convertToPbSimpleCourse 将实体转换为简化的 protobuf 消息
func convertToPbSimpleCourse(course *entity.Course) *pb.SimpleCourse {
	return &pb.SimpleCourse{
		Id:             uint32(course.ID),
		Title:          course.Title,
		Desc:           course.Description,
		CoverUrl:       course.CoverURL,
		Level:          convertStringToLevel(course.Level),
		Tags:           course.Tags,
		Resources:      course.Resources,
		RecommendedAge: course.RecommendedAge,
		StudyPlan:      course.StudyPlan,
	}
}

// convertCategoryToString 将 protobuf 枚举类型转换为领域实体的 CourseCategory
func convertCategoryToString(category pb.CourseCategory) entity.CourseCategory {
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

// convertStringToCategory 将领域实体的 CourseCategory 转换为 protobuf 枚举类型
func convertStringToCategory(category entity.CourseCategory) pb.CourseCategory {
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

// Search 搜索课程
func (s *CourseService) Search(ctx context.Context, req *pb.CourseServiceSearchRequest) (*pb.CourseServiceSearchResponse, error) {
	courses, total, err := s.courseService.SearchCourse(
		ctx,
		req.Keyword,
		int(req.Page),
		int(req.PageSize),
		uint(req.Level),
		convertCategoryToString(req.Category),
		req.Tags,
	)
	if err != nil {
		return nil, err
	}

	coursesPb := make([]*pb.SimpleCourse, 0, len(courses))
	for _, course := range courses {
		coursesPb = append(coursesPb, convertToSimpleCourseProto(course))
	}

	return &pb.CourseServiceSearchResponse{
		Courses: coursesPb,
		Total:   uint32(total),
	}, nil
}

// convertToSimpleCourseProto 将实体转换为简化的 protobuf 消息
func convertToSimpleCourseProto(course *entity.Course) *pb.SimpleCourse {
	return &pb.SimpleCourse{
		Id:             uint32(course.ID),
		Title:          course.Title,
		Desc:           course.Description,
		CoverUrl:       course.CoverURL,
		Level:          convertStringToLevel(course.Level),
		Category:       convertStringToCategory(course.Category),
		Tags:           course.Tags,
		Resources:      course.Resources,
		RecommendedAge: course.RecommendedAge,
		StudyPlan:      course.StudyPlan,
	}
}
