package replication

import (
	"bytes"
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/hrumst/go-cdb/internal/database/storage/wal"
)

type slave struct {
	client          tcpClient
	syncInterval    time.Duration
	lastSegmentName string
	logger          appLogger
	stream          chan []wal.Log
	fsDir           fsDir
}

func newSlave(
	client tcpClient,
	syncInterval time.Duration,
	lastSegmentName string,
	fsDir fsDir,
	stream chan []wal.Log,
	logger appLogger,
) *slave {
	return &slave{
		client:          client,
		syncInterval:    syncInterval,
		lastSegmentName: lastSegmentName,
		fsDir:           fsDir,
		logger:          logger,
		stream:          stream,
	}
}

func (s *slave) handleNewSegment(segmentName string, segmentData []byte) error {
	if _, err := s.fsDir.WriteSync(segmentName, segmentData); err != nil {
		return err
	}
	s.lastSegmentName = segmentName

	var logs []wal.Log
	buffer := bytes.NewBuffer(segmentData)
	for buffer.Len() > 0 {
		var decLog wal.Log
		if err := decLog.Decode(buffer); err != nil {
			return err
		}
		logs = append(logs, decLog)
	}
	s.stream <- logs
	return nil
}

func (s *slave) start(ctx context.Context) {
	ticker := time.NewTicker(s.syncInterval)
	defer func() {
		ticker.Stop()
		_ = s.client.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			var request Request
			request.PrevSegmentName = s.lastSegmentName
			requestEnc, err := respReqEncode(&request)
			if err != nil {
				s.logger.Error(ctx, "failed encode request", zap.Error(err))
				continue
			}

			respData, err := s.client.Request(string(requestEnc))
			if err != nil {
				s.logger.Error(ctx, "failed receive client response", zap.Error(err))
				continue
			}

			var response Response
			if err := respReqDecode(&response, []byte(respData)); err != nil {
				s.logger.Error(ctx, "failed decode response", zap.Error(err))
				continue
			}

			if !response.Ok || len(response.SegmentName) == 0 {
				continue
			}
			if err := s.handleNewSegment(response.SegmentName, response.Data); err != nil {
				s.logger.Error(ctx, "failed to handle new segment", zap.Error(err))
			}
		}
	}
}
