package wal

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestWal_AddLogRecordByLimit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testWalRepo := NewMocklogsRepository(cntrl)
	logSizes := make([]int, 0)
	testWalRepo.
		EXPECT().
		Write(
			tools.DoMatch(
				func(v []Log) bool {
					logSizes = append(logSizes, len(v))
					return len(v) >= 1
				},
			),
		).
		Return(nil).
		Times(3)

	testWal := InitWal(
		context.Background(),
		2,
		1*time.Millisecond,
		testWalRepo,
	)

	wg := sync.WaitGroup{}
	wg.Add(5)
	for i := 0; i < 5; i += 1 {
		go func() {
			defer wg.Done()

			err := testWal.AddLogRecord(
				context.Background(),
				compute.CommandTypeSet,
				[]string{fmt.Sprintf("arg_%d1", i), fmt.Sprintf("arg_%d2", i)},
			)
			assert.NoError(t, err)
		}()
	}
	wg.Wait()

	assert.Equal(t, logSizes[0], 2)
	assert.Equal(t, logSizes[1], 2)
	assert.Equal(t, logSizes[2], 1)

}

func TestWal_AddLogRecordBytimeout(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testWalRepo := NewMocklogsRepository(cntrl)
	logSizes := make([]int, 0)
	testWalRepo.
		EXPECT().
		Write(
			tools.DoMatch(
				func(v []Log) bool {
					logSizes = append(logSizes, len(v))
					return true
				},
			),
		).
		Return(nil).
		Times(2)

	testWal := InitWal(
		context.Background(),
		10,
		1*time.Millisecond,
		testWalRepo,
	)

	wg := sync.WaitGroup{}
	wg.Add(6)
	for i := 0; i < 6; i += 1 {
		go func(i int) {
			defer wg.Done()

			time.Sleep(time.Duration(i) * 350 * time.Microsecond)
			err := testWal.AddLogRecord(
				context.Background(),
				compute.CommandTypeSet,
				[]string{fmt.Sprintf("arg_%d1", i), fmt.Sprintf("arg_%d2", i)},
			)
			assert.NoError(t, err)
		}(i)
	}
	wg.Wait()

	assert.Equal(t, logSizes[0], 3)
	assert.Equal(t, logSizes[1], 3)
}

func TestWal_AddLogRecordError(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testWalRepo := NewMocklogsRepository(cntrl)
	testWalRepo.
		EXPECT().
		Write(
			tools.DoMatch(
				func(v []Log) bool {
					return len(v) == 1
				},
			),
		).
		Return(fmt.Errorf("some error")).
		Times(1)

	testWal := InitWal(
		context.Background(),
		10,
		1*time.Millisecond,
		testWalRepo,
	)

	err := testWal.AddLogRecord(
		context.Background(),
		compute.CommandTypeSet,
		[]string{"arg_1", "arg_2"},
	)
	assert.Error(t, err)
}
