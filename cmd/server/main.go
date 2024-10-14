package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/database"
	"github.com/hrumst/go-cdb/internal/database/compute/parser"
	dbStorage "github.com/hrumst/go-cdb/internal/database/storage"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/fs"
	"github.com/hrumst/go-cdb/internal/network"
	"github.com/hrumst/go-cdb/internal/tools"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	appLogger := tools.NewAppLogger(zapLogger)

	signalChannel := make(chan os.Signal, 1)
	signalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		cancel()
		fmt.Println("\nReceived an interrupt, exiting...")
		os.Exit(0)
	}()

	appConfig, err := config.ParseAppConfigFromEnv()
	if err != nil {
		zapLogger.Fatal("parse app config error", zap.Error(err))
	}

	storage, err := dbStorage.NewStorageWithRecover(
		signalCtx,
		engine.NewInMemoryEngine(),
		appConfig.Wal,
		fs.NewDir(appConfig.Wal.DataDirectoryPath),
	)
	if err != nil {
		zapLogger.Fatal("storage init error", zap.Error(err))
	}

	db := database.NewDatabase(storage, parser.NewCommandExecParserPlain(), appLogger)
	tcpServer := network.NewTCPServer(appConfig.Network.Address, appLogger)
	if err := tcpServer.Start(
		signalCtx,
		[]network.ConnInterceptor{
			network.NewConnLimiterWithTimeout(
				appConfig.Network.MaxConnections,
				appConfig.Network.AcceptTimeout,
				appLogger,
			).LimiterInterceptor,
		},
		network.NewConnHandlerQuery(
			db.Execute,
			config.ParserSize(appConfig.Network.MaxMessageSize),
			appConfig.Network.IdleTimeout,
			appLogger,
		).Handle,
	); err != nil {
		zapLogger.Fatal("start tcp server error", zap.Error(err))
	}
}
