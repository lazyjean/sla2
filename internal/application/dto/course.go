package dto

import "github.com/lazyjean/sla2/internal/domain/entity"

// CourseListInput 课程列表查询参数
type CourseListInput struct {
	Page     int
	PageSize int
	Level    uint
	Category entity.CourseCategory
	Tags     []string
	Status   string
}

// CourseSearchInput 课程搜索参数
type CourseSearchInput struct {
	Keyword  string
	Page     int
	PageSize int
	Level    uint
	Category entity.CourseCategory
	Tags     []string
}
