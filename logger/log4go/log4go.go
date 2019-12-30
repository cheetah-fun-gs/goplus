package log4go

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/alecthomas/log4go"
	filepathplus "github.com/cheetah-fun-gs/goplus/path/filepath"
)

// Config 日志器配置
type Config struct {
	Path             string `json:"path,omitempty"`
	Format           string `json:"format,omitempty"`
	IsDebugMode      bool   `json:"is_debug_mode,omitempty"`      // 是否打开debug模式, debug模式才输出 Debug 信息
	IsDisableConsole bool   `json:"is_disable_console,omitempty"` // 是否关闭控制台输出
	CallerDepth      int    `json:"caller_depth,omitempty"`       // 默认为4 -1为关闭 不直接引用, 每封装一次+1
}

// Logger 日志器
type Logger struct {
	log4go.Logger
	c *Config
}

// New ...
func New(name string, c *Config) *Logger {
	logger := make(log4go.Logger)
	fileFormat := "[%T %D] [%L] %M"
	if c.Format != "" {
		fileFormat = c.Format
	}
	consoleFormat := name + " " + fileFormat // 控制台输出多打印日志器名称

	if c.Path != "" {
		logDir := filepath.Dir(c.Path)
		ok, err := filepathplus.Exists(logDir)
		if err != nil {
			panic(err)
		}
		if !ok {
			if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
				panic(err)
			}
		}

		fileOut := log4go.NewFileLogWriter(c.Path, true)
		fileOut.SetFormat(fileFormat)
		fileOut.SetRotateDaily(true)
		logger.AddFilter("file", log4go.FINEST, fileOut)
	}

	if !c.IsDisableConsole {
		consoleOut := log4go.NewConsoleLogWriter()
		consoleOut.SetFormat(consoleFormat)
		logger.AddFilter("stdout", log4go.FINEST, consoleOut)
	}

	return &Logger{Logger: logger, c: c}
}

func callerSource(depth int) string {
	// Determine caller func
	pc, _, lineno, ok := runtime.Caller(depth)
	src := ""
	if ok {
		src = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), lineno)
	}
	return src
}

func (logger *Logger) to(ctx context.Context, format string, v ...interface{}) (string, []interface{}) {
	var headFormat string
	if logger.c.CallerDepth == 0 {
		headFormat = fmt.Sprintf("(%s)", callerSource(4))
	} else if logger.c.CallerDepth > 0 {
		headFormat = fmt.Sprintf("(%s)", callerSource(logger.c.CallerDepth))
	}

	if ctx == context.Background() || ctx == nil {
		headFormat += " - "
	} else {
		headFormat += " %v "
		v = append([]interface{}{ctx}, v...)
	}
	return headFormat + format, v
}

// Debug Debug
func (logger *Logger) Debug(ctx context.Context, format string, v ...interface{}) {
	if logger.c.IsDebugMode {
		format, v = logger.to(ctx, format, v...)
		logger.Logger.Debug(format, v...)
	}
}

// Info Info
func (logger *Logger) Info(ctx context.Context, format string, v ...interface{}) {
	format, v = logger.to(ctx, format, v...)
	logger.Logger.Info(format, v...)
}

// Warn Warn
func (logger *Logger) Warn(ctx context.Context, format string, v ...interface{}) {
	format, v = logger.to(ctx, format, v...)
	logger.Logger.Warn(format, v...)
}

// Error Error
func (logger *Logger) Error(ctx context.Context, format string, v ...interface{}) {
	format, v = logger.to(ctx, format, v...)
	logger.Logger.Error(format, v...)
}
