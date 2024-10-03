package test_integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
	"time"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database"
	"github.com/hrumst/go-cdb/internal/database/compute/parser"
	"github.com/hrumst/go-cdb/internal/database/storage"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/network"
	"github.com/hrumst/go-cdb/internal/tools"
)

func TestNewTCPClient_Integration(t *testing.T) {
	testAddr := "localhost:11114"
	testLogger := tools.NewAppLogger(zap.NewNop())
	testServer := network.NewTCPServer(testAddr, testLogger)

	db := database.NewDatabase(
		storage.NewStorage(engine.NewInMemoryEngine()),
		parser.NewCommandExecParserPlain(),
		testLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := testServer.Start(
			ctx,
			[]network.ConnInterceptor{
				network.NewConnLimiterWithTimeout(
					10,
					100*time.Millisecond,
					testLogger,
				).LimiterInterceptor,
			},
			network.NewConnHandlerQuery(
				db.Execute,
				config.ParserSize("100b"),
				100*time.Millisecond,
				testLogger,
			).Handle,
		)
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	testClient, err := network.NewTcpClient(testAddr, 100*time.Millisecond, 100)
	assert.NoError(t, err)

	out, err := testClient.Request("GET key1")
	assert.NoError(t, err)
	assert.Equal(t, "Error: execution error: inMemoryEngine.Get error: not found key: key1", out)

	time.Sleep(90 * time.Millisecond)
	out, err = testClient.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)

	time.Sleep(90 * time.Millisecond)
	out, err = testClient.Request("GET key1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: val1", out)

	time.Sleep(90 * time.Millisecond)
	out, err = testClient.Request("DEL key1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)

	// idle timeout
	time.Sleep(110 * time.Millisecond)
	out, err = testClient.Request("DEL key1")
	assert.Error(t, err)
	assert.Empty(t, out)
}

func TestNewTCPClient_IntegrationConnLimit(t *testing.T) {
	testAddr := "localhost:11112"
	testLogger := tools.NewAppLogger(zap.NewNop())
	testServer := network.NewTCPServer(testAddr, testLogger)

	db := database.NewDatabase(
		storage.NewStorage(engine.NewInMemoryEngine()),
		parser.NewCommandExecParserPlain(),
		testLogger,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := testServer.Start(
			ctx,
			[]network.ConnInterceptor{
				network.NewConnLimiterWithTimeout(
					2,
					10*time.Millisecond,
					testLogger,
				).LimiterInterceptor,
			},
			network.NewConnHandlerQuery(
				db.Execute,
				config.ParserSize("100b"),
				100*time.Millisecond,
				testLogger,
			).Handle,
		)
		assert.NoError(t, err)
	}()

	time.Sleep(100 * time.Millisecond)

	testClient1, err := network.NewTcpClient(testAddr, 100*time.Millisecond, 100)
	assert.NoError(t, err)

	testClient2, err := network.NewTcpClient(testAddr, 100*time.Millisecond, 100)
	assert.NoError(t, err)

	testClient3, err := network.NewTcpClient(testAddr, 100*time.Millisecond, 100)
	assert.NoError(t, err)

	out, err := testClient1.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)

	out, err = testClient2.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)

	out, err = testClient3.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "accept connection deadline exceeded", out)

	// release connection
	err = testClient2.Close()
	assert.NoError(t, err)

	// recreate prev failed connection
	testClient3, err = network.NewTcpClient(testAddr, 100*time.Millisecond, 100)
	assert.NoError(t, err)
	out, err = testClient3.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)
}
