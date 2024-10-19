package test_integration

import (
	"context"
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

func InitTestStorageRecover(t *testing.T, testDirPath string) testStorage {
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

func Test(t *testing.T) {
	testAddr := "localhost:11115"
	testLogger := tools.NewAppLogger(zap.NewNop())
	testServer := network.NewTCPServer(testAddr, testLogger)
	db := database.NewDatabase(
		InitTestStorageRecover(t, "./test_dir_2"),
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

	out, err := testClient.Request("GET key2")
	assert.NoError(t, err)
	assert.Equal(t, "Ok: val1", out)
}
