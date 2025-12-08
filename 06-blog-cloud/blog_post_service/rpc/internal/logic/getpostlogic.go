package logic

import (
	"context"

	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/models"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostLogic {
	return &GetPostLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取文章
func (l *GetPostLogic) GetPost(in *post.PostId) (*post.PostDto, error) {
	// 声明一个Post模型变量用于存储查询结果
	var postModel models.Post

	// 使用GORM根据ID查询文章
	result := inits.MysqlDb.First(&postModel, in.Id)
	if result.Error != nil {
		l.Error("查询文章失败：", result.Error)
		return nil, result.Error
	}

	// 将Post模型转换为PostDto
	postDto := &post.PostDto{
		ID:        postModel.ID,
		Title:     postModel.Title,
		Content:   postModel.Content,
		Summary:   postModel.Summary,
		Cover:     postModel.Cover,
		ViewCount: int64(postModel.ViewCount),
		LikeCount: int64(postModel.LikeCount),
		UserID:    postModel.UserID,
	}

	l.Info("查询文章成功，ID：", in.Id)
	return postDto, nil
}
