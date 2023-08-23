package database

import "sync"

// ----------------------------------------
//
//	GORM 客户端管理器
//
// ----------------------------------------
type ClientManager struct {
	clients map[string]*Client
	sync.RWMutex
}

// 创建 GORM 客户端管理器实例
func CreateClientManager() *ClientManager {
	return &ClientManager{clients: make(map[string]*Client)}
}

// 获取给定名称的 GORM 客户端实例（如果客户端不存在则返回 nil）
func (manager *ClientManager) Get(name string) *Client {
	manager.RLock()
	defer manager.RUnlock()

	if client, exists := manager.clients[name]; exists {
		return client
	}

	return nil
}

// 添加或更新 GORM 客户端实例
func (manager *ClientManager) Set(name string, client *Client) {
	manager.Lock()
	manager.clients[name] = client
	manager.Unlock()
}
