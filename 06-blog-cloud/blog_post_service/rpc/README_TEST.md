# 帖子服务RPC接口测试指南

本文档提供了测试帖子服务RPC模块接口的方法和步骤。

## 测试工具

- `test_post_client.go`: 测试客户端程序，用于测试帖子服务的所有RPC接口

## 测试步骤

### 1. 确保帖子服务RPC服务端已启动

帖子服务RPC服务端应该在端口8081上运行。如果服务未启动，请运行以下命令：

```bash
cd h:\go\work\06-blog-cloud\blog_post_service\rpc
go run post.go
```

### 2. 运行测试客户端

在另一个终端中运行测试客户端：

```bash
cd h:\go\work\06-blog-cloud\blog_post_service\rpc\client
go run test_post_client.go
```

## 测试内容

测试客户端会依次测试以下接口：

1. **创建帖子**: 测试`CreatePost`接口，创建一个测试帖子
2. **获取帖子列表**: 测试`GetPostsByConditions`接口，获取帖子列表
3. **获取单个帖子**: 测试`GetPost`接口，获取指定ID的帖子详情
4. **更新帖子**: 测试`UpdatePost`接口，更新指定ID的帖子
5. **删除帖子**: 测试`DeletePost`接口，删除指定ID的帖子

## 注意事项

- 确保在运行测试客户端前，帖子服务RPC服务端已经成功启动
- 测试客户端默认连接到`0.0.0.0:8081`，如果服务端配置有变化，请相应修改测试客户端的连接地址
- 测试使用的帖子ID默认为1，可能需要根据实际情况调整
