package vocabulary

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/service"
)

// VocabularyService 词汇服务实现
type VocabularyService struct {
	pb.UnimplementedVocabularyServiceServer
	service *service.VocabularyService
}

// NewVocabularyService 创建词汇服务实例
func NewVocabularyService(service *service.VocabularyService) *VocabularyService {
	return &VocabularyService{
		service: service,
	}
}
