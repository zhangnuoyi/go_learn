package main

import (
	"blog_user_service/rpc/types/user"
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 连接到gRPC服务器 - 假设服务在8080端口
	conn, err := grpc.Dial(":8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接到RPC服务器: %v", err)
	}
	defer conn.Close()

	log.Println("成功连接到RPC服务器")

	// 创建客户端
	client := user.NewUserServiceClient(conn)

	// 测试注册接口
	fmt.Println("=== 测试注册接口 ===")
	registerReq := &user.RegisterDto{
		Name:     "test_user_" + time.Now().Format("150405"), // 确保用户名唯一
		Password: "123456",
		Phone:    "13800138000",
		Email:    "test_" + time.Now().Format("150405") + "@example.com",
	}

	registerResp, err := client.Register(context.Background(), registerReq)
	if err != nil {
		log.Printf("注册失败: %v", err)
	} else {
		fmt.Printf("注册成功, 用户ID: %d\n", registerResp.Id)
	}

	// 测试登录接口
	fmt.Println("\n=== 测试登录接口 ===")
	loginReq := &user.LoginDto{
		Username: registerReq.Name, // 使用刚刚注册的用户名
		Password: "123456",
	}

	loginResp, err := client.Login(context.Background(), loginReq)
	if err != nil {
		log.Printf("Login failed: %v", err)
	} else {
		fmt.Printf("Login success, token: %s\n", loginResp.Token)
	}

	// 测试分页查询接口
	fmt.Println("\n=== 测试分页查询接口 ===")
	listReq := &user.PageInfoDto{
		PageNumber: 1,
		PageSize:   10,
	}

	listResp, err := client.ListUser(context.Background(), listReq)
	if err != nil {
		log.Printf("分页查询失败: %v", err)
	} else {
		fmt.Printf("分页查询成功, 总数: %d, 返回用户数: %d\n", listResp.Total, len(listResp.UserInfoVoList))
		// 打印用户列表
		for i, u := range listResp.UserInfoVoList {
			fmt.Printf("用户 %d: ID=%d, 名称=%s, 邮箱=%s\n", i+1, u.Id, u.Name, u.Email)
		}
	}

	// 测试根据ID查询接口
	fmt.Println("\n=== 测试根据ID查询接口 ===")
	if registerResp != nil {
		getByIdReq := &user.IdDto{
			Id: registerResp.Id,
		}

		getByIdResp, err := client.GetUserById(context.Background(), getByIdReq)
		if err != nil {
			log.Printf("根据ID查询失败: %v", err)
		} else {
			fmt.Printf("根据ID查询成功, 用户名称: %s, 邮箱: %s\n", getByIdResp.Name, getByIdResp.Email)
		}
	}

	// 测试更新用户信息接口
	fmt.Println("\n=== 测试更新用户信息接口 ===")
	if registerResp != nil {
		updateReq := &user.UserInfoVo{
			Id:     registerResp.Id,
			Name:   registerReq.Name + "_updated",
			Email:  registerReq.Email,
			Phone:  "13900139000",
			Gender: 1,
		}

		updateResp, err := client.UpdateUserInfo(context.Background(), updateReq)
		if err != nil {
			log.Printf("更新用户信息失败: %v", err)
		} else {
			fmt.Printf("更新用户信息结果: 成功=%v\n", updateResp.Success)
		}
	}

	log.Println("所有RPC接口测试完成")
}
