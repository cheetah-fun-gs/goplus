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
func RegisterConn(name string, pool redigo.Conn) error {
	if _, ok := mConn[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mConn[name] = pool
	return nil
}

// GetConn 获取连接池
func GetConn() (redigo.Conn, error) {
	return GetConnN(d)
}

// GetConnN get pool with name 获取连接池
func GetConnN(name string) (redigo.Conn, error) {
	if pool, ok := mConn[name]; ok {
		return pool, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}
