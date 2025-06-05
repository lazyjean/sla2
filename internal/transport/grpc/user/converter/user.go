package converter

import (
	pb "github.com/lazyjean/sla2/api/proto/v1"
	"github.com/lazyjean/sla2/internal/application/dto"
	"github.com/lazyjean/sla2/internal/domain/entity"
)

type UserConverter struct{}

func NewUserConverter() *UserConverter {
	return &UserConverter{}
}

// ToProtoUser 将实体转换为 Proto User 消息
func (c *UserConverter) ToProtoUser(user *entity.User) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:                uint64(user.ID),
		Username:          user.Username,
		Email:             user.Email,
		Nickname:          user.Nickname,
		Avatar:            user.Avatar,
		Status:            c.ConvertStatusToProto(user.Status),
		CreatedAt:         user.CreatedAt.Unix(),
		UpdatedAt:         user.UpdatedAt.Unix(),
		AddedCourseIds:    nil, // TODO: 实现课程相关功能
		CurrentLearningId: 0,   // TODO: 实现课程相关功能
	}
}

// ToProtoUserDTO 将 DTO 转换为 Proto User 消息
func (c *UserConverter) ToProtoUserDTO(user *dto.UserDTO) *pb.User {
	if user == nil {
		return nil
	}

	return &pb.User{
		Id:                uint64(user.ID),
		Username:          user.Username,
		Email:             user.Email,
		Nickname:          user.Nickname,
		Avatar:            user.Avatar,
		Status:            c.ConvertStatusStringToProto(user.Status),
		CreatedAt:         user.CreatedAt.Unix(),
		UpdatedAt:         user.UpdatedAt.Unix(),
		AddedCourseIds:    nil, // TODO: 实现课程相关功能
		CurrentLearningId: 0,   // TODO: 实现课程相关功能
	}
}

// PbToEntityUser 将 Proto User 消息转换为实体
func (c *UserConverter) PbToEntityUser(userPb *pb.User) *entity.User {
	if userPb == nil {
		return nil
	}

	return &entity.User{
		Username: userPb.Username,
		Email:    userPb.Email,
		Nickname: userPb.Nickname,
		Avatar:   userPb.Avatar,
		Status:   c.ConvertStatusToEntity(userPb.Status),
	}
}

// ConvertStatusToProto 将实体状态转换为 Proto 状态
func (c *UserConverter) ConvertStatusToProto(status entity.UserStatus) pb.UserStatus {
	switch status {
	case entity.UserStatusActive:
		return pb.UserStatus_USER_STATUS_ACTIVE
	case entity.UserStatusInactive:
		return pb.UserStatus_USER_STATUS_INACTIVE
	case entity.UserStatusSuspended:
		return pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// ConvertStatusToEntity 将 Proto 状态转换为实体状态
func (c *UserConverter) ConvertStatusToEntity(status pb.UserStatus) entity.UserStatus {
	switch status {
	case pb.UserStatus_USER_STATUS_ACTIVE:
		return entity.UserStatusActive
	case pb.UserStatus_USER_STATUS_INACTIVE:
		return entity.UserStatusInactive
	case pb.UserStatus_USER_STATUS_SUSPENDED:
		return entity.UserStatusSuspended
	default:
		return entity.UserStatusActive
	}
}

// ConvertStatusStringToProto 将状态字符串转换为 Proto 状态
func (c *UserConverter) ConvertStatusStringToProto(status string) pb.UserStatus {
	switch status {
	case "active":
		return pb.UserStatus_USER_STATUS_ACTIVE
	case "inactive":
		return pb.UserStatus_USER_STATUS_INACTIVE
	case "suspended":
		return pb.UserStatus_USER_STATUS_SUSPENDED
	default:
		return pb.UserStatus_USER_STATUS_UNSPECIFIED
	}
}

// ToProtoLoginResponse 将登录响应 DTO 转换为 Proto 消息
func (c *UserConverter) ToProtoLoginResponse(resp *dto.LoginResponse) *pb.LoginResponse {
	if resp == nil {
		return nil
	}

	return &pb.LoginResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
	}
}

// ToProtoRegisterResponse 将注册响应 DTO 转换为 Proto 消息
func (c *UserConverter) ToProtoRegisterResponse(resp *dto.RegisterResponse) *pb.RegisterResponse {
	if resp == nil {
		return nil
	}

	return &pb.RegisterResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: "", // TODO: 从 resp 中获取
			Email:    "", // TODO: 从 resp 中获取
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
	}
}

// ToProtoAppleLoginResponse 将苹果登录响应 DTO 转换为 Proto 消息
func (c *UserConverter) ToProtoAppleLoginResponse(resp *dto.AppleLoginResponse) *pb.AppleLoginResponse {
	if resp == nil {
		return nil
	}

	return &pb.AppleLoginResponse{
		User: &pb.User{
			Id:       uint64(resp.UserID),
			Username: resp.Username,
			Email:    resp.Email,
			Nickname: resp.Nickname,
			Avatar:   resp.Avatar,
			Status:   pb.UserStatus_USER_STATUS_ACTIVE,
		},
		Token:        resp.Token,
		RefreshToken: resp.RefreshToken,
		IsNewUser:    resp.IsNewUser,
	}
}
