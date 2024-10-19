package wal

import (
	"bytes"
	"encoding/gob"

	"github.com/hrumst/go-cdb/internal/database/compute"
	"github.com/hrumst/go-cdb/internal/tools"
)

type Log struct {
	UUId    string
	CmdType compute.CommandType
	Args    []string
}

func newLog(cmdType compute.CommandType, args []string) Log {
	return Log{
		UUId:    tools.GenerateUUIDv7(),
		CmdType: cmdType,
		Args:    args,
	}
}

func (l *Log) Encode(buffer *bytes.Buffer) error {
	encoder := gob.NewEncoder(buffer)
	return encoder.Encode(*l)
}

func (l *Log) Decode(buffer *bytes.Buffer) error {
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(l)
}
