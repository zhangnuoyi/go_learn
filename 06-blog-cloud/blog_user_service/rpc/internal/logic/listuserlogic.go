package logic

import (
	"context"

	"blog_user_service/rpc/inits"
	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/models"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListUserLogic {
	return &ListUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分页查询用户
func (l *ListUserLogic) ListUser(in *user.PageInfoDto) (*user.UserInfoVoList, error) {
	l.Logger.Infof("ListUser request received: Page=%d, Size=%d", in.PageNumber, in.PageSize)

	// 获取mysql配置
	db := inits.MysqlDb

	// 参数验证
	pageNumber := in.PageNumber
	if pageNumber <= 0 {
		pageNumber = 1
	}

	pageSize := in.PageSize
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}

	offset := (pageNumber - 1) * pageSize

	// 查询总数
	var total int64
	if err := db.Model(&models.User{}).Count(&total).Error; err != nil {
		l.Logger.Errorf("Failed to count users: %v", err)
		return nil, err
	}

	// 查询数据
	var users []models.User
	if err := db.Offset(int(offset)).Limit(int(pageSize)).Find(&users).Error; err != nil {
		l.Logger.Errorf("Failed to list users: %v", err)
		return nil, err
	}

	// 转换为响应格式
	userInfos := make([]*user.UserInfoVo, len(users))
	for i, u := range users {
		userInfos[i] = &user.UserInfoVo{
			Id:     uint32(u.ID),
			Name:   u.Name,
			Email:  u.Email,
			Gender: 1,
			Phone:  "13555554444",
		}
		l.Logger.Infof("User found: ID=%d, Name=%s, Email=%s", u.ID, u.Name, u.Email)
	}

	l.Logger.Infof("ListUser request completed: Total=%d, Returned=%d", total, len(userInfos))

	// 返回结果
	return &user.UserInfoVoList{
			Total:          int32(total),
			UserInfoVoList: userInfos,
		},
		nil
}
