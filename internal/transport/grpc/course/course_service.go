package course

import (
	"context"

	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
	"github.com/lazyjean/sla2/internal/domain/entity"
	"github.com/lazyjean/sla2/internal/transport/grpc/course/converter"
)

// Service 课程服务
type Service struct {
	pb.UnimplementedCourseServiceServer
	courseService service.CourseService
	converter     *converter.CourseConverter
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseService service.CourseService) *Service {
	return &Service{
		courseService: courseService,
		converter:     converter.NewCourseConverter(),
	}
}

// Create CreateCourse 创建课程
func (s *Service) Create(ctx context.Context, req *pb.CourseServiceCreateRequest) (*pb.CourseServiceCreateResponse, error) {
	course := s.converter.ToEntityCourse(req)

	result, err := s.courseService.CreateCourse(ctx, course)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceCreateResponse{
		Id: uint32(result.ID),
	}, nil
}

// Update 更新课程
func (s *Service) Update(ctx context.Context, req *pb.CourseServiceUpdateRequest) (*pb.CourseServiceUpdateResponse, error) {
	course := s.converter.ToEntityCourseForUpdate(req)
	course.ID = s.converter.ToEntityCourseID(req.Id)

	if err := s.courseService.UpdateCourse(ctx, course.ID, course); err != nil {
		return nil, err
	}

	return &pb.CourseServiceUpdateResponse{
		Id: req.Id,
	}, nil
}

// Get 获取课程详情
func (s *Service) Get(ctx context.Context, req *pb.CourseServiceGetRequest) (*pb.CourseServiceGetResponse, error) {
	course, err := s.courseService.GetCourse(ctx, s.converter.ToEntityCourseID(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceGetResponse{
		Course: s.converter.ToPB(course),
	}, nil
}

// List 获取课程列表
func (s *Service) List(ctx context.Context, req *pb.CourseServiceListRequest) (*pb.CourseServiceListResponse, error) {
	input := s.converter.ToEntityCourseListInput(req)
	courses, total, err := s.courseService.ListCourses(ctx, input)
	if err != nil {
		return nil, err
	}

	coursesPb := make([]*pb.SimpleCourse, 0, len(courses))
	for _, course := range courses {
		coursesPb = append(coursesPb, s.converter.ToSimplePB(course))
	}

	return &pb.CourseServiceListResponse{
		Courses: coursesPb,
		Total:   uint32(total),
	}, nil
}

// Delete 删除课程
func (s *Service) Delete(ctx context.Context, req *pb.CourseServiceDeleteRequest) (*pb.CourseServiceDeleteResponse, error) {
	err := s.courseService.DeleteCourse(ctx, s.converter.ToEntityCourseID(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceDeleteResponse{}, nil
}

// CreateSection 创建课程章节
func (s *Service) CreateSection(ctx context.Context, req *pb.CourseServiceCreateSectionRequest) (*pb.CourseServiceCreateSectionResponse, error) {
	id, err := s.courseService.CreateSection(ctx, s.converter.ToEntityCourseID(req.CourseId), int32(req.OrderIndex), req.Title, req.Desc)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceCreateSectionResponse{
		Id: int64(id),
	}, nil
}

// UpdateSection 更新课程章节
func (s *Service) UpdateSection(ctx context.Context, req *pb.CourseServiceUpdateSectionRequest) (*pb.CourseServiceUpdateSectionResponse, error) {
	err := s.courseService.UpdateSection(ctx, entity.CourseSectionID(req.Id), req.Title, req.Desc, req.OrderIndex, s.converter.ToEntitySectionStatus(req.Status))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceUpdateSectionResponse{}, nil
}

// DeleteSection 删除课程章节
func (s *Service) DeleteSection(ctx context.Context, req *pb.CourseServiceDeleteSectionRequest) (*pb.CourseServiceDeleteSectionResponse, error) {
	err := s.courseService.DeleteSection(ctx, entity.CourseSectionID(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceDeleteSectionResponse{}, nil
}

// DeleteUnit 删除章节单元
func (s *Service) DeleteUnit(ctx context.Context, req *pb.CourseServiceDeleteUnitRequest) (*pb.CourseServiceDeleteUnitResponse, error) {
	err := s.courseService.DeleteUnit(ctx, entity.CourseSectionUnitID(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.CourseServiceDeleteUnitResponse{}, nil
}

// BatchCreate 批量创建课程
func (s *Service) BatchCreate(ctx context.Context, req *pb.CourseServiceBatchCreateRequest) (*pb.CourseServiceBatchCreateResponse, error) {
	courses := s.converter.ToEntityBatchCreateCourseInput(req)
	ids, err := s.courseService.BatchCreateCourse(ctx, courses)
	if err != nil {
		return nil, err
	}

	return &pb.CourseServiceBatchCreateResponse{
		Ids: s.converter.ToPBCourseIDs(ids),
	}, nil
}

// Search 搜索课程
func (s *Service) Search(ctx context.Context, req *pb.CourseServiceSearchRequest) (*pb.CourseServiceSearchResponse, error) {
	input := s.converter.ToEntityCourseSearchInput(req)
	courses, total, err := s.courseService.SearchCourse(ctx, input)
	if err != nil {
		return nil, err
	}

	coursesPb := make([]*pb.SimpleCourse, 0, len(courses))
	for _, course := range courses {
		coursesPb = append(coursesPb, s.converter.ToSimplePB(course))
	}

	return &pb.CourseServiceSearchResponse{
		Courses: coursesPb,
		Total:   uint32(total),
	}, nil
}
