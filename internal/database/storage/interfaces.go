package storage

import (
	"context"
	"os"

	"go.uber.org/zap/zapcore"

	"github.com/hrumst/go-cdb/internal/database/compute"
	wal2 "github.com/hrumst/go-cdb/internal/database/storage/wal"
)

//go:generate mockgen -source=interfaces.go -package=storage -destination=mock.go

type storageEngine interface {
	Set(ctx context.Context, key string, val string) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
}

type wal interface {
	AddLogRecord(ctx context.Context, cmdType compute.CommandType, args []string) error
}

type fsDir interface {
	FilesStats() ([]os.FileInfo, error)
	WriteSync(filename string, data []byte) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

type appLogger interface {
	Debug(ctx context.Context, msg string, fields ...zapcore.Field)
	Info(ctx context.Context, msg string, fields ...zapcore.Field)
	Error(ctx context.Context, msg string, fields ...zapcore.Field)
}

type initReplica interface {
	InitSlaveReplica(ctx context.Context, stream chan []wal2.Log) error
	InitMasterReplica(ctx context.Context) error
}
