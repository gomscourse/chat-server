package logger

import (
	"log/slog"
)

var globalLogger *slog.Logger

func Init(handler slog.Handler) {
	globalLogger = slog.New(handler)
}

func Debug(msg string, args ...any) {
	globalLogger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	globalLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	globalLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	globalLogger.Error(msg, args...)
}

func WithArgs(args ...any) *slog.Logger {
	return globalLogger.With(args...)
}
