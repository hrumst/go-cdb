package replication

import (
	"context"
	"os"
	"sort"
	"sync"

	"go.uber.org/zap"
)

type master struct {
	fsDir        fsDir
	logger       appLogger
	nextSegments []os.FileInfo
	mu           sync.Mutex
}

func newMaster(fsDir fsDir, logger appLogger) *master {
	return &master{
		fsDir:  fsDir,
		logger: logger,
	}
}

func (m *master) nextFetchedSegment(fromSegment string) string {
	for fi := 0; fi < len(m.nextSegments); fi += 1 {
		if m.nextSegments[fi].Name() <= fromSegment {
			continue
		}
		result := m.nextSegments[fi]
		return result.Name()
	}
	return ""
}

func (m *master) nextSegment(fromSegment string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := m.nextFetchedSegment(fromSegment)
	if len(result) > 0 {
		return result, nil
	}

	fileStats, err := m.fsDir.FilesStats()
	if err != nil {
		return "", err
	}

	sort.Slice(
		fileStats,
		func(i, j int) bool {
			return fileStats[i].Name() < fileStats[j].Name()
		},
	)
	m.nextSegments = fileStats
	return m.nextFetchedSegment(fromSegment), nil
}

func (m *master) formatResponse(ctx context.Context, resp *Response) []byte {
	encResp, encErr := respReqEncode(&resp)
	if encErr != nil {
		m.logger.Error(ctx, "response encode failed", zap.Error(encErr))
		return nil
	}
	return encResp
}

func (m *master) handleWithPingCheck(ctx context.Context, input []byte) []byte {
	if ctx.Err() != nil {
		return nil
	}
	if len(input) == 4 && string(input) == "ping" {
		return []byte("pong")
	}

	return m.handle(ctx, input)
}

func (m *master) handle(ctx context.Context, input []byte) []byte {
	var (
		req  Request
		resp Response
	)

	if err := respReqDecode(&req, input); err != nil {
		m.logger.Error(ctx, "request decode failed", zap.Error(err))
		return m.formatResponse(ctx, &resp)
	}

	nsm, err := m.nextSegment(req.PrevSegmentName)
	if err != nil {
		m.logger.Error(ctx, "get next segment failed", zap.Error(err))
		return m.formatResponse(ctx, &resp)
	}

	if len(nsm) == 0 {
		resp.Ok = true
		return m.formatResponse(ctx, &resp)
	}

	data, err := m.fsDir.ReadFile(nsm)
	if err != nil {
		resp.Ok = false
		m.logger.Error(ctx, "read data from next segment failed", zap.Error(err))
		return m.formatResponse(ctx, &resp)
	}

	resp.Ok = true
	resp.Data = data
	resp.SegmentName = nsm
	return m.formatResponse(ctx, &resp)
}
