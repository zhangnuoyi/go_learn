package services

import (
	"context"
	"errors"
	"github.com/user/rpc/models"
	"github.com/user/rpc/proto"
)

// UserService gRPC 用户服务实现
type UserService struct {
	proto.UnimplementedUserServiceServer
	userStore models.UserStore
}

// NewUserService 创建用户服务实例
func NewUserService(userStore models.UserStore) *UserService {
	return &UserService{
		userStore: userStore,
	}
}

// CreateUser 创建用户实现
func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// 创建用户模型
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Age:   req.Age,
	}

	// 保存用户
	err := s.userStore.Create(user)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	respUser := &proto.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &proto.CreateUserResponse{
		User:    respUser,
		Message: "用户创建成功",
	}, nil
}

// GetUser 获取用户实现
func (s *UserService) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	// 获取用户
	user, err := s.userStore.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("用户不存在")
	}

	// 转换为响应格式
	respUser := &proto.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return &proto.GetUserResponse{
		User: respUser,
	}, nil
}

// UpdateUser 更新用户实现
func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// 检查用户是否存在
	existingUser, err := s.userStore.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("用户不存在")
	}

	// 更新用户信息
	updatedUser := &models.User{
		ID:        req.ID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Age:       req.Age,
		CreatedAt: existingUser.CreatedAt,
		UpdatedAt: existingUser.UpdatedAt, // 将在存储层更新
	}

	// 保存更新
	err = s.userStore.Update(updatedUser)
	if err != nil {
		return nil, err
	}

	// 重新获取更新后的用户
	updatedUser, err = s.userStore.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	respUser := &proto.User{
		ID:        updatedUser.ID,
		Name:      updatedUser.Name,
		Email:     updatedUser.Email,
		Phone:     updatedUser.Phone,
		Age:       updatedUser.Age,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
	}

	return &proto.UpdateUserResponse{
		User:    respUser,
		Message: "用户更新成功",
	}, nil
}

// DeleteUser 删除用户实现
func (s *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// 检查用户是否存在
	existingUser, err := s.userStore.GetByID(req.ID)
	if err != nil {
		return nil, err
	}

	if existingUser == nil {
		return nil, errors.New("用户不存在")
	}

	// 删除用户
	err = s.userStore.Delete(req.ID)
	if err != nil {
		return nil, err
	}

	return &proto.DeleteUserResponse{
		Message: "用户删除成功",
	}, nil
}

// ListUsers 列出用户实现
func (s *UserService) ListUsers(ctx context.Context, req *proto.ListUsersRequest) (*proto.ListUsersResponse, error) {
	// 参数验证
	page := int(req.Page)
	if page <= 0 {
		page = 1
	}

	pageSize := int(req.PageSize)
	if pageSize <= 0 {
		pageSize = 10
	}

	// 获取用户列表
	users, total, err := s.userStore.List(page, pageSize)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	respUsers := make([]*proto.User, 0, len(users))
	for _, user := range users {
		respUsers = append(respUsers, &proto.User{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Phone:     user.Phone,
			Age:       user.Age,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	return &proto.ListUsersResponse{
		Users:    respUsers,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
