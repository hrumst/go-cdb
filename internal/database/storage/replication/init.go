package replication

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.uber.org/zap"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database/storage/wal"
	"github.com/hrumst/go-cdb/internal/network"
)

const (
	ReplicaTypeMaster = "master"
	ReplicaTypeSlave  = "slave"
)

const (
	defaultSyncInterval = 10 * time.Second
	defaultMaxBodySize  = 1 << 20
	defaultIdleTimeout  = 1 * time.Hour
)

type InitReplica struct {
	config config.AppConfigReplica
	fsDir  fsDir
	logger appLogger
}

func NewInitReplica(
	config config.AppConfigReplica,
	fsDir fsDir,
	logger appLogger,
) *InitReplica {
	return &InitReplica{
		config: config,
		fsDir:  fsDir,
		logger: logger,
	}
}

func (ir InitReplica) InitSlaveReplica(ctx context.Context, stream chan []wal.Log) error {
	if len(ir.config.MasterAddress) == 0 {
		return fmt.Errorf("invalid address for master replica")
	}

	client, err := network.NewTcpClient(
		ir.config.MasterAddress,
		defaultIdleTimeout,
		defaultMaxBodySize,
	)
	if err != nil {
		return err
	}

	allSegments, err := ir.fsDir.FilesStats()
	if err != nil {
		return err
	}

	sort.Slice(
		allSegments,
		func(i, j int) bool {
			return allSegments[i].Name() < allSegments[j].Name()
		},
	)

	syncInterval := ir.config.SyncInterval
	if syncInterval < 1 {
		syncInterval = defaultSyncInterval
	}

	lastSegmentName := ""
	if len(allSegments) > 0 {
		lastSegmentName = allSegments[len(allSegments)-1].Name()
	}

	slaveInstance := newSlave(client, syncInterval, lastSegmentName, ir.fsDir, stream, ir.logger)
	go slaveInstance.start(ctx)
	return nil
}

func (ir InitReplica) InitMasterReplica(ctx context.Context) error {
	masterInstance := newMaster(ir.fsDir, ir.logger)
	if len(ir.config.MasterAddress) == 0 {
		return fmt.Errorf("invalid address for master replica")
	}

	tcpServer := network.NewTCPServer(ir.config.MasterAddress, ir.logger)
	go func() {
		if err := tcpServer.Start(
			ctx,
			[]network.ConnInterceptor{},
			network.NewConnHandlerQuery(
				masterInstance.handleWithPingCheck,
				defaultMaxBodySize,
				defaultIdleTimeout,
				ir.logger,
			).Handle,
		); err != nil {
			ir.logger.Error(ctx, "failed to start master replica server", zap.Error(err))
		}
	}()

	pingClient, err := network.NewTcpClient(
		ir.config.MasterAddress,
		defaultIdleTimeout,
		defaultMaxBodySize,
	)
	if err != nil {
		return err
	}

	resp, err := pingClient.Request("ping")
	if err != nil {
		return err
	}
	if resp != "pong" {
		return fmt.Errorf("invalid ping-pong resonse")
	}
	return nil
}
