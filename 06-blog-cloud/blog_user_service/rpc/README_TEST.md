# 用户RPC模块接口测试指南

本目录包含测试用户RPC模块接口的相关文件。

## 测试工具说明

1. `client/test_client.go` - 测试客户端程序，用于调用用户RPC接口
2. `run_tests.bat` - 一键测试脚本，自动启动服务端和测试客户端

## 测试步骤

### 方法一：使用一键测试脚本

1. 确保已经安装Go环境
2. 双击运行 `run_tests.bat` 脚本
3. 脚本会自动：
   - 启动用户服务RPC服务端
   - 等待服务端启动完成
   - 运行测试客户端调用各个RPC接口

### 方法二：手动测试

1. 启动服务端：
   ```bash
   cd h:\go\work\06-blog-cloud\blog_user_service\rpc
   go run user.go
   ```

2. 运行测试客户端：
   ```bash
   cd h:\go\work\06-blog-cloud\blog_user_service\rpc\client
   go run test_client.go
   ```

## 测试内容

测试客户端会按顺序测试以下接口：

1. **注册接口 (Register)** - 注册一个新用户
2. **登录接口 (Login)** - 使用刚注册的用户登录
3. **分页查询接口 (ListUser)** - 查询用户列表
4. **根据ID查询接口 (GetUserById)** - 根据用户ID查询详情
5. **更新用户信息接口 (UpdateUserInfo)** - 更新用户信息

## 注意事项

1. 确保服务端在8080端口正常运行
2. 测试程序会自动生成唯一的测试用户名和邮箱
3. 测试结果会在控制台输出详细信息

如果需要修改测试参数，请编辑 `client/test_client.go` 文件。
