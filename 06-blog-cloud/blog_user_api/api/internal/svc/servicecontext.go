// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"api/internal/config"
	"api/internal/middleware"
	"blog_user_service/rpc/userservice"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	AuthInterceptor rest.Middleware
	UserRpc         userservice.UserService
}

func NewServiceContext(c config.Config) *ServiceContext {
	client, err := zrpc.NewClient(c.UserRpc)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:          c,
		AuthInterceptor: middleware.NewAuthInterceptorMiddleware().Handle,
		UserRpc:         userservice.NewUserService(client),
	}
}
