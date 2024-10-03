package network

import (
	"fmt"
	"net"
	"time"
)

const (
	defaultMaxBufferSize     = 1 << 12
	defaultIdleTimeoutClient = time.Minute * 2
)

type tcpClient struct {
	conn        net.Conn
	maxBodySize int64
	idleTimeout time.Duration
}

func NewTcpClient(address string, idleTimeout time.Duration, maxBodySize int64) (*tcpClient, error) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("dial connection error: %w", err)
	}

	if maxBodySize < 1 {
		maxBodySize = int64(defaultMaxBufferSize)
	}
	if idleTimeout < 1 {
		idleTimeout = defaultIdleTimeoutClient
	}

	return &tcpClient{
		conn:        conn,
		maxBodySize: maxBodySize,
		idleTimeout: idleTimeout,
	}, nil
}

func (tc *tcpClient) Request(input string) (string, error) {
	if tc.idleTimeout > 0 {
		if err := tc.conn.SetDeadline(time.Now().Add(tc.idleTimeout)); err != nil {
			return "", fmt.Errorf("SetDeadline error: %w", err)
		}
	}

	if _, err := tc.conn.Write([]byte(input)); err != nil {
		return "", fmt.Errorf("sending message error: %w", err)
	}

	readBuffer := make([]byte, tc.maxBodySize)
	count, err := tc.conn.Read(readBuffer)
	if err != nil {
		return "", fmt.Errorf("reading message error: %w", err)
	}
	return string(readBuffer[:count]), nil
}

func (tc *tcpClient) Close() error {
	if err := tc.conn.Close(); err != nil {
		return fmt.Errorf("closing connection error: %w", err)
	}
	return nil
}
