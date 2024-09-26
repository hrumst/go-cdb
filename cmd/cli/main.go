package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"github.com/hrumst/go-cdb/internal/cli"
	"github.com/hrumst/go-cdb/internal/database"
	"github.com/hrumst/go-cdb/internal/database/compute/parser"
	"github.com/hrumst/go-cdb/internal/database/storage"
	"github.com/hrumst/go-cdb/internal/database/storage/engine"
	"github.com/hrumst/go-cdb/internal/tools"
)

const (
	cliName    = "cdbREPL"
	inputLimit = 1 << 12
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}
	appLogger := tools.NewAppLogger(zapLogger)

	db := database.NewDatabase(
		storage.NewStorage(engine.NewInMemoryEngine()),
		parser.NewCommandExecParserPlain(),
		appLogger,
	)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		fmt.Println("\nReceived an interrupt, exiting...")
		os.Exit(0)
	}()

	replServer := cli.NewREPL(
		cliName,
		inputLimit,
		os.Stdin,
		os.Stdout,
		func(input string) (string, error) {
			return db.Execute(context.Background(), input), nil
		},
	)
	if err := replServer.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
