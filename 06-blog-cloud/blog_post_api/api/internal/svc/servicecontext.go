// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"06-blog-cloud/blog_post_api/api/internal/config"
	"blog-post-service/rpc/postservice"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config  config.Config
	PostRpc postservice.PostService
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := zrpc.NewClient(c.PostRpc)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:  c,
		PostRpc: postservice.NewPostService(client),
	}
}
