package multiredigo

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
func RegisterPool(name string, pool *redigo.Pool) error {
	if _, ok := mPool[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mPool[name] = pool
	return nil
}

// GetPool 获取连接池
func GetPool() (*redigo.Pool, error) {
	return GetPoolN(d)
}

// GetPoolN get pool with name 获取连接池
func GetPoolN(name string) (*redigo.Pool, error) {
	if pool, ok := mPool[name]; ok {
		return pool, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}
