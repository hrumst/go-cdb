package storage

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database/storage/replication"
	walLogs "github.com/hrumst/go-cdb/internal/database/storage/wal"
)

func NewStorageWithRecover(
	ctx context.Context,
	engine storageEngine,
	walCfg config.AppConfigWal,
	dirFs fsDir,
	initReplica initReplica,
	replicaType string,
	logger appLogger,
) (*storage, error) {
	walRepo := walLogs.NewWalRepoFs(dirFs, config.ParserSize(walCfg.MaxSegmentSize))
	walFs := walLogs.InitWal(
		ctx,
		walCfg.FlushBatchSize,
		walCfg.FlushBatchTimeout,
		walRepo,
	)

	newStorage := &storage{
		engine:   engine,
		wal:      walFs,
		isMaster: replicaType == replication.ReplicaTypeMaster,
	}

	if replicaType == replication.ReplicaTypeMaster {
		if err := initReplica.InitMasterReplica(ctx); err != nil {
			return nil, fmt.Errorf("failed to init master replica: %w", err)
		}
	} else if replicaType == replication.ReplicaTypeSlave {
		replicaLogsStream := make(chan []walLogs.Log)
		if err := initReplica.InitSlaveReplica(ctx, replicaLogsStream); err != nil {
			return nil, fmt.Errorf("failed to init slave replica: %w", err)
		}
		go func() {
			for {
				select {
				case <-ctx.Done():
					break
				case logs := <-replicaLogsStream:
					if err := newStorage.applyLogs(ctx, logs); err != nil {
						logger.Error(ctx, "failed to apply replica logs", zap.Error(err))
					}
				}
			}
		}()
	} else {
		return nil, fmt.Errorf("unknowh replication instance type")
	}

	recoverLogsIter, err := walRepo.ReadIterator()
	if err != nil {
		return nil, err
	}
	for {
		logs, err := recoverLogsIter.ReadNext()
		if err != nil {
			if errors.Is(err, walLogs.NextSegmentNotExists) {
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
