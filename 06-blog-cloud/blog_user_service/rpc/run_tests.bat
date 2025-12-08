@echo off

REM 启动用户服务RPC服务端
start cmd /k "cd /d h:\go\work\06-blog-cloud\blog_user_service\rpc && go run user.go"

REM 等待服务端启动
ping -n 5 127.0.0.1 > nul

REM 运行测试客户端
echo 启动测试客户端...
c: && cd /d h:\go\work\06-blog-cloud\blog_user_service\rpc\client && go run test_client.go

pause
