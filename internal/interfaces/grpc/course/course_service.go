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

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(ctx context.Context, req *pb.CreateCourseRequest) (*pb.CreateCourseResponse, error) {
	course, err := s.courseService.CreateCourse(
		ctx,
		req.Title,
		req.Description,
		req.CoverUrl,
		req.Level,
		int(req.Duration),
		req.Tags,
	)
	if err != nil {
		return nil, err
	}

	return &pb.CreateCourseResponse{
		Course: convertToPbCourse(course),
	}, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(ctx context.Context, req *pb.UpdateCourseRequest) (*pb.UpdateCourseResponse, error) {
	course, err := s.courseService.UpdateCourse(
		ctx,
		uint(req.Id),
		req.Title,
		req.Description,
		req.CoverUrl,
		req.Level,
		int(req.Duration),
		req.Tags,
		req.Status,
	)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateCourseResponse{
		Course: convertToPbCourse(course),
	}, nil
}

// GetCourse 获取课程详情
func (s *CourseService) GetCourse(ctx context.Context, req *pb.GetCourseRequest) (*pb.GetCourseResponse, error) {
	course, err := s.courseService.GetCourse(ctx, uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &pb.GetCourseResponse{
		Course: convertToPbCourse(course),
	}, nil
}

// ListCourses 获取课程列表
func (s *CourseService) ListCourses(ctx context.Context, req *pb.ListCoursesRequest) (*pb.ListCoursesResponse, error) {
	courses, total, err := s.courseService.ListCourses(
		ctx,
		int(req.Page),
		int(req.PageSize),
		uint(req.Level),
		[]string{req.Tag},
		req.Status,
	)
	if err != nil {
		return nil, err
	}

	var pbCourses []*pb.Course
	for _, course := range courses {
		pbCourses = append(pbCourses, convertToPbCourse(course))
	}

	return &pb.ListCoursesResponse{
		Courses: pbCourses,
		Total:   uint32(total),
	}, nil
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(ctx context.Context, req *pb.DeleteCourseRequest) (*pb.DeleteCourseResponse, error) {
	id := req.Id

	err := s.courseService.DeleteCourse(ctx, uint(id))
	if err != nil {
		return nil, err
	}

	return &pb.DeleteCourseResponse{}, nil
}

// convertToPbCourse 将实体转换为 protobuf 消息
func convertToPbCourse(course *entity.Course) *pb.Course {
	return &pb.Course{
		Id:          uint32(course.ID),
		Title:       course.Title,
		Description: course.Description,
		CoverUrl:    course.CoverURL,
		Level:       course.Level,
		Duration:    uint32(course.Duration),
		Tags:        course.Tags,
		Status:      course.Status,
		CreatedAt:   timestamppb.New(course.CreatedAt),
		UpdatedAt:   timestamppb.New(course.UpdatedAt),
	}
}
