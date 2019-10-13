package multilogger

import (
	"fmt"
	"sync"

	"github.com/cheetah-fun-gs/goplus/logger"
)

const (
	d = "default"
)

type mutilLogger map[string]logger.Logger

var (
	once  sync.Once
	mutil mutilLogger
)

// Init 初始化
func Init(defaultLogger logger.Logger) {
	once.Do(func() {
		mutil = mutilLogger{
			d: defaultLogger,
		}
	})
}

// Register 注册日志器
func Register(name string, logger logger.Logger) error {
	if _, ok := mutil[name]; ok {
		return fmt.Errorf("duplicate name: %v", name)
	}
	mutil[name] = logger
	return nil
}

// Retrieve 获取 logger.Logger
func Retrieve() logger.Logger {
	return mutil[d]
}

// RetrieveN 获取 logger.Logger
func RetrieveN(name string) (logger.Logger, error) {
	if c, ok := mutil[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("name not found: %v", name)
}

// MustRetrieveN 获取 logger.Logger
func MustRetrieveN(name string) logger.Logger {
	c, err := RetrieveN(name)
	if err != nil {
		panic((err))
	}
	return c
}
