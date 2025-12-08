package main

import (
	"flag"
	"fmt"

	"blog-post-service/rpc/inits"
	"blog-post-service/rpc/internal/config"
	"blog-post-service/rpc/internal/server"
	"blog-post-service/rpc/internal/svc"
	"blog-post-service/rpc/types/post"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// var configFile = flag.String("f", "etc/post.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad("etc/post.yaml", &c)
	ctx := svc.NewServiceContext(c)
	inits.InitDB()
	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		post.RegisterPostServiceServer(grpcServer, server.NewPostServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
