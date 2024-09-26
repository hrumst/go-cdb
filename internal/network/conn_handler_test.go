package network

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hrumst/go-cdb/internal/tools"
)

type mockTcpConn struct {
	readData  [][]byte
	writeData [][]byte
	isClosed  int32
	deadline  time.Time
	readi     int
}

func (m *mockTcpConn) Read(b []byte) (n int, err error) {
	for {
		if time.Now().After(m.deadline) {
			return 0, io.EOF
		}
		if atomic.LoadInt32(&m.isClosed) == 1 {
			return 0, io.EOF
		}
		if m.readi == len(m.readData) {
			continue
		}

		rc := min(len(m.readData[m.readi]), len(b))
		copy(b, m.readData[m.readi][:rc])
		m.readi += 1
		return rc, nil
	}

}

func (m *mockTcpConn) Write(b []byte) (n int, err error) {
	m.writeData = append(m.writeData, b)
	return 0, err
}

func (m *mockTcpConn) Close() error {
	atomic.AddInt32(&m.isClosed, 1)
	return nil
}

func (m *mockTcpConn) LocalAddr() net.Addr {
	return nil
}

func (m *mockTcpConn) RemoteAddr() net.Addr {
	return nil
}

func (m *mockTcpConn) SetDeadline(t time.Time) error {
	m.deadline = t
	return nil
}

func (m *mockTcpConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockTcpConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestConnHandlerQueryHandler_TestDefaults(t *testing.T) {
	testLogger := tools.NewAppLogger(zap.NewNop())
	chq := NewConnHandlerQuery(nil, 0, 0, testLogger)
	assert.Equal(t, defaultIdleTimeoutServer, chq.connIdleTimeout)
	assert.Equal(t, int64(defaultMaxMessageSize), chq.maxBodySizeBytes)
}

func TestConnHandlerQueryHandler(t *testing.T) {
	testLogger := tools.NewAppLogger(zap.NewNop())

	var handlerCallCount int64
	testHandler := func(ctx context.Context, input string) string {
		atomic.AddInt64(&handlerCallCount, 1)
		return "response: " + input
	}

	chq := NewConnHandlerQuery(testHandler, 10, 100*time.Millisecond, testLogger)

	mockConn := mockTcpConn{}

	go chq.Handle(context.Background(), &mockConn)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		mockConn.readData = append(mockConn.readData, []byte("test"))
		time.Sleep(10 * time.Millisecond)
		mockConn.readData = append(mockConn.readData, []byte("test-more-than-limit"))
		time.Sleep(10 * time.Millisecond)
		mockConn.readData = append(mockConn.readData, []byte("test-2"))

		time.Sleep(110 * time.Millisecond)
		assert.Equal(t, int32(1), atomic.LoadInt32(&mockConn.isClosed))
	}()
	wg.Wait()

	assert.Equal(t, "response: test", string(mockConn.writeData[0]))
	assert.Equal(t, formatErrorResp(tooLongRequestError), string(mockConn.writeData[1]))
	assert.Equal(t, "response: test-2", string(mockConn.writeData[2]))

	assert.Equal(t, int64(2), handlerCallCount)
}
