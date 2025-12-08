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

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UserInfoVo) (*user.IsSuccess, error) {
	l.Logger.Infof("UpdateUserInfo request received: UserId=%d", in.Id)
	
	// 获取mysql配置
	db := inits.MysqlDb
	
	// 参数验证
	if in.Id <= 0 {
		l.Logger.Errorf("Invalid user ID: %d", in.Id)
		return &user.IsSuccess{Success: false}, errors.New("无效的用户ID")
	}
	
	// 检查用户是否存在
	var dbUser models.User
	result := db.First(&dbUser, in.Id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			l.Logger.Errorf("User not found for ID: %d", in.Id)
			return &user.IsSuccess{Success: false}, errors.New("用户不存在")
		}
		l.Logger.Errorf("Database query error: %v", result.Error)
		return &user.IsSuccess{Success: false}, result.Error
	}
	
	// 创建更新字段映射
	updates := make(map[string]interface{})
	
	// 只更新非空字段
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Email != "" {
		updates["email"] = in.Email
	}
	if in.Password != "" {
		// 密码加密
		hash := md5.Sum([]byte(in.Password))
		updates["password"] = hex.EncodeToString(hash[:])
	}
	if in.Phone != "" {
		updates["phone"] = in.Phone
	}
	if in.Gender > 0 {
		updates["gender"] = in.Gender
	}
	
	// 如果有需要更新的字段
	if len(updates) > 0 {
		result = db.Model(&dbUser).Updates(updates)
		if result.Error != nil {
			l.Logger.Errorf("Failed to update user: %v", result.Error)
			return &user.IsSuccess{Success: false}, result.Error
		}
		l.Logger.Infof("User updated successfully: ID=%d, updated fields=%v", in.Id, updates)
	} else {
		l.Logger.Infof("No fields provided to update for user: ID=%d", in.Id)
	}
	
	return &user.IsSuccess{Success: true}, nil
}
