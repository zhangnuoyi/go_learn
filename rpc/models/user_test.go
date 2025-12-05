package models

import (
	"testing"
)

// TestInMemoryUserStore_Create 测试创建用户
func TestInMemoryUserStore_Create(t *testing.T) {
	store := NewInMemoryUserStore()

	user := &User{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
	}

	// 创建用户
	err := store.Create(user)
	if err != nil {
		t.Errorf("创建用户失败: %v", err)
	}

	// 验证用户已创建
	if user.ID == 0 {
		t.Error("用户ID应该被设置")
	}

	if user.Name != "测试用户" {
		t.Errorf("用户名称不匹配: 期望 %s, 实际 %s", "测试用户", user.Name)
	}

	if user.Email != "test@example.com" {
		t.Errorf("用户邮箱不匹配: 期望 %s, 实际 %s", "test@example.com", user.Email)
	}

	if user.Age != 25 {
		t.Errorf("用户年龄不匹配: 期望 %d, 实际 %d", 25, user.Age)
	}
}

// TestInMemoryUserStore_GetByID 测试根据ID获取用户
func TestInMemoryUserStore_GetByID(t *testing.T) {
	store := NewInMemoryUserStore()

	// 先创建一个用户
	user := &User{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
	}
	_ = store.Create(user)

	// 根据ID获取用户
	retrievedUser, err := store.GetByID(user.ID)
	if err != nil {
		t.Errorf("根据ID获取用户失败: %v", err)
	}

	// 验证用户信息
	if retrievedUser.ID != user.ID {
		t.Errorf("用户ID不匹配: 期望 %d, 实际 %d", user.ID, retrievedUser.ID)
	}

	if retrievedUser.Name != "测试用户" {
		t.Errorf("用户名称不匹配: 期望 %s, 实际 %s", "测试用户", retrievedUser.Name)
	}

	// 测试获取不存在的用户
	retrievedUser, err = store.GetByID(999)
	if retrievedUser != nil {
		t.Error("获取不存在的用户应该返回nil")
	}
}

// TestInMemoryUserStore_Update 测试更新用户
func TestInMemoryUserStore_Update(t *testing.T) {
	store := NewInMemoryUserStore()

	// 先创建一个用户
	user := &User{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
	}
	_ = store.Create(user)

	// 更新用户信息
	user.Name = "更新后的用户"
	user.Email = "updated@example.com"
	user.Age = 30

	err := store.Update(user)
	if err != nil {
		t.Errorf("更新用户失败: %v", err)
	}

	// 重新获取用户，验证信息已更新
	updatedUser, _ := store.GetByID(user.ID)
	if updatedUser.Name != "更新后的用户" {
		t.Errorf("用户名称不匹配: 期望 %s, 实际 %s", "更新后的用户", updatedUser.Name)
	}

	if updatedUser.Email != "updated@example.com" {
		t.Errorf("用户邮箱不匹配: 期望 %s, 实际 %s", "updated@example.com", updatedUser.Email)
	}

	if updatedUser.Age != 30 {
		t.Errorf("用户年龄不匹配: 期望 %d, 实际 %d", 30, updatedUser.Age)
	}

	// 测试更新不存在的用户
	nonExistentUser := &User{
		ID:    999,
		Name:  "不存在的用户",
		Email: "nonexistent@example.com",
		Age:   25,
	}

	err = store.Update(nonExistentUser)
	// 根据实现，更新不存在的用户不会返回错误，只是不做任何操作
	if err != nil {
		t.Errorf("更新不存在的用户不应该返回错误: %v", err)
	}
}

// TestInMemoryUserStore_Delete 测试删除用户
func TestInMemoryUserStore_Delete(t *testing.T) {
	store := NewInMemoryUserStore()

	// 先创建一个用户
	user := &User{
		Name:  "测试用户",
		Email: "test@example.com",
		Age:   25,
	}
	_ = store.Create(user)

	// 删除用户
	err := store.Delete(user.ID)
	if err != nil {
		t.Errorf("删除用户失败: %v", err)
	}

	// 验证用户已删除
	retrievedUser, _ := store.GetByID(user.ID)
	if retrievedUser != nil {
		t.Error("获取已删除的用户应该返回nil")
	}

	// 测试删除不存在的用户
	err = store.Delete(999)
	// 根据实现，删除不存在的用户不会返回错误
	if err != nil {
		t.Errorf("删除不存在的用户不应该返回错误: %v", err)
	}
}

// TestInMemoryUserStore_List 测试获取用户列表
func TestInMemoryUserStore_List(t *testing.T) {
	store := NewInMemoryUserStore()

	// 创建多个用户
	users := []*User{
		{Name: "用户1", Email: "user1@example.com", Age: 20},
		{Name: "用户2", Email: "user2@example.com", Age: 25},
		{Name: "用户3", Email: "user3@example.com", Age: 30},
		{Name: "用户4", Email: "user4@example.com", Age: 35},
		{Name: "用户5", Email: "user5@example.com", Age: 40},
	}

	for _, user := range users {
		_ = store.Create(user)
	}

	// 获取用户列表
	userList, total, err := store.List(1, 10)
	if err != nil {
		t.Errorf("获取用户列表失败: %v", err)
	}

	// 验证用户列表
	if total != 5 {
		t.Errorf("用户总数不匹配: 期望 %d, 实际 %d", 5, total)
	}

	if len(userList) != 5 {
		t.Errorf("用户列表长度不匹配: 期望 %d, 实际 %d", 5, len(userList))
	}

	// 测试分页
	userList, total, err = store.List(1, 3)
	if err != nil {
		t.Errorf("获取用户列表失败: %v", err)
	}

	if total != 5 {
		t.Errorf("用户总数不匹配: 期望 %d, 实际 %d", 5, total)
	}

	if len(userList) != 3 {
		t.Errorf("用户列表长度不匹配: 期望 %d, 实际 %d", 3, len(userList))
	}

	// 测试第二页
	userList, total, err = store.List(2, 3)
	if err != nil {
		t.Errorf("获取用户列表失败: %v", err)
	}

	if total != 5 {
		t.Errorf("用户总数不匹配: 期望 %d, 实际 %d", 5, total)
	}

	if len(userList) != 2 {
		t.Errorf("用户列表长度不匹配: 期望 %d, 实际 %d", 2, len(userList))
	}
}
