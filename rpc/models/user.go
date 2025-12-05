package models

import (
	"sync"
	"time"
)

// User 用户模型
type User struct {
	ID        uint32    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Age       int32     `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserStore 用户存储接口
type UserStore interface {
	Create(user *User) error
	GetByID(id uint32) (*User, error)
	Update(user *User) error
	Delete(id uint32) error
	List(page, pageSize int) ([]*User, int64, error)
}

// InMemoryUserStore 内存中的用户存储实现
type InMemoryUserStore struct {
	users  map[uint32]*User
	mutex  sync.RWMutex
	nextID uint32
}

// NewInMemoryUserStore 创建内存用户存储实例
func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users:  make(map[uint32]*User),
		nextID: 1,
	}
}

// Create 创建用户
func (s *InMemoryUserStore) Create(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user.ID = s.nextID
	s.nextID++
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

// GetByID 根据ID获取用户
func (s *InMemoryUserStore) GetByID(id uint32) (*User, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return nil, nil // 用户不存在
	}
	return user, nil
}

// Update 更新用户
func (s *InMemoryUserStore) Update(user *User) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, exists := s.users[user.ID]
	if !exists {
		return nil // 用户不存在
	}

	user.UpdatedAt = time.Now()
	s.users[user.ID] = user
	return nil
}

// Delete 删除用户
func (s *InMemoryUserStore) Delete(id uint32) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.users, id)
	return nil
}

// List 获取用户列表
func (s *InMemoryUserStore) List(page, pageSize int) ([]*User, int64, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// 创建用户列表
	userList := make([]*User, 0, len(s.users))
	for _, user := range s.users {
		userList = append(userList, user)
	}

	// 计算分页
	total := int64(len(userList))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(userList) {
		return []*User{}, total, nil
	}

	if end > len(userList) {
		end = len(userList)
	}

	return userList[start:end], total, nil
}
