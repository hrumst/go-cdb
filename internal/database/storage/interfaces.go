package storage

import (
	"context"
	"os"

	"github.com/hrumst/go-cdb/internal/database/compute"
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
