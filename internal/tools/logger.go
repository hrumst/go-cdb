package tools

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type appLogger struct {
	logger *zap.Logger
}

func (a appLogger) Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	a.logger.Debug(msg, fields...)
}

func (a appLogger) Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	a.logger.Info(msg, fields...)
}

func (a appLogger) Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	a.logger.Info(msg, fields...)
}

func NewAppLogger(logger *zap.Logger) *appLogger {
	return &appLogger{
		logger: logger,
	}
}
