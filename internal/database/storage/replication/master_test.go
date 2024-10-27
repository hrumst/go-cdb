package replication

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/tools"
)

func TestMasterReplica(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testFsDir := NewMockfsDir(cntrl)
	testAppLogger := NewMockappLogger(cntrl)

	testFsDir.
		EXPECT().
		FilesStats().
		Return(
			[]os.FileInfo{
				tools.MockFile{100, "segment-0"},
				tools.MockFile{100, "segment-1"},
				tools.MockFile{100, "segment-2"},
			},
			nil,
		).
		Times(2)

	testFsDir.
		EXPECT().
		ReadFile("segment-2").
		Return([]byte("some-data"), nil).
		Times(1)

	testMaster := newMaster(testFsDir, testAppLogger)

	testReq1 := Request{PrevSegmentName: "segment-1"}
	testReqInput1, err := respReqEncode(&testReq1)
	assert.NoError(t, err)
	respOut1 := testMaster.handle(context.Background(), testReqInput1)

	var testResp Response
	err = respReqDecode(&testResp, respOut1)
	assert.NoError(t, err)

	expectResp := Response{
		true,
		"segment-2",
		[]byte("some-data"),
	}
	assert.Equal(t, expectResp, testResp)

	testReq2 := Request{PrevSegmentName: "segment-2"}
	testReqInput2, err := respReqEncode(&testReq2)
	assert.NoError(t, err)
	respOut2 := testMaster.handle(context.Background(), testReqInput2)

	var testResp2 Response
	err = respReqDecode(&testResp2, respOut2)
	assert.NoError(t, err)

	expectResp2 := Response{
		true,
		"",
		nil,
	}
	assert.Equal(t, expectResp2, testResp2)

	testReqInput3 := []byte("ping")
	respOut3 := testMaster.handleWithPingCheck(context.Background(), testReqInput3)
	assert.Equal(t, []byte("pong"), respOut3)
}
