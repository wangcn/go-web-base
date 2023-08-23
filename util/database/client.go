package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Client ----------------------------------------
//
//	GORM 客户端
//
// ----------------------------------------
type Client struct {
	name        string
	config      *Config
	orm         *gorm.DB
	connectedAt time.Time
	pingedAt    time.Time
	sync.Mutex
}

// 创建 GORM 客户端实例
func CreateClient(name string, config *Config) (*Client, error) {
	client := &Client{name: name, config: config}
	if err := client.Connect(); err != nil {
		return nil, err
	} else {
		return client, nil
	}
}

// 连接远程服务
func (client *Client) Connect() error {
	client.Lock()
	defer client.Unlock()

	if client.orm != nil {
		return client.Ping()
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		client.config.Username,
		client.config.Password,
		client.config.Host,
		client.config.Port,
		client.config.Database,
	)
	logLevel := logger.Silent
	if client.config.Debug {
		logLevel = logger.Info
	}

	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return errors.Wrap(err, "gorm open err")
	}

	db, err := orm.DB()
	if err != nil {
		return errors.Wrap(err, "gorm db err")
	}
	db.SetMaxOpenConns(client.config.MaxConn)
	db.SetMaxIdleConns(client.config.MaxIdleConn)

	client.orm = orm
	client.connectedAt = time.Now().Local()
	return client.Ping()
}

// 检查远程服务的可用性
func (client *Client) Ping() error {
	if client.config.Ping && client.orm != nil {
		client.pingedAt = time.Now().Local()
		db, err := client.orm.DB()
		if err != nil {
			return err
		}
		return db.Ping()
	}
	return nil
}

// 获取当前 GORM 客户端的第三方客户端实例
func (client *Client) GetORM() *gorm.DB {
	return client.orm
}
