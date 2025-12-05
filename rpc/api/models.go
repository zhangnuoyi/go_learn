package api

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
	Age   int32  `json:"age" binding:"gte=0,lte=200"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Phone string `json:"phone"`
	Age   int32  `json:"age" binding:"gte=0,lte=200"`
}

// UserDTO 用户数据传输对象
type UserDTO struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Age       int32  `json:"age"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// CreateUserResponse 创建用户响应
type CreateUserResponse struct {
	User    UserDTO `json:"user"`
	Message string  `json:"message"`
}

// UserResponse 用户响应
type UserResponse struct {
	User UserDTO `json:"user"`
}

// UpdateUserResponse 更新用户响应
type UpdateUserResponse struct {
	User    UserDTO `json:"user"`
	Message string  `json:"message"`
}

// DeleteUserResponse 删除用户响应
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// ListUsersResponse 列出用户响应
type ListUsersResponse struct {
	Users    []UserDTO `json:"users"`
	Total    int64     `json:"total"`
	Page     int32     `json:"page"`
	PageSize int32     `json:"page_size"`
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error string `json:"error"`
}
