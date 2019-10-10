package mulredigo

import (
	"fmt"
	"sync"

	redigo "github.com/gomodule/redigo/redis"
)

const (
	d = "default"
)

type mutilPool map[string]*redigo.Pool

var (
	oncePool sync.Once
	mPool    mutilPool
)

// InitPool 初始化
func InitPool(defaultPool *redigo.Pool) {
	oncePool.Do(func() {
		mPool = mutilPool{
			d: defaultPool,
		}
	})
}

// RegisterPool 注册连接池
func RegisterPool(key string, pool *redigo.Pool) error {
	if _, ok := mPool[key]; ok {
		return fmt.Errorf("duplicate key: %v", key)
	}
	mPool[key] = pool
	return nil
}

// GetPool 获取连接池
func GetPool() (*redigo.Pool, error) {
	return GetPoolK(d)
}

// GetPoolK get pool with key 获取连接池
func GetPoolK(key string) (*redigo.Pool, error) {
	if pool, ok := mPool[key]; ok {
		return pool, nil
	}
	return nil, fmt.Errorf("key not found: %v", key)
}
