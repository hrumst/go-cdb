package network

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net"
	"time"
)

const (
	defaultMaxConns      = 1 << 7
	defaultAcceptTimeout = time.Second * 5
)

var acceptConnDeadlineExceededError = errors.New("accept connection deadline exceeded")

type connLimiterWithTimeout struct {
	connLimit         chan struct{}
	connAcceptTimeout time.Duration
	waitCheckTimeout  time.Duration
	logger            logger
}

func NewConnLimiterWithTimeout(
	maxConns int64,
	connAcceptTimeout time.Duration,
	logger logger,
) *connLimiterWithTimeout {
	if maxConns < 1 {
		maxConns = int64(defaultMaxConns)
	}
	if connAcceptTimeout < 1 {
		connAcceptTimeout = defaultAcceptTimeout
	}

	waitCheckTimeout := connAcceptTimeout / 10
	if waitCheckTimeout < time.Millisecond*10 {
		waitCheckTimeout = time.Millisecond * 10
	}

	return &connLimiterWithTimeout{
		connLimit:         make(chan struct{}, maxConns),
		connAcceptTimeout: connAcceptTimeout,
		waitCheckTimeout:  waitCheckTimeout,
		logger:            logger,
	}
}

func (cl *connLimiterWithTimeout) tryAcquire() bool {
	select {
	case cl.connLimit <- struct{}{}:
		return true
	default:
		return false
	}
}

func (cl *connLimiterWithTimeout) release() {
	<-cl.connLimit
}

func (cl *connLimiterWithTimeout) LimiterInterceptor(nextHandler ConnHandler) ConnHandler {
	return func(ctx context.Context, conn net.Conn) {
		deadline := time.Now().Add(cl.connAcceptTimeout)
		for {
			if time.Now().After(deadline) {
				if _, err := conn.Write([]byte(acceptConnDeadlineExceededError.Error())); err != nil {
					cl.logger.Error(ctx, "conn write error", zap.Error(err))
				}
				_ = conn.Close()
				return
			}
			if cl.tryAcquire() {
				break
			}

			// if limit exceeded try acquire letter
			time.Sleep(cl.waitCheckTimeout)
		}

		nextHandler(ctx, conn)
		cl.release()
	}
}
