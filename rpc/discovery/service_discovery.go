package discovery

import (
	"errors"
	"sync"
	"time"
)

// ServiceInstance 服务实例
type ServiceInstance struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Address  string            `json:"address"`
	Port     int               `json:"port"`
	Metadata map[string]string `json:"metadata"`
	Tags     []string          `json:"tags"`
}

// ServiceDiscovery 服务发现接口
type ServiceDiscovery interface {
	// Register 注册服务实例
	Register(instance *ServiceInstance) error
	// Deregister 注销服务实例
	Deregister(id string) error
	// GetService 获取服务实例列表
	GetService(name string) ([]*ServiceInstance, error)
	// GetInstances 获取所有服务实例
	GetInstances() []*ServiceInstance
	// Close 关闭服务发现
	Close() error
}

// LocalServiceDiscovery 本地服务发现实现（用于开发和测试）
type LocalServiceDiscovery struct {
	instances map[string]*ServiceInstance
	mutex     sync.RWMutex
}

// NewLocalServiceDiscovery 创建本地服务发现实例
func NewLocalServiceDiscovery() *LocalServiceDiscovery {
	return &LocalServiceDiscovery{
		instances: make(map[string]*ServiceInstance),
	}
}

// Register 注册服务实例
func (d *LocalServiceDiscovery) Register(instance *ServiceInstance) error {
	if instance.ID == "" {
		return errors.New("服务实例 ID 不能为空")
	}
	if instance.Name == "" {
		return errors.New("服务实例名称不能为空")
	}
	if instance.Address == "" {
		return errors.New("服务实例地址不能为空")
	}
	if instance.Port <= 0 {
		return errors.New("服务实例端口必须大于 0")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	// 设置默认元数据
	if instance.Metadata == nil {
		instance.Metadata = make(map[string]string)
	}

	// 设置注册时间
	instance.Metadata["registered_at"] = time.Now().Format(time.RFC3339)

	d.instances[instance.ID] = instance
	return nil
}

// Deregister 注销服务实例
func (d *LocalServiceDiscovery) Deregister(id string) error {
	if id == "" {
		return errors.New("服务实例 ID 不能为空")
	}

	d.mutex.Lock()
	defer d.mutex.Unlock()

	if _, exists := d.instances[id]; !exists {
		return errors.New("服务实例不存在")
	}

	delete(d.instances, id)
	return nil
}

// GetService 获取服务实例列表
func (d *LocalServiceDiscovery) GetService(name string) ([]*ServiceInstance, error) {
	if name == "" {
		return nil, errors.New("服务名称不能为空")
	}

	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var instances []*ServiceInstance
	for _, instance := range d.instances {
		if instance.Name == name {
			instances = append(instances, instance)
		}
	}

	return instances, nil
}

// GetInstances 获取所有服务实例
func (d *LocalServiceDiscovery) GetInstances() []*ServiceInstance {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	var instances []*ServiceInstance
	for _, instance := range d.instances {
		instances = append(instances, instance)
	}

	return instances
}

// Close 关闭服务发现
func (d *LocalServiceDiscovery) Close() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.instances = make(map[string]*ServiceInstance)
	return nil
}
