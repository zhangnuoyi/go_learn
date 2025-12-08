package logic

import (
	"context"

	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByUsernameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByUsernameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByUsernameLogic {
	return &GetUserByUsernameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 根据用户名获取用户信息
func (l *GetUserByUsernameLogic) GetUserByUsername(in *user.UsernameDto) (*user.UserInfoVo, error) {
	// todo: add your logic here and delete this line

	return &user.UserInfoVo{}, nil
}
