// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"blog_user_service/rpc/types/user"
	"context"
	"errors"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterDto) (*types.RegisterResponse, error) {
	// 记录请求信息
	l.Logger.Infof("Register request received: Name=%s, Email=%s", req.Name, req.Email)

	// 创建RPC请求对象
	rpcReq := &user.RegisterDto{
		Name:     req.Name,
		Password: req.Password,
		Email:    req.Email,
	}
	l.Logger.Infof("Prepared RPC request: %+v", rpcReq)

	// 调用RPC服务
	userInfo, err := l.svcCtx.UserRpc.Register(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("RPC Register failed: %v", err)
		return nil, err
	}

	// 检查RPC响应
	if userInfo == nil {
		l.Logger.Error("RPC returned nil userInfo")
		return nil, errors.New("failed to register user")
	}

	l.Logger.Infof("RPC Register successful, received userInfo: Id=%d, Name=%s, Email=%s",
		userInfo.Id, userInfo.Name, userInfo.Email)

	// 创建RegisterResponse响应对象
	resp := &types.RegisterResponse{
		Id:    int64(userInfo.Id),
		Name:  userInfo.Name,
		Email: userInfo.Email,
	}

	l.Logger.Infof("Created response object: %+v", resp)
	return resp, nil
}
