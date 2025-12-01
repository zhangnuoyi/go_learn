package repositories

import (
	"04blog/models"

	"gorm.io/gorm"
)

// UserRepository 用户数据访问接口

type UserRepository interface {
	// CreateUser 创建用户
	CreateUser(user *models.User) error
	// GetUserByID 根据用户ID获取用户
	GetUserByID(userID int64) (*models.User, error)
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(username string) (*models.User, error)
	// GetUserByMobile 根据手机号获取用户
	GetUserByMobile(mobile string) (*models.User, error)
	// UpdateUser 更新用户信息
	UpdateUser(user *models.User) error
	// DeleteUser 删除用户
	DeleteUser(userID int64) error
	// GetUserByEmail 根据邮箱获取用户
	GetUserByEmail(email string) (*models.User, error)
}

// userRepository 用户数据访问实现

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户数据访问实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户

func (u *userRepository) CreateUser(user *models.User) error {
	return u.db.Create(user).Error
}

// GetUserByID 根据用户ID获取用户
func (u *userRepository) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	if err := u.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名获取用户
func (u *userRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func (u *userRepository) UpdateUser(user *models.User) error {
	return u.db.Save(user).Error
}

// DeleteUser 删除用户
func (u *userRepository) DeleteUser(userID int64) error {
	return u.db.Delete(&models.User{}, userID).Error
}

// GetUserByEmail 根据邮箱获取用户
func (u *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByMobile 根据手机号获取用户
func (u *userRepository) GetUserByMobile(mobile string) (*models.User, error) {
	var user models.User
	if err := u.db.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
