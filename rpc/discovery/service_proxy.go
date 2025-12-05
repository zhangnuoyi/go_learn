package discovery

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// ServiceProxy 服务代理接口
type ServiceProxy interface {
	// GetClientConn 获取 gRPC 客户端连接
	GetClientConn(serviceName string) (*grpc.ClientConn, error)
	// Close 关闭服务代理
	Close() error
}

// GRPCServiceProxy gRPC 服务代理实现
type GRPCServiceProxy struct {
	discovery   ServiceDiscovery
	clients     map[string]*grpc.ClientConn
	mutex       sync.RWMutex
	connTimeout time.Duration
	retryCount  int
}

// NewGRPCServiceProxy 创建 gRPC 服务代理实例
func NewGRPCServiceProxy(discovery ServiceDiscovery) *GRPCServiceProxy {
	return &GRPCServiceProxy{
		discovery:   discovery,
		clients:     make(map[string]*grpc.ClientConn),
		connTimeout: 5 * time.Second,
		retryCount:  3,
	}
}

// GetClientConn 获取 gRPC 客户端连接
func (p *GRPCServiceProxy) GetClientConn(serviceName string) (*grpc.ClientConn, error) {
	if serviceName == "" {
		return nil, fmt.Errorf("服务名称不能为空")
	}

	// 先检查是否已有连接
	p.mutex.RLock()
	conn, exists := p.clients[serviceName]
	p.mutex.RUnlock()

	if exists && !conn.GetState().String() == "TRANSIENT_FAILURE" {
		return conn, nil
	}

	// 获取服务实例
	instances, err := p.discovery.GetService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("获取服务实例失败: %v", err)
	}

	if len(instances) == 0 {
		return nil, fmt.Errorf("没有可用的服务实例: %s", serviceName)
	}

	// 使用第一个实例（实际应用中可能需要负载均衡）
	instance := instances[0]
	target := fmt.Sprintf("%s:%d", instance.Address, instance.Port)

	// 创建连接选项
	options := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(p.connTimeout),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10 * 1024 * 1024), // 10MB
			grpc.MaxCallSendMsgSize(10 * 1024 * 1024), // 10MB
		),
	}

	// 创建连接
	conn, err = p.createConnection(target, options)
	if err != nil {
		return nil, fmt.Errorf("创建连接失败: %v", err)
	}

	// 更新连接缓存
	p.mutex.Lock()
	p.clients[serviceName] = conn
	p.mutex.Unlock()

	return conn, nil
}

// createConnection 创建 gRPC 连接，带重试机制
func (p *GRPCServiceProxy) createConnection(target string, options []grpc.DialOption) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < p.retryCount; i++ {
		log.Printf("尝试连接到服务: %s, 重试次数: %d/%d", target, i+1, p.retryCount)
		conn, err = grpc.Dial(target, options...)
		if err == nil {
			log.Printf("成功连接到服务: %s", target)
			return conn, nil
		}

		log.Printf("连接服务失败: %s, 错误: %v", target, err)
		if i < p.retryCount-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	return nil, fmt.Errorf("重试连接失败: %v", err)
}

// Close 关闭服务代理
func (p *GRPCServiceProxy) Close() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	var lastErr error
	for serviceName, conn := range p.clients {
		if err := conn.Close(); err != nil {
			log.Printf("关闭服务连接失败: %s, 错误: %v", serviceName, err)
			lastErr = err
		} else {
			log.Printf("成功关闭服务连接: %s", serviceName)
		}
	}

	p.clients = make(map[string]*grpc.ClientConn)
	return lastErr
}
