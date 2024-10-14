package wal

import "os"

//go:generate mockgen -source=interfaces.go -package=wal -destination=mock.go

type fsDir interface {
	FilesStats() ([]os.FileInfo, error)
	WriteSync(filename string, data []byte) (os.FileInfo, error)
	ReadFile(filename string) ([]byte, error)
}

type logsRepository interface {
	Write(logs []Log) error
}

type logReadIterator interface {
	ReadNext() ([]Log, error)
}
