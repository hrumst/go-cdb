package storage

import (
	"context"

	"github.com/hrumst/go-cdb/internal/database/compute"
	wal2 "github.com/hrumst/go-cdb/internal/database/storage/wal"
)

type storage struct {
	engine storageEngine
	wal    wal
}

func (s storage) Set(ctx context.Context, key string, val string) error {
	if err := s.wal.AddLogRecord(ctx, compute.CommandTypeSet, []string{key, val}); err != nil {
		return err
	}
	return s.engine.Set(ctx, key, val)
}

func (s storage) Del(ctx context.Context, key string) error {
	if err := s.wal.AddLogRecord(ctx, compute.CommandTypeDel, []string{key}); err != nil {
		return err
	}
	return s.engine.Del(ctx, key)
}

func (s storage) Get(ctx context.Context, key string) (string, error) {
	return s.engine.Get(ctx, key)
}

func (s storage) applyLogs(ctx context.Context, logs []wal2.Log) error {
	for _, log := range logs {
		switch log.CmdType {
		case compute.CommandTypeSet:
			if err := s.engine.Set(ctx, log.Args[0], log.Args[1]); err != nil {
				return err
			}
		case compute.CommandTypeDel:
			if err := s.engine.Del(ctx, log.Args[0]); err != nil {
				return err
			}
		}
	}
	return nil
}
