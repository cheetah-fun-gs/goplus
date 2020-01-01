package multilogger

import "context"

// 注意 非N 方法不能引用 N 方法, 会造成 caller depth 不一致

// Debugc Debug
func Debugc(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Debugc(ctx, format, v...)
	}
}

// Infoc Info
func Infoc(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Infoc(ctx, format, v...)
	}
}

// Warnc Warn
func Warnc(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Warnc(ctx, format, v...)
	}
}

// Errorc Error
func Errorc(ctx context.Context, format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Errorc(ctx, format, v...)
	}
}

// DebugcN Debug with name
func DebugcN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Debugc(ctx, format, v...)
	}
}

// InfocN Info with name
func InfocN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Infoc(ctx, format, v...)
	}
}

// WarncN Warn with name
func WarncN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Warnc(ctx, format, v...)
	}
}

// ErrorcN Error with name
func ErrorcN(ctx context.Context, name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Errorc(ctx, format, v...)
	}
}

// Debug Debug
func Debug(format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Debug(format, v...)
	}
}

// Info Info
func Info(format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Info(format, v...)
	}
}

// Warn Warn
func Warn(format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Warn(format, v...)
	}
}

// Error Error
func Error(format string, v ...interface{}) {
	if logger, ok := mutil[d]; ok {
		logger.Error(format, v...)
	}
}

// DebugN Debug with name
func DebugN(name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Debug(format, v...)
	}
}

// InfoN Info with name
func InfoN(name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Info(format, v...)
	}
}

// WarnN Warn with name
func WarnN(name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Warn(format, v...)
	}
}

// ErrorN Error with name
func ErrorN(name, format string, v ...interface{}) {
	if logger, ok := mutil[name]; ok {
		logger.Error(format, v...)
	}
}
