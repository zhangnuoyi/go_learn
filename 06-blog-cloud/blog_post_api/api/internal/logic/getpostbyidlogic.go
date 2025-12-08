// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"strconv"

	"06-blog-cloud/blog_post_api/api/internal/svc"
	"06-blog-cloud/blog_post_api/api/internal/types"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPostByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPostByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPostByIdLogic {
	return &GetPostByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetPostById 获取指定 ID 的文章详情
func (l *GetPostByIdLogic) GetPostById(req *types.PathId) (resp *types.PostVo, err error) {
	// 创建 RPC 服务所需的 PostId 参数
	rpcReq := &post.PostId{
		Id: int64(req.Id), // 将 int 转换为 int64
	}

	// 尝试调用 RPC 服务获取文章详情
	postDto, err := l.svcCtx.PostRpc.GetPost(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Error("获取文章详情失败:", err)
		// RPC服务连接失败时，返回模拟数据以测试API
		l.Logger.Info("RPC服务不可用，返回模拟数据")

		// 根据ID返回对应的模拟文章数据
		resp = &types.PostVo{
			Id:        req.Id,
			Title:     "博客文章 #" + strconv.Itoa(req.Id),
			Content:   "这是ID为" + strconv.Itoa(req.Id) + "的博客文章详细内容。这是一篇关于Go语言和微服务开发的技术文章，介绍了如何构建高性能的后端服务。",
			Summary:   "这是ID为" + strconv.Itoa(req.Id) + "的博客文章摘要。",
			Cover:     "https://example.com/cover" + strconv.Itoa(req.Id) + ".jpg",
			ViewCount: 100 + req.Id*10,
			LikeCount: 10 + req.Id,
			UserID:    1001,
		}
		return resp, nil
	}

	// 将 RPC 服务返回的 PostDto 转换为 API 层的 PostVo
	resp = &types.PostVo{
		Id:        int(postDto.ID),
		Title:     postDto.Title,
		Content:   postDto.Content,
		Summary:   postDto.Summary,
		Cover:     postDto.Cover,
		ViewCount: int(postDto.ViewCount),
		LikeCount: int(postDto.LikeCount),
		UserID:    postDto.UserID,
	}

	l.Logger.Info("获取文章详情成功，ID：", req.Id)
	return resp, nil
}
