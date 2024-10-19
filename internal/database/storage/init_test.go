package storage

import (
	"bytes"
	"context"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	wal2 "github.com/hrumst/go-cdb/internal/database/storage/wal"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestInitWithRecover(t *testing.T) {
	var (
		testKey1 = "key1"
		testKey2 = "key2"
		testVal1 = "val1"
		testVal2 = "val2"
	)

	buffer1 := bytes.NewBuffer(make([]byte, 0))
	testCmd1 := wal2.Log{
		UUId:    tools.GenerateUUIDv7(),
		CmdType: compute.CommandTypeSet,
		Args:    []string{testKey1, testVal1},
	}
	err := testCmd1.Encode(buffer1)
	assert.NoError(t, err)

	testCmd2 := wal2.Log{
		UUId:    tools.GenerateUUIDv7(),
		CmdType: compute.CommandTypeSet,
		Args:    []string{testKey2, testVal2},
	}
	err = testCmd2.Encode(buffer1)
	assert.NoError(t, err)

	testCmd3 := wal2.Log{
		UUId:    tools.GenerateUUIDv7(),
		CmdType: compute.CommandTypeDel,
		Args:    []string{testKey1},
	}
	err = testCmd3.Encode(buffer1)
	assert.NoError(t, err)

	buffer2 := bytes.NewBuffer(make([]byte, 0))
	testCmd4 := wal2.Log{
		UUId:    tools.GenerateUUIDv7(),
		CmdType: compute.CommandTypeSet,
		Args:    []string{testKey2, testVal1},
	}
	err = testCmd4.Encode(buffer2)
	assert.NoError(t, err)

	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testFilename1, testFilename2 := "filename1.data", "filename2.data"
	testDirFs := NewMockfsDir(cntrl)
	testDirFs.
		EXPECT().
		FilesStats().
		Return(
			[]os.FileInfo{
				tools.MockFile{200, testFilename1},
				tools.MockFile{200, testFilename2},
			},
			nil,
		).
		Times(1)

	testDirFs.
		EXPECT().
		ReadFile(testFilename1).
		Return(buffer1.Bytes(), nil).
		Times(1)
	testDirFs.
		EXPECT().
		ReadFile(testFilename2).
		Return(buffer2.Bytes(), nil).
		Times(1)

	initStorage, err := NewStorageWithRecover(
		context.Background(),
		engine.NewInMemoryEngine(),
		config.AppConfigWal{
			FlushBatchSize:    100,
			FlushBatchTimeout: 10 * time.Millisecond,
			MaxSegmentSize:    "30kb",
		},
		testDirFs,
	)
	assert.NoError(t, err)

	res, err := initStorage.Get(context.Background(), testKey2)
	assert.NoError(t, err)
	assert.Equal(t, testVal1, res)

	res, err = initStorage.Get(context.Background(), testKey1)
	assert.ErrorIs(t, err, engine.NotFoundKey)
	assert.Empty(t, res)
}

func TestInitWithEmptyLogs(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testDirFs := NewMockfsDir(cntrl)
	testDirFs.
		EXPECT().
		FilesStats().
		Return(
			[]os.FileInfo{},
			nil,
		).
		Times(1)

	initStorage, err := NewStorageWithRecover(
		context.Background(),
		engine.NewInMemoryEngine(),
		config.AppConfigWal{
			FlushBatchSize:    100,
			FlushBatchTimeout: 10 * time.Millisecond,
			MaxSegmentSize:    "30kb",
		},
		testDirFs,
	)
	assert.NoError(t, err)
	res, err := initStorage.Get(context.Background(), "key")
	assert.ErrorIs(t, err, engine.NotFoundKey)
	assert.Empty(t, res)
}
