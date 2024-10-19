package wal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hrumst/go-cdb/internal/database/compute"
)

func TestWalLogInit(t *testing.T) {
	testLog1 := newLog(compute.CommandTypeSet, []string{"arg1", "arg2"})
	assert.Equal(t, compute.CommandTypeSet, testLog1.CmdType)
	assert.Equal(t, []string{"arg1", "arg2"}, testLog1.Args)
	assert.Len(t, testLog1.UUId, 36)
}

func TestWalLogEncodeDecode(t *testing.T) {
	testLog1 := newLog(compute.CommandTypeSet, []string{"arg1", "arg2"})
	testLog2 := newLog(compute.CommandTypeGet, []string{"arg3", "arg4"})
	testBuffer := bytes.NewBuffer([]byte{})

	err := testLog1.Encode(testBuffer)
	assert.NoError(t, err)

	err = testLog2.Encode(testBuffer)
	assert.NoError(t, err)
	assert.Len(t, testBuffer.Bytes(), 254)

	var (
		testLogDec1,
		testLogDec2 Log
	)
	err = testLogDec1.Decode(testBuffer)
	assert.NoError(t, err)
	assert.Equal(t, testLogDec1.CmdType, compute.CommandTypeSet)
	assert.Equal(t, testLogDec1.Args, []string{"arg1", "arg2"})

	err = testLogDec2.Decode(testBuffer)
	assert.Equal(t, testLogDec2.CmdType, compute.CommandTypeGet)
	assert.Equal(t, testLogDec2.Args, []string{"arg3", "arg4"})
	assert.NoError(t, err)
}
