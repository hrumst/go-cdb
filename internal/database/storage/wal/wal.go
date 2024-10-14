package wal

import (
	"context"
	"sync"
	"time"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/tools"
)

type logBuffer struct {
	resultChan chan error
	logs       []Log
}

func newLogBuffer() *logBuffer {
	return &logBuffer{
		resultChan: make(chan error),
	}
}

type wal struct {
	logBufferMaxSize int64

	logsRepo     logsRepository
	bufLock      *sync.Mutex
	curLogBuffer *logBuffer
}

const (
	defaultLogBufferMaxSize       = 1 << 7
	defaultLogBufferFlushInterval = 10 * time.Millisecond
)

func InitWal(
	ctx context.Context,
	logBufferMaxSize int64,
	logBufferFlushInterval time.Duration,
	logsRepo logsRepository,
) *wal {
	if logBufferMaxSize < 1 {
		logBufferMaxSize = defaultLogBufferMaxSize
	}
	if logBufferFlushInterval < 1*time.Millisecond {
		logBufferFlushInterval = defaultLogBufferFlushInterval
	}

	wl := &wal{
		logBufferMaxSize: logBufferMaxSize,
		logsRepo:         logsRepo,
		bufLock:          &sync.Mutex{},
		curLogBuffer:     newLogBuffer(),
	}

	go func() {
		flushLogTimeTick := time.NewTicker(logBufferFlushInterval)
		defer flushLogTimeTick.Stop()

		for {
			select {
			case <-flushLogTimeTick.C:
				tools.WithLock(
					wl.bufLock,
					func() {
						if len(wl.curLogBuffer.logs) > 0 {
							wl.flushRotateBuffer()
						}
					},
				)
			case <-ctx.Done():
				return
			}
		}
	}()

	return wl
}

func (w *wal) AddLogRecord(ctx context.Context, cmdType compute.CommandType, args []string) error {
	log := newLog(cmdType, args)
	var currBuffer *logBuffer
	tools.WithLock(
		w.bufLock,
		func() {
			currBuffer = w.curLogBuffer
			w.curLogBuffer.logs = append(w.curLogBuffer.logs, log)
			if len(w.curLogBuffer.logs) >= int(w.logBufferMaxSize) {
				w.flushRotateBuffer()
			}
		},
	)
	return <-currBuffer.resultChan
}

func (w *wal) flushRotateBuffer() {
	curBuffer := w.curLogBuffer
	w.curLogBuffer = newLogBuffer()

	go func() {
		logs := curBuffer.logs
		err := w.logsRepo.Write(logs)
		for li := 0; li < len(logs); li += 1 {
			curBuffer.resultChan <- err
		}
	}()
}
