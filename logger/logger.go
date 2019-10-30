package logger

import (
	"context"
	"log"
)

// Logger 日志器
type Logger interface {
	Debug(ctx context.Context, format string, v ...interface{})
	Info(ctx context.Context, format string, v ...interface{})
	Warn(ctx context.Context, format string, v ...interface{})
	Error(ctx context.Context, format string, v ...interface{})
}

func logPrintf(ctx context.Context, format string, v ...interface{}) {
	if ctx == nil || ctx == context.Background() {
		log.Printf("- "+format, v...)
	} else {
		log.Printf("%v "+format, append([]interface{}{ctx}, v...)...)
	}
}

// DefaultLogger 默认日志器
type DefaultLogger struct {
	IsDebugMode bool // 是否 debug 模式
}

// Debug 级别日志
func (logger *DefaultLogger) Debug(ctx context.Context, format string, v ...interface{}) {
	if logger.IsDebugMode {
		log.SetPrefix("[Debug] ")
		logPrintf(ctx, format, v...)
	}
	return
}

// Info 级别日志
func (logger *DefaultLogger) Info(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Info] ")
	logPrintf(ctx, format, v...)
	return
}

// Warn 级别日志
func (logger *DefaultLogger) Warn(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Warn] ")
	logPrintf(ctx, format, v...)
	return
}

// Error 级别日志
func (logger *DefaultLogger) Error(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Error] ")
	logPrintf(ctx, format, v...)
	return
}

// New 一个新的日志器
func New() *DefaultLogger {
	return &DefaultLogger{IsDebugMode: true}
}
