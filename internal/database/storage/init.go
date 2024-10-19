package storage

import (
	"context"
	"errors"

	"github.com/hrumst/go-cdb/internal/config"
	wal2 "github.com/hrumst/go-cdb/internal/database/storage/wal"
)

func NewStorageWithRecover(
	ctx context.Context,
	engine storageEngine,
	walCfg config.AppConfigWal,
	dirFs fsDir,
) (*storage, error) {
	walRepo := wal2.NewWalRepoFs(dirFs, config.ParserSize(walCfg.MaxSegmentSize))
	walFs := wal2.InitWal(
		ctx,
		walCfg.FlushBatchSize,
		walCfg.FlushBatchTimeout,
		walRepo,
	)

	newStorage := &storage{
		engine: engine,
		wal:    walFs,
	}

	recoverLogsIter, err := walRepo.ReadIterator()
	if err != nil {
		return nil, err
	}
	for {
		logs, err := recoverLogsIter.ReadNext()
		if err != nil {
			if errors.Is(err, wal2.NextSegmentNotExists) {
				break
			}
			return nil, err
		}
		if err := newStorage.applyLogs(ctx, logs); err != nil {
			return nil, err
		}
	}
	return newStorage, nil
}
