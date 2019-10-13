package multilogger

import "context"

// Debug Debug
func Debug(ctx context.Context, format string, v ...interface{}) {
	DebugN(ctx, d, format, v...)
}

// Info Info
func Info(ctx context.Context, format string, v ...interface{}) {
	InfoN(ctx, d, format, v...)
}

// Warn Warn
func Warn(ctx context.Context, format string, v ...interface{}) {
	WarnN(ctx, d, format, v...)
}

// Error Error
func Error(ctx context.Context, format string, v ...interface{}) {
	ErrorN(ctx, d, format, v...)
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
