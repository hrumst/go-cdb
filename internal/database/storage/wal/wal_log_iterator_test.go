package wal

import (
	"bytes"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestWalLogIterator(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testDirFs := NewMockfsDir(cntrl)

	testLog1 := newLog(compute.CommandTypeSet, []string{"arg1", "arg2"})
	testLog2 := newLog(compute.CommandTypeGet, []string{"arg3", "arg4"})
	testBuffer := bytes.NewBuffer([]byte{})
	err := testLog1.Encode(testBuffer)
	assert.NoError(t, err)
	err = testLog2.Encode(testBuffer)
	assert.NoError(t, err)

	testDirFs.
		EXPECT().
		ReadFile("test1.data").
		Return(testBuffer.Bytes(), nil).
		Times(1)

	testLog3 := newLog(compute.CommandTypeSet, []string{"arg5", "arg6"})
	testLog4 := newLog(compute.CommandTypeGet, []string{"arg7", "arg8"})
	testBuffer2 := bytes.NewBuffer([]byte{})
	err = testLog3.Encode(testBuffer2)
	assert.NoError(t, err)
	err = testLog4.Encode(testBuffer2)
	assert.NoError(t, err)

	testDirFs.
		EXPECT().
		ReadFile("test2.data").
		Return(testBuffer2.Bytes(), nil).
		Times(1)

	testReadIterator := NewLogReadIteratorFs(
		testDirFs,
		[]os.FileInfo{
			tools.MockFile{200, "test1.data"},
			tools.MockFile{200, "test2.data"},
		},
	)
	res1, err := testReadIterator.ReadNext()
	assert.NoError(t, err)
	assert.Equal(t, []Log{testLog1, testLog2}, res1)

	res2, err := testReadIterator.ReadNext()
	assert.NoError(t, err)
	assert.Equal(t, []Log{testLog3, testLog4}, res2)

	res3, err := testReadIterator.ReadNext()
	assert.Errorf(t, err, "next segment not exists")
	assert.Empty(t, res3)
}

func TestWalLogIteratorEmpty(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testDirFs := NewMockfsDir(cntrl)

	testReadIterator := NewLogReadIteratorFs(
		testDirFs,
		[]os.FileInfo{},
	)

	res1, err := testReadIterator.ReadNext()
	assert.Errorf(t, err, "next segment not exists")
	assert.Empty(t, res1)
}
