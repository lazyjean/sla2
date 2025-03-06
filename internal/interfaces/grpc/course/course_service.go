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
	case pb.CourseLevel_COURSE_LEVEL_BEGINNER:
		return "beginner"
	case pb.CourseLevel_COURSE_LEVEL_INTERMEDIATE:
		return "intermediate"
	case pb.CourseLevel_COURSE_LEVEL_ADVANCED:
		return "advanced"
	default:
		return "beginner"
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
	case "beginner":
		return pb.CourseLevel_COURSE_LEVEL_BEGINNER
	case "intermediate":
		return pb.CourseLevel_COURSE_LEVEL_INTERMEDIATE
	case "advanced":
		return pb.CourseLevel_COURSE_LEVEL_ADVANCED
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
		req.Tags,
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
		req.Tags,
		convertStatusToString(req.Status),
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
		req.Tags,
		convertStatusToString(req.Status),
	)
	if err != nil {
		return nil, err
	}

	var pbCourses []*pb.SimpleCourse
	for _, course := range courses {
		pbCourses = append(pbCourses, &pb.SimpleCourse{
			Id:       uint32(course.ID),
			Title:    course.Title,
			CoverUrl: course.CoverURL,
			Level:    convertStringToLevel(course.Level),
			Tags:     course.Tags,
			Desc:     course.Description,
		})
	}

	return &pb.CourseServiceListResponse{
		Courses: pbCourses,
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
		Id:        uint32(course.ID),
		Title:     course.Title,
		Desc:      course.Description,
		CoverUrl:  course.CoverURL,
		Level:     convertStringToLevel(course.Level),
		Tags:      course.Tags,
		Status:    convertStringToStatus(course.Status),
		Version:   0, // 添加版本号字段
		Sections:  convertToPbCourseSections(course.Sections),
		CreatedAt: timestamppb.New(course.CreatedAt),
		UpdatedAt: timestamppb.New(course.UpdatedAt),
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
