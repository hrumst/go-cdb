package replication

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/storage/wal"
)

func TestSlaveBase(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testFsDir := NewMockfsDir(cntrl)
	testTcpClient := NewMocktcpClient(cntrl)
	testAppLogger := NewMockappLogger(cntrl)

	testNextSegment1 := "segment-2"
	testNextSegment2 := "segment-3"
	testNextSegment3 := "segment-4"

	testNextSegmentData1 := []wal.Log{
		{
			UUId:    "uuid-1",
			CmdType: 1,
			Args:    []string{"test-1"},
		}, {
			UUId:    "uuid-2",
			CmdType: 1,
			Args:    []string{"test-2"},
		},
	}
	var testNextSegmentData1Enc = bytes.NewBuffer(make([]byte, 0))
	err := testNextSegmentData1[0].Encode(testNextSegmentData1Enc)
	assert.NoError(t, err)
	err = testNextSegmentData1[1].Encode(testNextSegmentData1Enc)
	assert.NoError(t, err)

	testNextSegmentData2 := []wal.Log{
		{
			UUId:    "uuid-3",
			CmdType: 1,
			Args:    []string{"test-3"},
		},
	}
	var testNextSegmentData2Enc = bytes.NewBuffer(make([]byte, 0))
	err = testNextSegmentData2[0].Encode(testNextSegmentData2Enc)
	assert.NoError(t, err)

	testReq1 := Request{
		PrevSegmentName: testNextSegment1,
	}
	testResp1 := Response{
		Ok:          true,
		SegmentName: testNextSegment2,
		Data:        testNextSegmentData1Enc.Bytes(),
	}
	testReq1Enc, err := respReqEncode(&testReq1)
	assert.NoError(t, err)
	testResp1Enc, err := respReqEncode(&testResp1)
	assert.NoError(t, err)

	testReq2 := Request{
		PrevSegmentName: testNextSegment2,
	}
	testResp2 := Response{
		Ok:          true,
		SegmentName: testNextSegment3,
		Data:        testNextSegmentData2Enc.Bytes(),
	}
	testReq2Enc, err := respReqEncode(&testReq2)
	assert.NoError(t, err)
	testResp2Enc, err := respReqEncode(&testResp2)
	assert.NoError(t, err)

	testReq3 := Request{
		PrevSegmentName: testNextSegment3,
	}
	testResp3 := Response{
		Ok:          false,
		SegmentName: testNextSegment3,
		Data:        nil,
	}
	testReq3Enc, err := respReqEncode(&testReq3)
	assert.NoError(t, err)
	testResp3Enc, err := respReqEncode(&testResp3)
	assert.NoError(t, err)

	testTcpClient.
		EXPECT().
		Request(string(testReq1Enc)).
		Return(string(testResp1Enc), nil).
		Times(1)
	assert.NoError(t, err)

	testTcpClient.
		EXPECT().
		Request(string(testReq2Enc)).
		Return(string(testResp2Enc), nil).
		Times(1)
	assert.NoError(t, err)

	testTcpClient.
		EXPECT().
		Close().
		Return(nil).
		Times(1)
	assert.NoError(t, err)

	testTcpClient.
		EXPECT().
		Request(string(testReq3Enc)).
		Return(string(testResp3Enc), nil).
		Times(1)
	assert.NoError(t, err)

	testFsDir.
		EXPECT().
		WriteSync(testNextSegment2, testNextSegmentData1Enc.Bytes()).
		Return(nil, nil).
		Times(1)

	testFsDir.
		EXPECT().
		WriteSync(testNextSegment3, testNextSegmentData2Enc.Bytes()).
		Return(nil, nil).
		Times(1)

	testCtx, cancel := context.WithTimeout(
		context.Background(),
		1*time.Second,
	)
	defer cancel()

	testSyncInterval := time.Millisecond * 300
	testDataStream := make(chan []wal.Log, 2)

	testSlaveReplica := newSlave(
		testTcpClient,
		testSyncInterval,
		testNextSegment1,
		testFsDir,
		testDataStream,
		testAppLogger,
	)

	testSlaveReplica.start(testCtx)

	assert.Len(t, testDataStream, 2)
	assert.Equal(t, testNextSegmentData1, <-testDataStream)
	assert.Equal(t, testNextSegmentData2, <-testDataStream)
}
