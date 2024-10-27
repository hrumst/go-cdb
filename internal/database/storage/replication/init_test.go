package replication

import (
	"context"
	"io/fs"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database/storage/wal"
	"github.com/hrumst/go-cdb/internal/network"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestSlaveInit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	mockFiles := []fs.FileInfo{
		tools.MockFile{
			Filename: "name-1",
		}, tools.MockFile{
			Filename: "name-3",
		}, tools.MockFile{
			Filename: "name-2",
		},
	}

	testFsDir := NewMockfsDir(cntrl)
	testFsDir.EXPECT().
		FilesStats().
		Return(mockFiles, nil).
		Times(1)

	testCtx, cancel := context.WithTimeout(
		context.Background(),
		600*time.Millisecond,
	)
	defer cancel()

	testAppLogger := NewMockappLogger(cntrl)
	testAppLogger.EXPECT().Debug(gomock.Any(), gomock.Any(), gomock.Any())
	testSyncInter := 500 * time.Millisecond

	testReq1 := Request{
		PrevSegmentName: "name-3",
	}
	testResp1 := Response{
		Ok:          false,
		SegmentName: "name-3",
		Data:        nil,
	}
	testResp1Enc, err := respReqEncode(&testResp1)
	assert.NoError(t, err)

	testAddr := "127.0.0.1:9999"
	var recordedReq Request

	go func() {
		err = network.NewTCPServer(testAddr, testAppLogger).
			Start(
				testCtx,
				[]network.ConnInterceptor{},
				network.NewConnHandlerQuery(
					func(ctx context.Context, input []byte) []byte {
						var req Request
						err = respReqDecode(&req, input)
						assert.NoError(t, err)
						recordedReq = req

						return testResp1Enc
					},
					defaultMaxBodySize,
					defaultIdleTimeout,
					testAppLogger,
				).Handle,
			)
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Millisecond)

	testInitReplica := NewInitReplica(
		config.AppConfigReplica{
			MasterAddress: testAddr,
			SyncInterval:  testSyncInter,
		},
		testFsDir,
		testAppLogger,
	)
	testDataStream := make(chan []wal.Log, 1)
	err = testInitReplica.InitSlaveReplica(testCtx, testDataStream)
	assert.NoError(t, err)

	time.Sleep(1 * time.Second)
	assert.Empty(t, testDataStream)
	assert.Equal(t, testReq1, recordedReq)
}

func TestMasterInit(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testFsDir := NewMockfsDir(cntrl)

	testAppLogger := NewMockappLogger(cntrl)
	testAddr := "127.0.0.1:9998"
	testInitReplica := NewInitReplica(
		config.AppConfigReplica{
			MasterAddress: testAddr,
		},
		testFsDir,
		testAppLogger,
	)

	testCtx, cancel := context.WithTimeout(
		context.Background(),
		600*time.Millisecond,
	)
	defer cancel()

	err := testInitReplica.InitMasterReplica(testCtx)
	assert.NoError(t, err)
}
