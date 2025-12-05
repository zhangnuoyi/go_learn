# 微服务架构案例：RESTful API + gRPC

## 项目概述

本项目是一个基于 Go 语言实现的微服务架构案例，结合了 RESTful API 和 gRPC 技术，旨在使内部系统多模块的调用像本地调用一样便捷。通过服务发现和调用封装层，实现了服务之间的无缝通信。

## 架构设计

### 核心技术栈

- **语言**: Go 1.18+
- **Web框架**: Gin
- **RPC框架**: gRPC
- **服务发现**: 本地服务发现（支持扩展到 Consul、Etcd 等）
- **数据存储**: 内存存储（支持扩展到数据库）

### 架构图

```
+----------------+     +----------------+     +----------------+
|   客户端应用    |     |   RESTful API  |     |   gRPC 服务    |
|                | <--> |                | <--> |                |
|  (client/main) |     |  (api/server)  |     | (server/main)  |
+----------------+     +----------------+     +----------------+
                                ^                     ^
                                |                     |
                        +----------------+     +----------------+
                        | 服务工厂 & 代理  |     |   业务逻辑层    |
                        | (discovery/*)   |     | (services/*)   |
                        +----------------+     +----------------+
                                ^                     ^
                                |                     |
                        +----------------+     +----------------+
                        |    服务发现     |     |   数据模型层    |
                        | (discovery/*)   |     |  (models/*)    |
                        +----------------+     +----------------+
```

## 目录结构

```
rpc/
├── api/                # RESTful API 接口
│   ├── main.go         # API 服务器入口
│   ├── server.go       # API 服务器实现
│   ├── user_handler.go # 用户相关 API 处理器
│   └── models.go       # API 数据模型
├── client/             # 客户端示例
│   └── main.go         # 客户端入口
├── discovery/          # 服务发现和调用封装
│   ├── service_discovery.go # 服务发现接口和实现
│   ├── service_proxy.go     # 服务代理实现
│   └── service_factory.go   # 服务工厂实现
├── models/             # 数据模型
│   ├── user.go         # 用户模型和存储
│   └── user_test.go    # 用户模型测试
├── proto/              # Protocol Buffers 定义
│   ├── user.proto      # 用户服务 proto 定义
│   ├── user.pb.go      # 自动生成的消息代码
│   └── user_grpc.pb.go # 自动生成的 gRPC 代码
├── services/           # 业务逻辑层
│   └── user_service.go # 用户服务实现
├── server/             # gRPC 服务器
│   └── main.go         # gRPC 服务器入口
├── go.mod              # Go 模块定义
└── README.md           # 项目说明文档
```

## 安装和运行

### 环境要求

- Go 1.18 或更高版本
- gRPC 相关依赖
- Gin Web 框架

### 安装依赖

```bash
go mod tidy
```

### 运行服务

#### 1. 启动 gRPC 服务器

```bash
cd server
go run main.go
```

gRPC 服务器将在 `localhost:50051` 端口启动。

#### 2. 启动 RESTful API 服务器

```bash
cd api
go run main.go
```

RESTful API 服务器将在 `localhost:8080` 端口启动。

#### 3. 运行客户端示例

```bash
cd client
go run main.go
```

客户端示例将演示如何调用 gRPC 服务。

### 运行测试

```bash
go test ./models -v
```

## API 接口说明

### RESTful API

#### 用户相关接口

| 方法 | URL | 描述 | 请求体 (JSON) | 成功响应 (200 OK) |
|------|-----|------|---------------|-------------------|
| POST | /api/users | 创建用户 | `{"name": "张三", "email": "zhangsan@example.com", "age": 30}` | `{"id": 1, "name": "张三", "email": "zhangsan@example.com", "age": 30}` |
| GET | /api/users/:id | 获取用户 | N/A | `{"id": 1, "name": "张三", "email": "zhangsan@example.com", "age": 30}` |
| PUT | /api/users/:id | 更新用户 | `{"name": "张三（更新）", "email": "zhangsan_updated@example.com", "age": 31}` | `{"id": 1, "name": "张三（更新）", "email": "zhangsan_updated@example.com", "age": 31}` |
| DELETE | /api/users/:id | 删除用户 | N/A | `{"message": "用户删除成功"}` |
| GET | /api/users | 获取用户列表 | N/A | `{"total": 5, "users": [{"id": 1, "name": "张三", ...}, ...]}` |

### gRPC API

#### 用户服务接口

```protobuf
service UserService {
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);
  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);
  rpc ListUsers(ListUsersRequest) returns (ListUsersResponse);
}
```

## 使用示例

### 1. 创建用户

#### RESTful API

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name": "张三", "email": "zhangsan@example.com", "age": 30}'
```

#### gRPC 客户端

```go
// 创建服务工厂
factory, _ := discovery.NewLocalServiceFactory("localhost:50051")

// 获取用户服务客户端
userClient, _ := factory.GetUserServiceClient()

// 调用创建用户方法
createResp, _ := userClient.CreateUser(ctx, &proto.CreateUserRequest{
    Name:  "张三",
    Email: "zhangsan@example.com",
    Age:   30,
})
```

### 2. 获取用户

#### RESTful API

```bash
curl http://localhost:8080/api/users/1
```

#### gRPC 客户端

```go
// 调用获取用户方法
getResp, _ := userClient.GetUser(ctx, &proto.GetUserRequest{
    Id: 1,
})
```

### 3. 更新用户

#### RESTful API

```bash
curl -X PUT http://localhost:8080/api/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "张三（更新）", "email": "zhangsan_updated@example.com", "age": 31}'
```

#### gRPC 客户端

```go
// 调用更新用户方法
updateResp, _ := userClient.UpdateUser(ctx, &proto.UpdateUserRequest{
    Id:    1,
    Name:  "张三（更新）",
    Email: "zhangsan_updated@example.com",
    Age:   31,
})
```

### 4. 删除用户

#### RESTful API

```bash
curl -X DELETE http://localhost:8080/api/users/1
```

#### gRPC 客户端

```go
// 调用删除用户方法
_, _ = userClient.DeleteUser(ctx, &proto.DeleteUserRequest{
    Id: 1,
})
```

### 5. 获取用户列表

#### RESTful API

```bash
curl http://localhost:8080/api/users?page=1&page_size=10
```

#### gRPC 客户端

```go
// 调用获取用户列表方法
listResp, _ := userClient.ListUsers(ctx, &proto.ListUsersRequest{
    Page:     1,
    PageSize: 10,
})
```

## 服务发现和调用封装

本项目实现了一个灵活的服务发现和调用封装层，主要包括：

1. **服务发现接口**：定义了服务注册、注销和查询的标准接口
2. **本地服务发现实现**：适用于开发和测试环境
3. **服务代理**：封装了 gRPC 客户端连接的管理和重试逻辑
4. **服务工厂**：提供了获取各种服务客户端的统一入口

### 使用服务工厂

```go
// 创建本地服务工厂
factory, err := discovery.NewLocalServiceFactory("localhost:50051")
if err != nil {
    log.Fatalf("创建服务工厂失败: %v", err)
}
defer factory.Close()

// 获取用户服务客户端（像调用本地方法一样简单）
userClient, err := factory.GetUserServiceClient()
if err != nil {
    log.Fatalf("获取用户服务客户端失败: %v", err)
}

// 直接调用服务方法
resp, err := userClient.GetUser(ctx, &proto.GetUserRequest{Id: 1})
```

## 扩展和定制

### 添加新的服务

1. 在 `proto/` 目录下创建新的 proto 文件
2. 在 `models/` 目录下实现数据模型
3. 在 `services/` 目录下实现业务逻辑
4. 在 `api/` 目录下添加 RESTful API 接口
5. 在 `discovery/` 目录下更新服务工厂

### 切换到实际的服务发现

本项目使用本地服务发现进行开发和测试，可以轻松切换到实际的服务发现系统：

```go
// 使用 Consul 作为服务发现（示例代码）
func NewConsulServiceDiscovery(addr string) (ServiceDiscovery, error) {
    // 实现 Consul 服务发现
}

// 使用 Etcd 作为服务发现（示例代码）
func NewEtcdServiceDiscovery(endpoints []string) (ServiceDiscovery, error) {
    // 实现 Etcd 服务发现
}
```

### 切换到数据库存储

本项目使用内存存储进行开发和测试，可以轻松切换到数据库存储：

```go
// 使用 GORM 连接数据库（示例代码）
func NewDatabaseUserStore(db *gorm.DB) UserStore {
    // 实现数据库用户存储
}
```

## 注意事项

1. **开发环境**：本项目使用内存存储和本地服务发现，仅适用于开发和测试环境
2. **生产环境**：生产环境中应使用实际的数据库和服务发现系统（如 Consul、Etcd 等）
3. **错误处理**：示例代码简化了错误处理，生产环境中应加强错误处理和日志记录
4. **并发安全**：服务发现和调用封装层已实现并发安全，但在使用时仍需注意并发场景
5. **性能优化**：对于高并发场景，可考虑添加连接池、缓存等优化措施

## 许可证

本项目采用 MIT 许可证，详见 LICENSE 文件。

## 联系方式

如有问题或建议，欢迎提交 Issue 或 Pull Request。
