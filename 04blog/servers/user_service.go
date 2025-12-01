package servers

import (
	"04blog/models"
	"04blog/repositories"
	"04blog/utils"
	"errors"
)

// UserService 用户服务接口

type UserService interface {
	// RegisterUser 注册用户
	RegisterUser(user *models.User) error
	// LoginUser 登录用户
	LoginUser(username, password string) (string, *models.User, error)
	// GetUserByID 根据用户ID获取用户信息
	GetUserByID(userID int64) (*models.User, error)
}

// userService 用户服务实现

type userService struct {
	userDao repositories.UserRepository
}

func NewUserService(userDao repositories.UserRepository) UserService {
	return &userService{userDao: userDao}
}

//用户注册

func (u *userService) RegisterUser(user *models.User) error {
	//1.检查用户是否存在
	existingUser, err := u.userDao.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return errors.New("用户名已存在")
	}
	//检查邮箱是否存在（仅当邮箱不为空时）
	if user.Email != "" {
		existingUser, err = u.userDao.GetUserByEmail(user.Email)
		if err == nil && existingUser != nil {
			return errors.New("邮箱已注册")
		}
	}
	//密码加密
	user.Password = utils.GetMD5Hash(user.Password)

	//2.创建用户
	return u.userDao.CreateUser(user)
}

//用户登录

func (u *userService) LoginUser(usernameOrMobile string, password string) (string, *models.User, error) {
	//1.检查用户是否存在（先尝试用户名）
	existingUser, err := u.userDao.GetUserByUsername(usernameOrMobile)
	if err != nil || existingUser == nil {
		//尝试手机号
		existingUser, err = u.userDao.GetUserByMobile(usernameOrMobile)
		if err != nil || existingUser == nil {
			return "", nil, errors.New("用户名或密码错误")
		}
	}
	//2.检查密码是否正确
	if !utils.CheckPasswordHash(password, existingUser.Password) {
		return "", nil, errors.New("密码错误")
	}
	//3.生成token
	token, err := utils.GenerateToken(existingUser.ID, existingUser.Username)
	if err != nil {
		return "", nil, err
	}
	return token, existingUser, nil
}

// GetUserByID implements UserService.
func (u *userService) GetUserByID(userID int64) (*models.User, error) {
	return u.userDao.GetUserByID(userID)
}
