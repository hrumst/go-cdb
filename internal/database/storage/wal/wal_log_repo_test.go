package wal

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestWalLogRepoWriteMultiFiles(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	var prevFname string
	testDirFs := NewMockfsDir(cntrl)
	testDirFs.
		EXPECT().
		WriteSync(
			tools.DoMatch(
				func(fname string) bool {
					if len(fname) != 49 {
						return false
					}
					if len(prevFname) == 0 {
						prevFname = fname
						return true
					}
					return fname != prevFname
				},
			),
			gomock.Any(),
		).
		Return(tools.MockFile{250, "test"}, nil).
		Times(2)

	testWalDirFs := NewWalRepoFs(testDirFs, 250)
	testLogs := []Log{
		newLog(compute.CommandTypeSet, []string{"arg1", "arg2"}),
		newLog(compute.CommandTypeSet, []string{"arg3", "arg4"}),
	}
	err := testWalDirFs.Write(testLogs)
	assert.NoError(t, err)

	testLogs2 := []Log{
		newLog(compute.CommandTypeSet, []string{"arg5", "arg6"}),
		newLog(compute.CommandTypeSet, []string{"arg7", "arg8"}),
	}
	err = testWalDirFs.Write(testLogs2)
	assert.NoError(t, err)
}

func TestWalLogRepoWriteSingleFiles(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	var prevFname string
	testDirFs := NewMockfsDir(cntrl)
	testDirFs.
		EXPECT().
		WriteSync(
			tools.DoMatch(
				func(fname string) bool {
					if len(prevFname) == 0 {
						prevFname = fname
						return true
					}
					return fname == "test"
				},
			),
			gomock.Any(),
		).
		Return(tools.MockFile{250, "test"}, nil).
		Times(2)

	testWalDirFs := NewWalRepoFs(testDirFs, 400)
	testLogs := []Log{
		newLog(compute.CommandTypeSet, []string{"arg1", "arg2"}),
		newLog(compute.CommandTypeSet, []string{"arg3", "arg4"}),
	}
	err := testWalDirFs.Write(testLogs)
	assert.NoError(t, err)

	testLogs2 := []Log{
		newLog(compute.CommandTypeSet, []string{"arg5", "arg6"}),
		newLog(compute.CommandTypeSet, []string{"arg7", "arg8"}),
	}
	err = testWalDirFs.Write(testLogs2)
	assert.NoError(t, err)
}

func TestWalLogRepoReadIterator(t *testing.T) {
	cntrl := gomock.NewController(t)
	defer cntrl.Finish()

	testDirFs := NewMockfsDir(cntrl)
	testFile1 := tools.MockFile{200, "test2.data"}
	testFile2 := tools.MockFile{200, "test1.data"}
	testDirFs.
		EXPECT().
		FilesStats().
		Return(
			[]os.FileInfo{testFile1, testFile2},
			nil,
		).
		Times(1)

	testWalDirFs := NewWalRepoFs(testDirFs, 250)
	ri, err := testWalDirFs.ReadIterator()
	assert.NoError(t, err)

	assert.Equal(t, []os.FileInfo{testFile2, testFile1}, ri.(*logReadIteratorFs).segmentFiles)
	assert.Equal(t, testDirFs, ri.(*logReadIteratorFs).segmentDirFs)
	assert.Equal(t, 0, ri.(*logReadIteratorFs).nextSegmentIdx)
}
