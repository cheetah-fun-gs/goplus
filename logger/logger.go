package logger

import (
	"context"
	"log"
)

// Logger 日志器
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Debugc(ctx context.Context, format string, v ...interface{})
	Infoc(ctx context.Context, format string, v ...interface{})
	Warnc(ctx context.Context, format string, v ...interface{})
	Errorc(ctx context.Context, format string, v ...interface{})
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

// Debugc 级别日志
func (logger *DefaultLogger) Debugc(ctx context.Context, format string, v ...interface{}) {
	if logger.IsDebugMode {
		log.SetPrefix("[Debug] ")
		logPrintf(ctx, format, v...)
	}
	return
}

// Infoc 级别日志
func (logger *DefaultLogger) Infoc(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Info] ")
	logPrintf(ctx, format, v...)
	return
}

// Warnc 级别日志
func (logger *DefaultLogger) Warnc(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Warn] ")
	logPrintf(ctx, format, v...)
	return
}

// Errorc 级别日志
func (logger *DefaultLogger) Errorc(ctx context.Context, format string, v ...interface{}) {
	log.SetPrefix("[Error] ")
	logPrintf(ctx, format, v...)
	return
}

// Debug 级别日志
func (logger *DefaultLogger) Debug(format string, v ...interface{}) {
	if logger.IsDebugMode {
		log.SetPrefix("[Debug] ")
		log.Printf(format, v...)
	}
	return
}

// Info 级别日志
func (logger *DefaultLogger) Info(format string, v ...interface{}) {
	log.SetPrefix("[Info] ")
	log.Printf(format, v...)
	return
}

// Warn 级别日志
func (logger *DefaultLogger) Warn(format string, v ...interface{}) {
	log.SetPrefix("[Warn] ")
	log.Printf(format, v...)
	return
}

// Error 级别日志
func (logger *DefaultLogger) Error(format string, v ...interface{}) {
	log.SetPrefix("[Error] ")
	log.Printf(format, v...)
	return
}

// New 一个新的日志器
func New() *DefaultLogger {
	return &DefaultLogger{IsDebugMode: true}
}
