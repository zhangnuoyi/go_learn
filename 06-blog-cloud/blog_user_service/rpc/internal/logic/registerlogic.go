package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"blog_user_service/rpc/inits"
	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/models"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 注册
func (l *RegisterLogic) Register(in *user.RegisterDto) (userVo *user.UserInfoVo, err error) {
	// 创建一个新的UserInfoVo实例
	userVo = new(user.UserInfoVo)
	l.Logger.Infof("Register: %v", in)
	//获取mysql配置
	db := inits.MysqlDb
	username := in.Name
	var existingUser models.User
	// 检查用户名是否已存在，如果不存在则会返回RecordNotFound错误
	result := db.Where("name = ?", username).First(&existingUser)
	if result.Error == nil {
		// 用户已存在
		return nil, errors.New("用户名已存在")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		// 非记录不存在的其他错误
		return nil, result.Error
	}
	// 用户名不存在，可以继续创建
	//密码加密
	hash := md5.Sum([]byte(in.Password))
	in.Password = hex.EncodeToString(hash[:])
	//创建用户
	user := models.User{
		Name:     in.Name,
		Password: in.Password,
		Email:    in.Email,
	}
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	// 设置UserInfoVo字段
	l.Logger.Infof("User ID from database: %d", user.ID)
	userVo.Id = uint32(user.ID)
	l.Logger.Infof("Set userVo.Id to: %d", userVo.Id)
	userVo.Name = user.Name
	userVo.Email = user.Email
	l.Logger.Infof("Returning userVo: Id=%d, Name=%s, Email=%s", userVo.Id, userVo.Name, userVo.Email)
	return userVo, nil
}
