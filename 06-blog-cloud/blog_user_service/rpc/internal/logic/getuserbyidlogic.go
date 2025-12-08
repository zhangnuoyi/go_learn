package logic

import (
	"context"
	"errors"

	"blog_user_service/rpc/inits"
	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/models"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据id查询用户
func (l *GetUserByIdLogic) GetUserById(in *user.IdDto) (*user.UserInfoVo, error) {
	l.Logger.Infof("GetUserById request received: Id=%d", in.Id)

	// 获取mysql配置
	db := inits.MysqlDb

	// 参数验证
	if in.Id <= 0 {
		l.Logger.Errorf("Invalid user ID: %d", in.Id)
		return nil, errors.New("无效的用户ID")
	}

	// 查询用户
	var dbUser models.User
	result := db.First(&dbUser, in.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			l.Logger.Errorf("User not found for ID: %d", in.Id)
			return nil, errors.New("用户不存在")
		}
		l.Logger.Errorf("Database query error: %v", result.Error)
		return nil, result.Error
	}

	// 创建用户信息响应
	userVo := &user.UserInfoVo{
		Id:     uint32(dbUser.ID),
		Name:   dbUser.Name,
		Email:  dbUser.Email,
		Gender: 1,
		Phone:  "13555554444",
	}

	l.Logger.Infof("User found: ID=%d, Name=%s, Email=%s", dbUser.ID, dbUser.Name, dbUser.Email)

	return userVo, nil
}
