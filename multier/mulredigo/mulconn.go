package mulredigo

import (
	"fmt"
	"sync"

	redigo "github.com/gomodule/redigo/redis"
)

type mutilConn map[string]redigo.Conn

var (
	onceConn sync.Once
	mConn    mutilConn
)

// InitConn 初始化
func InitConn(defaultConn redigo.Conn) {
	onceConn.Do(func() {
		mConn = mutilConn{
			d: defaultConn,
		}
	})
}

// RegisterConn 注册连接池
func RegisterConn(key string, pool redigo.Conn) error {
	if _, ok := mConn[key]; ok {
		return fmt.Errorf("duplicate key: %v", key)
	}
	mConn[key] = pool
	return nil
}

// GetConn 获取连接池
func GetConn() (redigo.Conn, error) {
	return GetConnK(d)
}

// GetConnK get pool with key 获取连接池
func GetConnK(key string) (redigo.Conn, error) {
	if pool, ok := mConn[key]; ok {
		return pool, nil
	}
	return nil, fmt.Errorf("key not found: %v", key)
}
