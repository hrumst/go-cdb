package network

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hrumst/go-cdb/internal/tools"
)

func TestConnLimiterWithTimeout_DefaultValues(t *testing.T) {
	testLogger := tools.NewAppLogger(zap.NewNop())
	testLimiter := NewConnLimiterWithTimeout(0, 0, testLogger)
	assert.Equal(t, defaultMaxConns, cap(testLimiter.connLimit))
	assert.Equal(t, defaultAcceptTimeout, testLimiter.connAcceptTimeout)
	assert.Equal(t, defaultAcceptTimeout/10, testLimiter.waitCheckTimeout)
}

func TestConnLimiterWithTimeout_AcquireRelease(t *testing.T) {
	testLogger := tools.NewAppLogger(zap.NewNop())
	testLimiter := NewConnLimiterWithTimeout(10, 100*time.Millisecond, testLogger)

	var acquireTrue, acquireFalse int64
	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i += 1 {
		go func() {
			defer wg.Done()
			if testLimiter.tryAcquire() {
				atomic.AddInt64(&acquireTrue, 1)
			} else {
				atomic.AddInt64(&acquireFalse, 1)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(10), atomic.LoadInt64(&acquireTrue))
	assert.Equal(t, int64(10), atomic.LoadInt64(&acquireFalse))

	for i := 0; i < 5; i += 1 {
		testLimiter.release()
	}

	var acquireTrue2, acquireFalse2 int64
	wg.Add(20)
	for i := 0; i < 20; i += 1 {
		go func() {
			defer wg.Done()
			if testLimiter.tryAcquire() {
				atomic.AddInt64(&acquireTrue2, 1)
			} else {
				atomic.AddInt64(&acquireFalse2, 1)
			}
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(5), atomic.LoadInt64(&acquireTrue2))
	assert.Equal(t, int64(15), atomic.LoadInt64(&acquireFalse2))
}

func TestConnLimiterWithTimeout_AcquireWithTimeout(t *testing.T) {
	testLogger := tools.NewAppLogger(zap.NewNop())
	testLimiter := NewConnLimiterWithTimeout(10, time.Millisecond*100, testLogger)

	var nextReachedCounter int64
	mockNext := func(ctx context.Context, conn net.Conn) {
		atomic.AddInt64(&nextReachedCounter, 1)
		time.Sleep(time.Millisecond * 90)
	}
	testLimiterInterceptor := testLimiter.LimiterInterceptor(mockNext)

	// total requests: 30
	// first 20 will pass by 10 in 2 waves
	// last 10 will not pass because of timeout
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(30)
		for i := 0; i < 30; i += 1 {
			go func() {
				defer wg.Done()

				mockConn, _ := net.Pipe()
				testLimiterInterceptor(context.Background(), mockConn)
			}()
		}
		wg.Wait()
	}()

	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, int64(10), atomic.LoadInt64(&nextReachedCounter))

	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, int64(20), atomic.LoadInt64(&nextReachedCounter))
}
