// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"api/internal/svc"
	"api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PathExampleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPathExampleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PathExampleLogic {
	return &PathExampleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PathExampleLogic) PathExample(req *types.QueryPathDto) (resp *types.QueryVo, err error) {
	// todo: add your logic here and delete this line

	return
}
