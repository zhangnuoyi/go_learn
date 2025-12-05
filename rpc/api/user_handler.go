package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/user/rpc/proto"
	"github.com/user/rpc/services"
)

// UserHandler 用户 API 处理器
type UserHandler struct {
	userService *services.UserService
}

// NewUserHandler 创建用户 API 处理器实例
func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "用户信息"
// @Success 201 {object} CreateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 调用 gRPC 服务
	grpcReq := &proto.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Age:   req.Age,
	}

	resp, err := h.userService.CreateUser(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 转换为 API 响应格式
	c.JSON(http.StatusCreated, CreateUserResponse{
		User:    convertToUserDTO(resp.User),
		Message: resp.Message,
	})
}

// GetUser 获取用户
// @Summary 获取用户
// @Description 根据 ID 获取用户信息
// @Tags users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的用户 ID"})
		return
	}

	// 调用 gRPC 服务
	grpcReq := &proto.GetUserRequest{
		ID: uint32(id),
	}

	resp, err := h.userService.GetUser(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// 转换为 API 响应格式
	c.JSON(http.StatusOK, UserResponse{
		User: convertToUserDTO(resp.User),
	})
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户信息
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户 ID"
// @Param user body UpdateUserRequest true "用户信息"
// @Success 200 {object} UpdateUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的用户 ID"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// 调用 gRPC 服务
	grpcReq := &proto.UpdateUserRequest{
		ID:    uint32(id),
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
		Age:   req.Age,
	}

	resp, err := h.userService.UpdateUser(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// 转换为 API 响应格式
	c.JSON(http.StatusOK, UpdateUserResponse{
		User:    convertToUserDTO(resp.User),
		Message: resp.Message,
	})
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags users
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} DeleteUserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "无效的用户 ID"})
		return
	}

	// 调用 gRPC 服务
	grpcReq := &proto.DeleteUserRequest{
		ID: uint32(id),
	}

	resp, err := h.userService.DeleteUser(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: err.Error()})
		return
	}

	// 转换为 API 响应格式
	c.JSON(http.StatusOK, DeleteUserResponse{
		Message: resp.Message,
	})
}

// ListUsers 列出用户
// @Summary 列出用户
// @Description 分页列出用户
// @Tags users
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} ListUsersResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 调用 gRPC 服务
	grpcReq := &proto.ListUsersRequest{
		Page:     int32(page),
		PageSize: int32(pageSize),
	}

	resp, err := h.userService.ListUsers(c.Request.Context(), grpcReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	// 转换为 API 响应格式
	users := make([]UserDTO, 0, len(resp.Users))
	for _, user := range resp.Users {
		users = append(users, convertToUserDTO(user))
	}

	c.JSON(http.StatusOK, ListUsersResponse{
		Users:    users,
		Total:    resp.Total,
		Page:     resp.Page,
		PageSize: resp.PageSize,
	})
}

// convertToUserDTO 转换用户模型为 DTO
func convertToUserDTO(user *proto.User) UserDTO {
	return UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Phone:     user.Phone,
		Age:       user.Age,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
