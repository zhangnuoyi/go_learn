package logic

import (
	"context"

	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/models"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostsByConditionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPostsByConditionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostsByConditionsLogic {
	return &GetPostsByConditionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 多条件查询文章 不传入就是查询所有
func (l *GetPostsByConditionsLogic) GetPostsByConditions(in *post.PostDtoConditions) (*post.PostDtoList, error) {
	// 声明一个Post模型切片用于存储查询结果
	var postModels []models.Post
	var total int64

	// 构建查询条件
	query := inits.MysqlDb.Model(&models.Post{})

	// 根据条件过滤
	if in.Title != "" {
		query = query.Where("title LIKE ?", "%"+in.Title+"%")
	}
	if in.Content != "" {
		query = query.Where("content LIKE ?", "%"+in.Content+"%")
	}
	if in.UserID > 0 {
		query = query.Where("user_id = ?", in.UserID)
	}

	// 统计符合条件的总数
	if err := query.Count(&total).Error; err != nil {
		l.Error("统计文章总数失败：", err)
		return nil, err
	}

	// 执行查询
	if err := query.Find(&postModels).Error; err != nil {
		l.Error("查询文章列表失败：", err)
		return nil, err
	}

	// 将Post模型切片转换为PostDto切片
	postDtos := make([]*post.PostDto, 0, len(postModels))
	for _, postModel := range postModels {
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
		postDtos = append(postDtos, postDto)
	}

	l.Info("查询文章列表成功，共找到：", total, "篇")
	return &post.PostDtoList{
		Posts: postDtos,
		Total: int32(total),
	}, nil
}
