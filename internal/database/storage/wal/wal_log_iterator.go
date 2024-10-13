package wal

import (
	"bytes"
	"errors"
	"os"
)

var NextSegmentNotExists = errors.New("next segment not exists")

type logReadIteratorFs struct {
	segmentDirFs   fsDir
	segmentFiles   []os.FileInfo
	nextSegmentIdx int
}

func NewLogReadIteratorFs(
	segmentDirFs fsDir,
	segmentFiles []os.FileInfo,
) *logReadIteratorFs {
	return &logReadIteratorFs{
		segmentDirFs: segmentDirFs,
		segmentFiles: segmentFiles,
	}
}

func (li *logReadIteratorFs) ReadNext() ([]Log, error) {
	nsi := li.nextSegmentIdx
	if nsi >= len(li.segmentFiles) {
		return nil, NextSegmentNotExists
	}

	data, err := li.segmentDirFs.ReadFile(li.segmentFiles[nsi].Name())
	if err != nil {
		return nil, err
	}

	result := make([]Log, 0)
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var decLog Log
		if err := decLog.Decode(buffer); err != nil {
			return nil, err
		}
		result = append(result, decLog)
	}
	li.nextSegmentIdx += 1
	return result, nil
}
