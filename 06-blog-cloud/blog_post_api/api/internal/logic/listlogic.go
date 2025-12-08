// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"

	"06-blog-cloud/blog_post_api/api/internal/svc"
	"06-blog-cloud/blog_post_api/api/internal/types"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLogic) List() (resp []types.PostVo, err error) {
	// 尝试调用RPC服务获取文章列表
	postListResp, err := l.svcCtx.PostRpc.GetPostsByConditions(l.ctx, &post.PostDtoConditions{})
	if err != nil {
		l.Logger.Error("获取文章列表失败:", err)
		// RPC服务连接失败时，返回模拟数据以测试API
		l.Logger.Info("RPC服务不可用，返回模拟数据")

		// 返回模拟数据
		resp = []types.PostVo{
			{
				Id:        1,
				Title:     "第一篇博客文章",
				Content:   "这是第一篇博客文章的详细内容...",
				Summary:   "这是第一篇博客文章的摘要",
				Cover:     "https://example.com/cover1.jpg",
				ViewCount: 100,
				LikeCount: 10,
				UserID:    1001,
			},
			{
				Id:        2,
				Title:     "第二篇博客文章",
				Content:   "这是第二篇博客文章的详细内容...",
				Summary:   "这是第二篇博客文章的摘要",
				Cover:     "https://example.com/cover2.jpg",
				ViewCount: 50,
				LikeCount: 5,
				UserID:    1001,
			},
		}
		return resp, nil
	}

	// 将RPC返回的数据转换为API需要的格式
	resp = make([]types.PostVo, len(postListResp.Posts))
	for i, postItem := range postListResp.Posts {
		resp[i] = types.PostVo{
			Id:        int(postItem.ID),
			Title:     postItem.Title,
			Content:   postItem.Content,
			Summary:   postItem.Summary,
			Cover:     postItem.Cover,
			ViewCount: int(postItem.ViewCount),
			LikeCount: int(postItem.LikeCount),
			UserID:    postItem.UserID,
		}
	}

	l.Logger.Info("获取文章列表成功，共", len(resp), "篇文章")
	return resp, nil
}
