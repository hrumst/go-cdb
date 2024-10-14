package tools

import (
	"io/fs"
	"time"
)

type MockFile struct {
	Filesize int64
	Filename string
}

func (m MockFile) Name() string {
	return m.Filename
}

func (m MockFile) Size() int64 {
	return m.Filesize
}

func (m MockFile) Mode() fs.FileMode {
	return 0777
}

func (m MockFile) ModTime() time.Time {
	return time.Now()
}

func (m MockFile) IsDir() bool {
	return false
}

func (m MockFile) Sys() any {
	return nil
}
