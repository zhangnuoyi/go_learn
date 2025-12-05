package main

import (
	"context"
	"fmt"
	"log"

	"github.com/user/rpc/discovery"
	"github.com/user/rpc/proto"
)

func main() {
	// 创建本地服务工厂（连接到本地 gRPC 服务） 进阶
	factory, err := discovery.NewLocalServiceFactory("localhost:50051")
	if err != nil {
		log.Fatalf("创建服务工厂失败: %v", err)
	}
	defer factory.Close()

	// 获取用户服务客户端
	userClient, err := factory.GetUserServiceClient()
	if err != nil {
		log.Fatalf("获取用户服务客户端失败: %v", err)
	}

	ctx := context.Background()

	// 1. 创建用户示例
	fmt.Println("\n1. 创建用户示例:")
	createResp, err := userClient.CreateUser(ctx, &proto.CreateUserRequest{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   30,
	})
	if err != nil {
		log.Printf("创建用户失败: %v", err)
	} else {
		fmt.Printf("创建用户成功，用户ID: %d\n", createResp.User.Id)
	}

	// 2. 获取用户示例
	fmt.Println("\n2. 获取用户示例:")
	getResp, err := userClient.GetUser(ctx, &proto.GetUserRequest{
		Id: createResp.User.Id,
	})
	if err != nil {
		log.Printf("获取用户失败: %v", err)
	} else {
		fmt.Printf("获取用户成功: ID=%d, Name=%s, Email=%s, Age=%d\n",
			getResp.User.Id, getResp.User.Name, getResp.User.Email, getResp.User.Age)
	}

	// 3. 更新用户示例
	fmt.Println("\n3. 更新用户示例:")
	updateResp, err := userClient.UpdateUser(ctx, &proto.UpdateUserRequest{
		Id:    createResp.User.Id,
		Name:  "张三（更新）",
		Email: "zhangsan_updated@example.com",
		Age:   31,
	})
	if err != nil {
		log.Printf("更新用户失败: %v", err)
	} else {
		fmt.Printf("更新用户成功: ID=%d, Name=%s, Email=%s, Age=%d\n",
			updateResp.User.Id, updateResp.User.Name, updateResp.User.Email, updateResp.User.Age)
	}

	// 4. 创建更多用户用于列表展示
	fmt.Println("\n4. 创建更多用户用于列表展示:")
	for i := 1; i <= 3; i++ {
		userName := fmt.Sprintf("用户%d", i)
		userEmail := fmt.Sprintf("user%d@example.com", i)
		_, err := userClient.CreateUser(ctx, &proto.CreateUserRequest{
			Name:  userName,
			Email: userEmail,
			Age:   20 + i,
		})
		if err != nil {
			log.Printf("创建用户 %s 失败: %v", userName, err)
		} else {
			fmt.Printf("创建用户 %s 成功\n", userName)
		}
	}

	// 5. 列出用户示例
	fmt.Println("\n5. 列出用户示例:")
	listResp, err := userClient.ListUsers(ctx, &proto.ListUsersRequest{
		Page:     1,
		PageSize: 10,
	})
	if err != nil {
		log.Printf("列出用户失败: %v", err)
	} else {
		fmt.Printf("总用户数: %d\n", listResp.Total)
		fmt.Println("用户列表:")
		for _, user := range listResp.Users {
			fmt.Printf("  ID=%d, Name=%s, Email=%s, Age=%d\n",
				user.Id, user.Name, user.Email, user.Age)
		}
	}

	// 6. 删除用户示例
	fmt.Println("\n6. 删除用户示例:")
	_, err = userClient.DeleteUser(ctx, &proto.DeleteUserRequest{
		Id: createResp.User.Id,
	})
	if err != nil {
		log.Printf("删除用户失败: %v", err)
	} else {
		fmt.Printf("删除用户成功，用户ID: %d\n", createResp.User.Id)
	}

	// 7. 验证用户已删除
	fmt.Println("\n7. 验证用户已删除:")
	_, err = userClient.GetUser(ctx, &proto.GetUserRequest{
		Id: createResp.User.Id,
	})
	if err != nil {
		fmt.Printf("预期的错误（用户已删除）: %v\n", err)
	} else {
		fmt.Printf("错误: 用户应该已被删除，但仍然存在\n")
	}

	fmt.Println("\n客户端示例运行完成！")
}
