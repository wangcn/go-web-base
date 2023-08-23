package util

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var redisManager = CreateRedisManager()

const (
	RedisInsight string = "redis.master"
)

// ----------------------------------------
//  获取 Redis 客户端
// ----------------------------------------

func Redis(redisName string) *redis.Client {

	log := Log().With(zap.String("redis_name", redisName))
	redisConfig := new(RedisConfig)

	key := redisName
	if client := redisManager.Get(key); client != nil {
		return client
	}
	if err := Cfg("app").UnmarshalKey(redisName, redisConfig); err != nil {
		log.With(zap.Error(err)).Error("load redis def config err")
		return nil
	}

	// 构建 v8 配置
	options := &redis.Options{
		Network:    "tcp",
		Addr:       fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password:   redisConfig.Password,
		DB:         redisConfig.Database,
		MaxRetries: redisConfig.MaxRetries,
		PoolSize:   redisConfig.PoolSize,
	}
	if redisConfig.DialTimeout > 0 {
		options.DialTimeout = time.Duration(redisConfig.DialTimeout) * time.Millisecond
	}
	if redisConfig.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(redisConfig.ReadTimeout) * time.Millisecond
	}
	if redisConfig.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(redisConfig.WriteTimeout) * time.Millisecond
	}
	client := redis.NewClient(options)

	// 尝试连接
	if err := client.Ping(context.Background()).Err(); err != nil {
		log.With(zap.Error(err)).Error("redis ping err")
		return nil
	}

	redisManager.Set(key, client)

	return client
}

// ----------------------------------------
//  Redis V8 客户端配置
// ----------------------------------------

type RedisConfig struct {
	Host         string `toml:"host" mapstructure:"host"`
	Port         int    `toml:"port" mapstructure:"port"`
	Password     string `toml:"password" mapstructure:"password"`
	Database     int    `toml:"database" mapstructure:"database"`
	MaxRetries   int    `toml:"max_retries" mapstructure:"max_retries"`
	PoolSize     int    `toml:"pool_size" mapstructure:"pool_size"`
	Ping         bool   `toml:"ping" mapstructure:"ping"`
	DialTimeout  int    `toml:"dial_timeout" mapstructure:"dial_timeout"`
	ReadTimeout  int    `toml:"read_timeout" mapstructure:"read_timeout"`
	WriteTimeout int    `toml:"write_timeout" mapstructure:"write_timeout"`
}

// ----------------------------------------
//
//	Redis 客户端管理器
//
// ----------------------------------------
type ClientManager struct {
	clients map[string]*redis.Client
	sync.RWMutex
}

// 创建 Redis 客户端管理器实例
func CreateRedisManager() *ClientManager {
	return &ClientManager{clients: make(map[string]*redis.Client)}
}

// 获取给定名称的 Redis 客户端实例（如果客户端不存在则返回 nil）
func (manager *ClientManager) Get(name string) *redis.Client {
	manager.RLock()
	defer manager.RUnlock()

	if client, exists := manager.clients[name]; exists {
		return client
	}

	return nil
}

// 添加或更新 GORM 客户端实例
func (manager *ClientManager) Set(name string, client *redis.Client) {
	manager.Lock()
	manager.clients[name] = client
	manager.Unlock()
}
