package logger

import (
	"strings"

	"go.uber.org/zap"
)

type Logger struct {
	lg *zap.SugaredLogger
}

func New(env string) *Logger {
	var lg *zap.Logger

	switch strings.ToLower(strings.TrimSpace(env)) {
	case "prod", "production":
		lg, _ = zap.NewProduction()
	default:
		lg, _ = zap.NewDevelopment()
	}

	return &Logger{lg: lg.Sugar()}
}

func (l *Logger) Sync() error {
	return l.lg.Sync()
}

func (l *Logger) Debug(msg string, fields ...any) {
	l.lg.Debugw(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...any) {
	l.lg.Infow(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...any) {
	l.lg.Warnw(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...any) {
	l.lg.Errorw(msg, fields...)
}
