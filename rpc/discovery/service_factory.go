package discovery

import (
	"context"
	"fmt"
	"sync"

	"google.golang.org/grpc"
	"github.com/user/rpc/proto"
)

// ServiceFactory 服务工厂接口
type ServiceFactory interface {
	// GetUserServiceClient 获取用户服务客户端
	GetUserServiceClient() (proto.UserServiceClient, error)
	// Close 关闭服务工厂
	Close() error
}

// GRPCServiceFactory gRPC 服务工厂实现
type GRPCServiceFactory struct {
	proxy     ServiceProxy
	userMutex sync.Mutex
	userClient proto.UserServiceClient
	userConn   *grpc.ClientConn
}

// NewGRPCServiceFactory 创建 gRPC 服务工厂实例
func NewGRPCServiceFactory(proxy ServiceProxy) *GRPCServiceFactory {
	return &GRPCServiceFactory{
		proxy: proxy,
	}
}

// GetUserServiceClient 获取用户服务客户端
func (f *GRPCServiceFactory) GetUserServiceClient() (proto.UserServiceClient, error) {
	// 先检查是否已有客户端
	if f.userClient != nil {
		return f.userClient, nil
	}

	// 加锁确保并发安全
	f.userMutex.Lock()
	defer f.userMutex.Unlock()

	// 双重检查
	if f.userClient != nil {
		return f.userClient, nil
	}

	// 获取连接
	conn, err := f.proxy.GetClientConn("user-service")
	if err != nil {
		return nil, fmt.Errorf("获取用户服务连接失败: %v", err)
	}

	// 创建客户端
	f.userConn = conn
	f.userClient = proto.NewUserServiceClient(conn)

	return f.userClient, nil
}

// Close 关闭服务工厂
func (f *GRPCServiceFactory) Close() error {
	// 关闭用户服务连接
	if f.userConn != nil {
		if err := f.userConn.Close(); err != nil {
			return fmt.Errorf("关闭用户服务连接失败: %v", err)
		}
	}

	// 关闭代理
	if err := f.proxy.Close(); err != nil {
		return fmt.Errorf("关闭服务代理失败: %v", err)
	}

	return nil
}

// LocalServiceFactory 本地服务工厂（用于开发和测试，直接连接到本地服务）
type LocalServiceFactory struct {
	userClient proto.UserServiceClient
	userConn   *grpc.ClientConn
}

// NewLocalServiceFactory 创建本地服务工厂实例
func NewLocalServiceFactory(userServiceAddr string) (*LocalServiceFactory, error) {
	// 连接到本地 gRPC 服务
	conn, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("连接到本地用户服务失败: %v", err)
	}

	// 创建客户端
	userClient := proto.NewUserServiceClient(conn)

	return &LocalServiceFactory{
		userClient: userClient,
		userConn:   conn,
	}, nil
}

// GetUserServiceClient 获取用户服务客户端
func (f *LocalServiceFactory) GetUserServiceClient() (proto.UserServiceClient, error) {
	return f.userClient, nil
}

// Close 关闭服务工厂
func (f *LocalServiceFactory) Close() error {
	if f.userConn != nil {
		if err := f.userConn.Close(); err != nil {
			return fmt.Errorf("关闭用户服务连接失败: %v", err)
		}
	}

	return nil
}
