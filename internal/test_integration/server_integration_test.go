package test_integration

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database"
	"github.com/hrumst/go-cdb/internal/database/compute/parser"
	"github.com/hrumst/go-cdb/internal/database/storage"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/fs"
	"github.com/hrumst/go-cdb/internal/network"
	"github.com/hrumst/go-cdb/internal/tools"
)

type testStorage interface {
	Set(ctx context.Context, key string, val string) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
}

func InitTestStorage(t *testing.T, testDirPath string) testStorage {
	err := os.RemoveAll(testDirPath)
	if err != nil {
		t.Fatal(err)
	}
	err = os.Mkdir(testDirPath, 0755)
	if err != nil {
		t.Fatal(err)
	}

	dirFs := fs.NewDir(testDirPath)
	ts, err := storage.NewStorageWithRecover(
		context.Background(),
		engine.NewInMemoryEngine(),
		config.AppConfigWal{
			FlushBatchSize:    1 << 12,
			FlushBatchTimeout: 10 * time.Millisecond,
			MaxSegmentSize:    "10kb",
		},
		dirFs,
	)
	if err != nil {
		t.Fatal(err)
	}
	return ts
}

func TestNewTCPClient_Integration(t *testing.T) {
	testAddr := "localhost:11114"
	testLogger := tools.NewAppLogger(zap.NewNop())
	testServer := network.NewTCPServer(testAddr, testLogger)
	db := database.NewDatabase(
		InitTestStorage(t, "./test_dir"),
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
					10000*time.Millisecond,
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

	time.Sleep(80 * time.Millisecond)
	out, err = testClient.Request("SET key1 val1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: <empty>", out)

	time.Sleep(80 * time.Millisecond)
	out, err = testClient.Request("GET key1")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: val1", out)

	time.Sleep(80 * time.Millisecond)
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
		InitTestStorage(t, "./test_dir"),
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
