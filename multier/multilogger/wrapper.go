package multilogger

import "context"

// 注意 非N 方法不能引用 N 方法, 会造成 caller depth 不一致

// Debug Debug
func Debug(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Debug(ctx, format, v...)
	}
}

// Info Info
func Info(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Info(ctx, format, v...)
	}
}

// Warn Warn
func Warn(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Warn(ctx, format, v...)
	}
}

// Error Error
func Error(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Error(ctx, format, v...)
	}
}

// DebugN Debug with name
func DebugN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Debug(ctx, format, v...)
	}
}

// InfoN Info with name
func InfoN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Info(ctx, format, v...)
	}
}

// WarnN Warn with name
func WarnN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Warn(ctx, format, v...)
	}
}

// ErrorN Error with name
func ErrorN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Error(ctx, format, v...)
	}
}
