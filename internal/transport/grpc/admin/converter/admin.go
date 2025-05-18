package converter

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

// AdminConverter 管理员转换器
type AdminConverter struct{}

// NewAdminConverter 创建新的管理员转换器
func NewAdminConverter() *AdminConverter {
	return &AdminConverter{}
}

// ToPB 将 dto.AdminInfo 转换为 proto.AdminInfo
func (c *AdminConverter) ToPB(admin *dto.AdminInfo) *pb.AdminInfo {
	if admin == nil {
		return nil
	}
	return &pb.AdminInfo{
		Id:            uint64(admin.ID),
		Username:      admin.Username,
		Nickname:      admin.Nickname,
		Email:         admin.Email,
		EmailVerified: admin.EmailVerified,
		Roles:         admin.Roles,
		CreatedAt:     admin.CreatedAt.Unix(),
		UpdatedAt:     admin.UpdatedAt.Unix(),
	}
}

// ToPBFromResponse 将 dto.AdminInfoResponse 转换为 proto.AdminInfo
func (c *AdminConverter) ToPBFromResponse(admin *dto.AdminInfoResponse) *pb.AdminInfo {
	if admin == nil {
		return nil
	}
	return &pb.AdminInfo{
		Id:            uint64(admin.ID),
		Username:      admin.Username,
		Nickname:      admin.Nickname,
		Email:         admin.Email,
		EmailVerified: admin.EmailVerified,
		Roles:         admin.Roles,
	}
}

// ToEntity 将 proto.AdminInfo 转换为 entity.Admin
func (c *AdminConverter) ToEntity(admin *pb.AdminInfo) *entity.Admin {
	if admin == nil {
		return nil
	}
	return &entity.Admin{
		ID:            entity.UID(admin.Id),
		Username:      admin.Username,
		Nickname:      admin.Nickname,
		Email:         admin.Email,
		EmailVerified: admin.EmailVerified,
		Roles:         admin.Roles,
	}
}

// ToPBList 将 dto.AdminInfo 切片转换为 proto.AdminInfo 切片
func (c *AdminConverter) ToPBList(admins []*dto.AdminInfo) []*pb.AdminInfo {
	if admins == nil {
		return nil
	}
	result := make([]*pb.AdminInfo, len(admins))
	for i, admin := range admins {
		result[i] = c.ToPB(admin)
	}
	return result
}
