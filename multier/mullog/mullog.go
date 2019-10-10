package mullog

import (
	"context"
	"fmt"
	"sync"
)

const (
	d = "default"
)

// Logger 日志器
type Logger interface {
	Debug(ctx context.Context, format string, v ...interface{})
	Info(ctx context.Context, format string, v ...interface{})
	Warn(ctx context.Context, format string, v ...interface{})
	Error(ctx context.Context, format string, v ...interface{})
}

type mutilLogger map[string]Logger

var (
	once  sync.Once
	mutil mutilLogger
)

// Init 初始化
func Init(defaultLogger Logger) {
	once.Do(func() {
		mutil = mutilLogger{
			d: defaultLogger,
		}
	})
}

// Register 注册日志器
func Register(key string, logger Logger) error {
	if _, ok := mutil[key]; ok {
		return fmt.Errorf("duplicate key: %v", key)
	}
	mutil[key] = logger
	return nil
}

// Debug Debug
func Debug(ctx context.Context, format string, v ...interface{}) {
	DebugK(ctx, d, format, v...)
}

// Info Info
func Info(ctx context.Context, format string, v ...interface{}) {
	InfoK(ctx, d, format, v...)
}

// Warn Warn
func Warn(ctx context.Context, format string, v ...interface{}) {
	WarnK(ctx, d, format, v...)
}

// Error Error
func Error(ctx context.Context, format string, v ...interface{}) {
	ErrorK(ctx, d, format, v...)
}

// DebugK Debug with key
func DebugK(ctx context.Context, key, format string, v ...interface{}) {
	if logger, ok := mutil[key]; ok {
		logger.Debug(ctx, format, v...)
	}
}

// InfoK Info with key
func InfoK(ctx context.Context, key, format string, v ...interface{}) {
	if logger, ok := mutil[key]; ok {
		logger.Info(ctx, format, v...)
	}
}

// WarnK Warn with key
func WarnK(ctx context.Context, key, format string, v ...interface{}) {
	if logger, ok := mutil[key]; ok {
		logger.Warn(ctx, format, v...)
	}
}

// ErrorK Error with key
func ErrorK(ctx context.Context, key, format string, v ...interface{}) {
	if logger, ok := mutil[key]; ok {
		logger.Error(ctx, format, v...)
	}
}
