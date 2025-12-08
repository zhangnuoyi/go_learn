package logic

import (
	"context"

	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/models"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePostLogic {
	return &UpdatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新文章
func (l *UpdatePostLogic) UpdatePost(in *post.PostDto) (*post.IsSuccess, error) {
	// 先检查文章是否存在
	var existingPost models.Post

	// 构建更新数据
	updateData := map[string]interface{}{
		"title":      in.Title,
		"content":    in.Content,
		"summary":    in.Summary,
		"cover":      in.Cover,
		"view_count": in.ViewCount,
		"like_count": in.LikeCount,
		"user_id":    in.UserID,
	}
	result := inits.MysqlDb.First(&existingPost, in.ID)
	if result.Error != nil {
		//保存文章
		inits.MysqlDb.Model(&existingPost).Save(updateData)
		return &post.IsSuccess{Success: false}, result.Error
	}

	// 执行更新操作
	result = inits.MysqlDb.Model(&existingPost).Updates(updateData)
	if result.Error != nil {
		l.Error("更新文章失败：", result.Error)
		return &post.IsSuccess{Success: false}, result.Error
	}

	// 检查是否真的有行被更新
	if result.RowsAffected == 0 {
		l.Errorf("未找到需要更新的文章，ID：", in.ID)
		return &post.IsSuccess{Success: false}, nil
	}

	l.Info("更新文章成功，ID：", in.ID)
	return &post.IsSuccess{Success: true}, nil
}
