package multiredigopool

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
	once  sync.Once
	mutil mutilPool
)

// Init 初始化
func Init(defaultPool *redigo.Pool) {
	once.Do(func() {
		mutil = mutilPool{
			d: defaultPool,
		}
	})
}

// Register 注册连接池
func Register(name string, pool *redigo.Pool) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = pool
	return nil
}

// Retrieve 获取 *redigo.Pool
func Retrieve() *redigo.Pool {
	return mutil[d]
}

// RetrieveN 获取 *redigo.Pool
func RetrieveN(name string) (*redigo.Pool, error) {
	if c, ok := mutil[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// MustRetrieveN 获取 *redigo.Pool
func MustRetrieveN(name string) *redigo.Pool {
	c, err := RetrieveN(name)
	if err != nil {
		panic((err))
	}
	return c
}
