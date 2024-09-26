package network

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net"
	"sync"
)

type (
	ConnHandler     func(ctx context.Context, conn net.Conn)
	ConnInterceptor func(nextHandler ConnHandler) ConnHandler
)

type tcpServer struct {
	address string
	logger  logger
}

func NewTCPServer(
	address string,
	logger logger,
) *tcpServer {
	return &tcpServer{
		address: address,
		logger:  logger,
	}
}

func (ts *tcpServer) Start(ctx context.Context, interceptors []ConnInterceptor, handler ConnHandler) error {
	listener, err := net.Listen("tcp", ts.address)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		if err := listener.Close(); err != nil {
			ts.logger.Error(ctx, "listener close error", zap.Error(err))
		}
	}()

	interceptChain := handler
	for ci := len(interceptors) - 1; ci >= 0; ci -= 1 {
		interceptChain = interceptors[ci](interceptChain)
	}

	wg := sync.WaitGroup{}
	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			interceptChain(ctx, conn)
		}()
	}
	wg.Wait()
	return nil
}
