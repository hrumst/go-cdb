package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"

	"github.com/hrumst/go-cdb/internal/cli"
	"github.com/hrumst/go-cdb/internal/config"
	"github.com/hrumst/go-cdb/internal/network"
)

const (
	cliName = "cdbTCPClient"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		os.Exit(1)
	}

	var netAppCfp config.AppConfigNetwork
	flag.StringVar(
		&netAppCfp.Address,
		"address",
		"localhost:3223",
		"host database address",
	)
	flag.DurationVar(
		&netAppCfp.IdleTimeout,
		"idle_timeout",
		time.Minute,
		"idle timeout for database conn",
	)
	flag.StringVar(
		&netAppCfp.MaxMessageSize,
		"max_message_size",
		"4KB",
		"max message size for database conn",
	)
	flag.Parse()

	tcpClient, err := network.NewTcpClient(
		netAppCfp.Address,
		netAppCfp.IdleTimeout,
		config.ParserSize(netAppCfp.MaxMessageSize),
	)
	if err != nil {
		zapLogger.Fatal("init tcp client error", zap.Error(err))
	}
	defer func() {
		_ = tcpClient.Close()
	}()

	replServer := cli.NewREPL(
		cliName,
		int(config.ParserSize(netAppCfp.MaxMessageSize)),
		os.Stdin,
		os.Stdout,
		func(input string) (string, error) {
			return tcpClient.Request(input)
		},
	)
	if err := replServer.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
