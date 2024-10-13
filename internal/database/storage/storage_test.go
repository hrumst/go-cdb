package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
)

func TestStorage_InMemoryEngine(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	cntrl := gomock.NewController(t)
	defer cntrl.Finish()
	testWal := NewMockwal(cntrl)

	testWal.EXPECT().
		AddLogRecord(
			gomock.Any(),
			compute.CommandTypeSet,
			[]string{testKey, testVal},
		).
		Return(nil).
		Times(1)

	testWal.EXPECT().
		AddLogRecord(
			gomock.Any(),
			compute.CommandTypeDel,
			[]string{testKey},
		).
		Return(nil).
		Times(2)

	testStorage := storage{
		engine: engine.NewInMemoryEngine(),
		wal:    testWal,
	}
	result, err := testStorage.Get(context.Background(), testKey)
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, engine.NotFoundKey)

	err = testStorage.Del(context.Background(), testKey)
	assert.NoError(t, err, engine.NotFoundKey)

	err = testStorage.Set(context.Background(), testKey, testVal)
	assert.NoError(t, err)

	result, err = testStorage.Get(context.Background(), testKey)
	assert.Equal(t, testVal, result)
	assert.NoError(t, err)

	err = testStorage.Del(context.Background(), testKey)
	assert.NoError(t, err)

	result, err = testStorage.Get(context.Background(), testKey)
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, engine.NotFoundKey)
}

func TestStorage_WalError(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	cntrl := gomock.NewController(t)
	defer cntrl.Finish()
	testWal := NewMockwal(cntrl)

	testWal.EXPECT().
		AddLogRecord(
			gomock.Any(),
			compute.CommandTypeSet,
			[]string{testKey, testVal},
		).
		Return(fmt.Errorf("some error")).
		Times(1)

	testStorage := storage{
		engine: engine.NewInMemoryEngine(),
		wal:    testWal,
	}
	err := testStorage.Set(context.Background(), testKey, testVal)
	assert.Error(t, err)
}
