package network

import (
	"context"
	"go.uber.org/zap"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/hrumst/go-cdb/internal/tools"
	"github.com/stretchr/testify/assert"
)

func TestTCPServer(t *testing.T) {
	testAddr := "localhost:11113"
	testLogger := tools.NewAppLogger(zap.NewNop())
	testServer := NewTCPServer(testAddr, testLogger)

	ctx, cancel := context.WithCancel(context.Background())

	handlerResult := make([]string, 0)
	testHandler := func(ctx context.Context, conn net.Conn) {
		readBuff := make([]byte, 100)

		count, err := conn.Read(readBuff)
		assert.NoError(t, err)
		handlerResult = append(handlerResult, string(readBuff[:count]))

		count, err = conn.Read(readBuff)
		assert.NoError(t, err)
		handlerResult = append(handlerResult, string(readBuff[:count]))
	}

	connCnt := 0
	testInterceptor := func(nextHandler ConnHandler) ConnHandler {
		return func(ctx context.Context, conn net.Conn) {
			connCnt += 1
			nextHandler(ctx, conn)
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := testServer.Start(ctx, []ConnInterceptor{testInterceptor}, testHandler)
		assert.NoError(t, err)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(100 * time.Millisecond)
		testConn, err := net.Dial("tcp", testAddr)
		assert.NoError(t, err)

		_, err = testConn.Write([]byte("test-1"))
		assert.NoError(t, err)

		time.Sleep(100 * time.Millisecond)

		_, err = testConn.Write([]byte("test-2"))
		assert.NoError(t, err)
		testConn.Close()

		// stop server
		cancel()
	}()

	wg.Wait()

	assert.Equal(t, []string{"test-1", "test-2"}, handlerResult)
	assert.Equal(t, 1, connCnt)
}
