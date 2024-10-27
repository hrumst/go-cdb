package replication

import (
	"context"
	"os"

	"go.uber.org/zap/zapcore"
)

//go:generate mockgen -source=interfaces.go -package=replication -destination=mock.go

type fsDir interface {
	FilesStats() ([]os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
	WriteSync(filename string, data []byte) (os.FileInfo, error)
}

type tcpClient interface {
	Request(input string) (string, error)
	Close() error
}

type appLogger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, fields ...zapcore.Field)
}
