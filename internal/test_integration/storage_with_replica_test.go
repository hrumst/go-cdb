package test_integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database/storage"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/database/storage/replication"
	"github.com/hrumst/go-cdb/internal/fs"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestStorageWithReplica(t *testing.T) {
	appLogger := tools.NewAppLogger(zap.NewNop())
	testConnAddress := "127.0.0.1:8855"

	testDirPathMaster := "./test_dir_master"
	err := os.RemoveAll(testDirPathMaster)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDirPathMaster, 0755)
	if err != nil {
		t.Fatal(err)
	}

	testDirPathSlave := "./test_dir_slave"
	err = os.RemoveAll(testDirPathSlave)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDirPathSlave, 0755)
	if err != nil {
		t.Fatal(err)
	}

	dirFsMaster := fs.NewDir(testDirPathMaster)
	masterStorage, err := storage.NewStorageWithRecover(
		context.Background(),
		engine.NewInMemoryEngine(),
		config.AppConfigWal{
			FlushBatchSize:    1 << 12,
			FlushBatchTimeout: 10 * time.Millisecond,
			MaxSegmentSize:    "10kb",
		},
		dirFsMaster,
		replication.NewInitReplica(
			config.AppConfigReplica{
				ReplicaType:   replication.ReplicaTypeMaster,
				MasterAddress: testConnAddress,
			},
			dirFsMaster,
			appLogger,
		),
		replication.ReplicaTypeMaster,
		appLogger,
	)
	assert.NoError(t, err)

	dirFsSlave := fs.NewDir(testDirPathSlave)
	slaveStorage, err := storage.NewStorageWithRecover(
		context.Background(),
		engine.NewInMemoryEngine(),
		config.AppConfigWal{
			FlushBatchSize:    1 << 12,
			FlushBatchTimeout: 10 * time.Millisecond,
			MaxSegmentSize:    "10kb",
		},
		dirFsSlave,
		replication.NewInitReplica(
			config.AppConfigReplica{
				ReplicaType:   replication.ReplicaTypeMaster,
				MasterAddress: testConnAddress,
				SyncInterval:  100 * time.Millisecond,
			},
			dirFsSlave,
			appLogger,
		),
		replication.ReplicaTypeSlave,
		appLogger,
	)
	assert.NoError(t, err)

	const (
		testKey1 = "test_key_1"
		testVal1 = "test_val_1"
	)

	err = masterStorage.Set(context.Background(), testKey1, testVal1)
	assert.NoError(t, err)

	mResult1, err := masterStorage.Get(context.Background(), testKey1)
	assert.NoError(t, err)
	assert.Equal(t, testVal1, mResult1)

	slvResult, err := slaveStorage.Get(context.Background(), testKey1)
	assert.ErrorIs(t, err, engine.NotFoundKey)
	assert.Empty(t, slvResult)

	time.Sleep(150 * time.Millisecond)

	slvResult2, err := slaveStorage.Get(context.Background(), testKey1)
	assert.NoError(t, err)
	assert.Equal(t, testVal1, slvResult2)

	err = masterStorage.Del(context.Background(), testKey1)
	assert.NoError(t, err)

	mResult, err := masterStorage.Get(context.Background(), testKey1)
	assert.ErrorIs(t, err, engine.NotFoundKey)
	assert.Empty(t, mResult)

	slvResult3, err := slaveStorage.Get(context.Background(), testKey1)
	assert.NoError(t, err)
	assert.Equal(t, testVal1, slvResult3)

	time.Sleep(100 * time.Millisecond)

	slvResult4, err := masterStorage.Get(context.Background(), testKey1)
	assert.ErrorIs(t, err, engine.NotFoundKey)
	assert.Empty(t, slvResult4)

	err = slaveStorage.Set(context.Background(), testKey1, testVal1)
	assert.Error(t, err)
}
