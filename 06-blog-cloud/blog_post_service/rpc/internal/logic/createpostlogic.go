package logic

import (
	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/models"
	"blog-post-service/rpc/types/post"
	"context"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePostLogic {
	return &CreatePostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 新增文章
func (l *CreatePostLogic) CreatePost(in *post.PostDto) (*post.IsSuccess, error) {
	// 将PostDto转换为Post模型
	postModel := &models.Post{
		Title:     in.Title,
		Content:   in.Content,
		Summary:   in.Summary,
		Cover:     in.Cover,
		UserID:    in.UserID,
		ViewCount: int(in.ViewCount),
		LikeCount: int(in.LikeCount),
	}

	// 使用GORM保存到数据库
	result := inits.MysqlDb.Create(postModel)
	if result.Error != nil {
		l.Error("创建文章失败：", result.Error)
		return &post.IsSuccess{Success: false}, result.Error
	}

	l.Info("创建文章成功，ID：", postModel.ID)
	return &post.IsSuccess{Success: true}, nil
}
