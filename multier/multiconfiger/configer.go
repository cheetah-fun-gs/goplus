package multiconfiger

import (
	"fmt"
	"sync"

	"github.com/cheetah-fun-gs/goplus/configer"
)

const (
	d = "default"
)

type mutilConfiger map[string]configer.Configer

var (
	once  sync.Once
	mutil mutilConfiger
)

// Init 初始化
func Init(defaultConfiger configer.Configer) {
	once.Do(func() {
		mutil = mutilConfiger{
			d: defaultConfiger,
		}
	})
}

// Register 注册配置器
func Register(name string, configer configer.Configer) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = configer
	return nil
}

// Retrieve 获取Configer
func Retrieve() configer.Configer {
	return mutil[d]
}

// RetrieveAll 获取所有Configer
func RetrieveAll() map[string]configer.Configer {
	return mutil
}

// RetrieveN 获取Configer
func RetrieveN(name string) (configer.Configer, error) {
	if c, ok := mutil[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// MustRetrieveN 获取Configer
func MustRetrieveN(name string) configer.Configer {
	c, err := RetrieveN(name)
	if err != nil {
		panic((err))
	}
	return c
}
