# 博客系统后端服务

## 项目简介

这是一个基于Go语言开发的博客系统后端服务，使用Gin框架作为Web框架，GORM作为ORM框架，MySQL作为数据库，Redis作为缓存系统。

## 技术栈

- **语言**: Go
- **Web框架**: Gin
- **ORM框架**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **认证**: JWT

## 项目结构

```
04blog/
├── api/          # API处理器层，负责处理HTTP请求和响应
├── config/       # 配置管理，负责加载和管理系统配置
├── constant/     # 常量定义，包含Redis键常量等
├── middleware/   # 中间件，包含认证、日志、错误处理等
├── models/       # 数据模型层，定义数据库表结构
├── repositories/ # 数据访问层，负责与数据库交互
├── routes/       # 路由配置，定义API路由
├── servers/      # 业务逻辑层，实现核心业务逻辑
├── utils/        # 工具函数，包含JWT、密码加密等
├── go.mod        # Go模块定义
├── go.sum        # 依赖版本锁定
└── main.go       # 程序入口
```
## 开发步骤

1. **环境准备**
   - 安装Go语言环境
   - wls 环境安装MySQL数据库
   - wls 安装Redis缓存

2. **项目初始化**
   - 克隆项目仓库
   - 进入项目目录
   - 运行 `go mod init 04blog` `go mod tidy` 安装依赖
3. **项目初始化**
   - 运行 `go run main.go` 初始化项目,编写对应的目录结构，
   - 初始化 gin gorm 框架，配置路由，数据库连接，等
   - 编写 `models` 包，定义数据库表结构，通过配置进行数据库迁移
   - 按照 repositories - services - api 编写用户模块 完成进行测试 （添加JWT认证 ,将用户id放入上下文）
   - 按照 repositories - services - api 编写文章模块 完成进行测试
   - 按照 repositories - services - api 编写评论模块 完成进行测试
   - 添加redis缓存，实现文章阅读量统计（已完成），点赞量缓存，评论数缓存（待完成）

3. **数据库迁移**
   - 运行数据库迁移脚本，创建必要的数据库表

4. **配置项目**
   - 编辑 `config/config.yaml` 文件，配置数据库连接、Redis连接等

5. **启动服务**
   - 运行 `main.go` 文件，启动博客系统后端服务

## 核心功能

### 1. 用户管理
- 用户注册
- 用户登录（支持用户名/手机号登录）
- 获取用户信息

### 2. 文章管理
- 创建文章
- 获取文章详情（带阅读量统计）
- 获取文章列表（支持多条件查询）
- 更新文章
- 删除文章

### 3. 评论管理
- 发表评论
- 获取文章评论列表
- 删除评论

### 4. 点赞功能
- 文章点赞（Redis实现）

### 5. 缓存功能
- 文章阅读量统计（Redis自增）
- 文章点赞量缓存（Redis）
- 文章评论数缓存（Redis）

## 详细模块说明

### config包

配置管理模块，负责加载和管理系统配置，包括服务器配置、数据库配置和Redis配置。

**主要文件**:
- `config.go`: 定义配置结构和初始化函数

**核心功能**:
- 加载系统配置（默认值或环境变量）
- 初始化数据库连接
- 初始化Redis连接

### models包

数据模型层，定义数据库表结构和关系。

**主要文件**:
- `user.go`: 用户模型
- `post.go`: 文章模型
- `comment.go`: 评论模型
- `like.go`: 点赞模型

**核心功能**:
- 定义数据模型结构
- 定义表名和字段映射
- 支持软删除（通过BaseModel）

### repositories包

数据访问层，负责与数据库交互，实现数据的CRUD操作。

**主要文件**:
- `user_repository.go`: 用户数据访问
- `post_repository.go`: 文章数据访问
- `comment_repository.go`: 评论数据访问

**核心功能**:
- 封装数据库操作
- 提供统一的数据访问接口
- 支持条件查询

### servers包

业务逻辑层，实现核心业务逻辑，调用repositories层进行数据操作。

**主要文件**:
- `user_service.go`: 用户服务
- `post_service.go`: 文章服务
- `comment_service.go`: 评论服务

**核心功能**:
- 实现业务逻辑
- 调用数据访问层
- 与Redis缓存交互

### api包

API处理器层，负责处理HTTP请求和响应，调用servers层进行业务处理。

**主要文件**:
- `user_api.go`: 用户API
- `post_api.go`: 文章API
- `comment_api.go`: 评论API

**核心功能**:
- 处理HTTP请求
- 参数验证和绑定
- 调用业务逻辑层
- 返回响应结果

### middleware包

中间件层，提供各种中间件功能，如认证、日志、错误处理等。

**主要文件**:
- `auth.go`: 认证中间件
- `logger.go`: 日志中间件
- `error.go`: 错误处理中间件

**核心功能**:
- JWT认证
- 请求日志记录
- 统一错误处理

### routes包

路由配置，定义API路由和请求处理函数的映射关系。

**主要文件**:
- `routes.go`: 路由配置

**核心功能**:
- 定义API路由
- 配置中间件
- 依赖注入

### utils包

工具函数，提供各种辅助功能。

**主要文件**:
- `jwt.go`: JWT工具
- `md5Utils.go`: 密码加密工具

**核心功能**:
- JWT令牌生成和验证
- 密码加密和验证

### constant包

常量定义，包含系统中使用的各种常量。

**主要文件**:
- `redis_constants.go`: Redis键常量

**核心功能**:
- 定义Redis缓存键格式
- 统一常量管理

## API接口

### 用户接口

| 方法 | URL | 功能 | 认证 |
|------|-----|------|------|
| POST | /v1/register | 用户注册 | 否 |
| POST | /v1/login | 用户登录 | 否 |
| GET | /v1/user | 获取用户信息 | 是 |

### 文章接口

| 方法 | URL | 功能 | 认证 |
|------|-----|------|------|
| POST | /v1/post | 创建/更新文章 | 是 |
| GET | /v1/post/:id | 获取文章详情 | 否 |
| GET | /v1/posts | 获取文章列表 | 否 |
| DELETE | /v1/post/:id | 删除文章 | 是 |

### 评论接口

| 方法 | URL | 功能 | 认证 |
|------|-----|------|------|
| POST | /v1/comment | 发表评论 | 是 |
| GET | /v1/comments/:postID | 获取文章评论 | 否 |
| DELETE | /v1/comment/:commentID | 删除评论 | 是 |
| GET | /v1/comment/:commentID | 获取评论详情 | 否 |

## 配置说明

系统配置包括服务器配置、数据库配置和Redis配置，通过`config.LoadConfig()`函数加载。

### 服务器配置
- 端口：默认8080
- 主机：默认localhost

### 数据库配置
- 主机：172.18.112.82
- 端口：3306
- 用户名：root
- 密码：root
- 数据库名：moon_blog

### Redis配置
- 主机：172.18.112.82
- 端口：6379
- 密码：123456
- 数据库：0

## 运行项目

1. 确保已安装Go 1.18+
2. 确保MySQL和Redis服务正常运行
3. 执行以下命令：

```bash
# 安装依赖
go mod tidy

# 运行项目
go run main.go
```

## Redis缓存设计

### 文章阅读量
- 键格式：`post_views:{postID}`
- 类型：String
- 功能：使用自增操作记录文章阅读量

### 文章点赞量
- 键格式：`post_likes:{postID}`
- 类型：String
- 功能：记录文章点赞数量

### 文章评论数
- 键格式：`post_comments:{postID}`
- 类型：String
- 功能：记录文章评论数量

## 代码规范

1. **命名规范**：
   - 包名：使用小写字母
   - 结构体名：使用PascalCase
   - 函数名：使用PascalCase
   - 变量名：使用camelCase
   - 常量名：使用全大写，单词间用下划线分隔

2. **错误处理**：
   - 所有错误都应该被处理，不应该被忽略
   - 错误信息应该清晰、准确
   - Redis操作失败时应该有合理的降级策略

3. **代码组织**：
   - 遵循三层架构：API层 -> 服务层 -> 数据访问层
   - 使用依赖注入模式管理对象依赖
   - 接口和实现分离，便于测试和扩展

## 扩展建议

1. **增加单元测试和集成测试**
2. **添加Swagger文档**
3. **实现分页功能**
4. **添加搜索功能**
5. **实现图片上传功能**
6. **添加消息通知系统**
7. **实现标签功能**
8. **添加文章分类功能**
9. **实现邮件发送功能**
10. **添加用户关注功能**

## 部署说明

1. **编译项目**：
   ```bash
   go build -o blog_server
   ```

2. **配置环境变量**：
   - SERVER_PORT：服务器端口
   - SERVER_HOST：服务器主机
   - DB_HOST：数据库主机
   - DB_PORT：数据库端口
   - DB_USER：数据库用户名
   - DB_PASS：数据库密码
   - DB_NAME：数据库名
   - REDIS_HOST：Redis主机
   - REDIS_PORT：Redis端口
   - REDIS_PASSWORD：Redis密码

3. **运行服务**：
   ```bash
   ./blog_server
   ```

## 注意事项

1. 项目使用了JWT认证，需要在请求头中添加`Authorization: Bearer {token}`
2. 数据库连接信息和Redis连接信息需要根据实际环境修改
3. 生产环境中应该关闭Gin的Debug模式
4. 生产环境中应该使用HTTPS协议
5. 定期备份数据库，防止数据丢失
