package network

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"io"
	"net"
	"time"
)

const (
	defaultMaxMessageSize    = 1 << 12
	defaultIdleTimeoutServer = time.Minute * 2
)

type connHandlerQuery struct {
	maxBodySizeBytes int64
	connIdleTimeout  time.Duration
	logger           logger
	handle           func(ctx context.Context, input string) string
}

func NewConnHandlerQuery(
	handle func(ctx context.Context, input string) string,
	maxBodySizeBytes int64,
	connIdleTimeout time.Duration,
	logger logger,
) *connHandlerQuery {
	if maxBodySizeBytes < 1 {
		maxBodySizeBytes = int64(defaultMaxMessageSize)
	}
	if connIdleTimeout < 1 {
		connIdleTimeout = defaultIdleTimeoutServer
	}

	return &connHandlerQuery{
		handle:           handle,
		maxBodySizeBytes: maxBodySizeBytes,
		connIdleTimeout:  connIdleTimeout,
		logger:           logger,
	}
}

func formatErrorResp(err error) string {
	return "handle request error: " + err.Error()
}

var tooLongRequestError = errors.New("too long request")

func (chq *connHandlerQuery) Handle(ctx context.Context, conn net.Conn) {
	defer func() {
		if v := recover(); v != nil {
			chq.logger.Error(ctx, "panic while connection processing", zap.Any("panic", v))
		}
		if err := conn.Close(); err != nil {
			chq.logger.Error(ctx, "connection close error", zap.Error(err))
		}
	}()

	for {
		if chq.connIdleTimeout != 0 {
			if err := conn.SetDeadline(time.Now().Add(chq.connIdleTimeout)); err != nil {
				chq.logger.Error(ctx, "connection setDeadline error", zap.Error(err))
				return
			}
		}

		readBuffer := make([]byte, chq.maxBodySizeBytes+1)
		count, err := conn.Read(readBuffer)
		if err != nil {
			if err == io.EOF {
				chq.logger.Debug(ctx, "connection closed", zap.Error(err))
				break
			}
			chq.logger.Error(ctx, "connection read error", zap.Error(err))
			return
		}

		if int64(count) == chq.maxBodySizeBytes+1 {
			chq.logger.Debug(ctx, "connection read error", zap.Error(tooLongRequestError))
			if _, err := conn.Write([]byte(formatErrorResp(tooLongRequestError))); err != nil {
				chq.logger.Debug(ctx, "connection write error", zap.Error(err))
				return
			}
			continue
		}

		result := chq.handle(ctx, string(readBuffer[:count]))
		if _, err := conn.Write([]byte(result)); err != nil {
			chq.logger.Debug(ctx, "connection write error", zap.Error(err))
			return
		}
	}
}
