package logic

import (
	"context"

	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/models"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePostLogic {
	return &DeletePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 删除文章
func (l *DeletePostLogic) DeletePost(in *post.PostId) (*post.IsSuccess, error) {
	// 检查文章是否存在
	var existingPost models.Post
	result := inits.MysqlDb.First(&existingPost, in.Id)
	if result.Error != nil {
		l.Error("文章不存在或查询失败：", result.Error)
		return &post.IsSuccess{Success: false}, result.Error
	}

	// 执行删除操作（软删除，因为模型中定义了DeletedAt字段）
	result = inits.MysqlDb.Delete(&existingPost)
	if result.Error != nil {
		l.Error("删除文章失败：", result.Error)
		return &post.IsSuccess{Success: false}, result.Error
	}

	// 检查是否真的有行被删除
	if result.RowsAffected == 0 {
		l.Errorf("未找到需要删除的文章，ID：", in.Id)
		return &post.IsSuccess{Success: false}, nil
	}

	l.Info("删除文章成功，ID：", in.Id)
	return &post.IsSuccess{Success: true}, nil
}
