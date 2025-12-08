package main

import (
	"fmt"

	"blog_user_service/rpc/inits"
	"blog_user_service/rpc/internal/config"
	"blog_user_service/rpc/internal/server"
	"blog_user_service/rpc/internal/svc"
	"blog_user_service/rpc/types/user"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var c config.Config
	conf.MustLoad("etc/config.yaml", &c)
	ctx := svc.NewServiceContext(c)
	//初始化数据库
	inits.InitDB()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		user.RegisterUserServiceServer(grpcServer, server.NewUserServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
