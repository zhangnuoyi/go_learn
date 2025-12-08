package logic

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"

	"blog-common/utils"
	"blog_user_service/rpc/inits"
	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/models"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 登录
func (l *LoginLogic) Login(in *user.LoginDto) (*user.LoginResponse, error) {
	l.Logger.Infof("Login request received: Username=%s", in.Username)

	//获取mysql配置
	db := inits.MysqlDb

	// 查询用户
	var dbUser models.User
	result := db.Where("name = ?", in.Username).First(&dbUser)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			l.Logger.Errorf("User not found: %s", in.Username)
			return nil, errors.New("用户名或密码错误")
		}
		l.Logger.Errorf("Database query error: %v", result.Error)
		return nil, result.Error
	}

	// 密码加密验证
	hash := md5.Sum([]byte(in.Password))
	encryptedPassword := hex.EncodeToString(hash[:])

	if encryptedPassword != dbUser.Password {
		l.Logger.Errorf("Password mismatch for user: %s", in.Username)
		return nil, errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	token, err := utils.GenerateToken(dbUser.ID, dbUser.Name)
	if err != nil {
		l.Logger.Errorf("Failed to generate token: %v", err)
		return nil, errors.New("生成令牌失败")
	}

	l.Logger.Infof("Login successful for user: %s, generated token", in.Username)

	// 返回登录响应
	return &user.LoginResponse{
		Token: token,
	}, nil
}
