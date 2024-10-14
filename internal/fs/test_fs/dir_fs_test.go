package test_fs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/fs"
)

func TestDirFs(t *testing.T) {
	testFiledir1 := "./test_dir"
	_ = os.RemoveAll(testFiledir1)
	if err := os.MkdirAll(testFiledir1, 0777); err != nil {
		t.Fatal(err)
	}

	testFilename1 := "test_file_1"
	testFilename2 := "test_file_2"

	testDirFs := fs.NewDir(testFiledir1)

	dstats1, err := testDirFs.FilesStats()
	assert.NoError(t, err)
	assert.Empty(t, dstats1)

	stats, err := testDirFs.WriteSync(testFilename1, []byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, stats.Size(), int64(4))
	assert.Equal(t, stats.Name(), testFilename1)

	stats2, err := testDirFs.WriteSync(testFilename1, []byte("test-2"))
	assert.NoError(t, err)
	assert.Equal(t, stats2.Size(), int64(10))
	assert.Equal(t, stats2.Name(), testFilename1)

	stats3, err := testDirFs.WriteSync(testFilename2, []byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, stats3.Size(), int64(4))
	assert.Equal(t, stats3.Name(), testFilename2)

	data, err := testDirFs.ReadFile(testFilename1)
	assert.NoError(t, err)
	assert.Equal(t, "testtest-2", string(data))

	dstats2, err := testDirFs.FilesStats()
	assert.NoError(t, err)
	assert.Len(t, dstats2, 2)
}
