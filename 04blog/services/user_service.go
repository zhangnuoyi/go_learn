package services

import (
	"04blog/models"
	"04blog/repositories"
	"04blog/utils"
	"context"
	"errors"
)

// UserInfo 用户信息结构体
type UserInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Mobile   string `json:"mobile"`
	Email    string `json:"email"`
	Role     int    `json:"role"`
}

// UserService 用户服务接口
type UserService interface {
	// RegisterUser 注册用户
	RegisterUser(ctx context.Context, user *models.User) error
	// LoginUser 登录用户
	LoginUser(ctx context.Context, username, password string) (string, *UserInfo, error)
	// GetUserInfo 获取用户信息
	GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// RegisterUser 注册用户
func (s *userService) RegisterUser(ctx context.Context, user *models.User) error {
	//1.检查用户是否存在
	existingUser, err := s.userRepo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}
	//检查邮箱是否存在
	if user.Email != "" {
		existingUser, err = s.userRepo.GetUserByEmail(user.Email)
		if err == nil && existingUser != nil {
			return errors.New("邮箱已注册")
		}
	}
	//密码加密
	user.Password = utils.GetMD5Hash(user.Password)

	//2.创建用户
	return s.userRepo.CreateUser(user)
}

// LoginUser 登录用户
func (s *userService) LoginUser(ctx context.Context, usernameOrMobile, password string) (string, *UserInfo, error) {
	//1.检查用户是否存在（先尝试用户名）
	existingUser, err := s.userRepo.GetUserByUsername(usernameOrMobile)
	if err != nil || existingUser == nil {
		//尝试手机号
		existingUser, err = s.userRepo.GetUserByMobile(usernameOrMobile)
		if err != nil || existingUser == nil {
			return "", nil, errors.New("用户名或密码错误")
		}
	}
	//2.检查密码是否正确
	if !utils.CheckPasswordHash(password, existingUser.Password) {
		return "", nil, errors.New("用户名或密码错误")
	}
	//3.生成token
	token, err := utils.GenerateToken(existingUser.ID, existingUser.Username)
	if err != nil {
		return "", nil, err
	}
	//4.构建返回信息
	userInfo := &UserInfo{
		ID:       existingUser.ID,
		Username: existingUser.Username,
		Nickname: existingUser.Nickname,
		Mobile:   existingUser.Mobile,
		Email:    existingUser.Email,
		Role:     existingUser.Role,
	}
	return token, userInfo, nil
}

// GetUserInfo 获取用户信息
func (s *userService) GetUserInfo(ctx context.Context, userID int64) (*UserInfo, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		Mobile:   user.Mobile,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
