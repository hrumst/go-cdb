package wal

import (
	"bytes"
	"os"
	"sort"

	"github.com/hrumst/go-cdb/internal/tools"
)

type walRepoFs struct {
	segmentDirFs   fsDir
	lastSegment    os.FileInfo
	maxSegmentSize int64
}

const (
	defaultMaxSegmentSize = 10 << 20
)

func NewWalRepoFs(segmentDirFs fsDir, maxSegmentSize int64) *walRepoFs {
	if maxSegmentSize < 1 {
		maxSegmentSize = defaultMaxSegmentSize
	}
	return &walRepoFs{
		segmentDirFs:   segmentDirFs,
		maxSegmentSize: maxSegmentSize,
	}
}

func (w *walRepoFs) nextSegmentName() string {
	return "segment." + tools.GenerateUUIDv7() + ".data"
}

func (w *walRepoFs) Write(logs []Log) error {
	var data bytes.Buffer
	for _, log := range logs {
		if err := log.Encode(&data); err != nil {
			return err
		}
	}

	var segmentFilename string
	if w.lastSegment == nil {
		// create new segment
		segmentFilename = w.nextSegmentName()
	} else {
		segmentFilename = w.lastSegment.Name()
	}

	stats, err := w.segmentDirFs.WriteSync(segmentFilename, data.Bytes())
	if err != nil {
		return err
	}

	if stats.Size() >= w.maxSegmentSize {
		// create new segment in the next time
		w.lastSegment = nil
	} else {
		w.lastSegment = stats
	}
	return nil
}

func (w *walRepoFs) ReadIterator() (logReadIterator, error) {
	segmentsStats, err := w.segmentDirFs.FilesStats()
	if err != nil {
		return nil, err
	}
	sort.Slice(
		segmentsStats,
		func(i, j int) bool {
			return segmentsStats[i].Name() < segmentsStats[j].Name()
		},
	)

	return NewLogReadIteratorFs(w.segmentDirFs, segmentsStats), nil
}
