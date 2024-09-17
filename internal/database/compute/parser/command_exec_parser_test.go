package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hrumst/go-cdb/internal/database/compute"
)

func TestCommandExecParser_WithPlainInputParser(t *testing.T) {
	type testCase struct {
		input             string
		expectCommandExec compute.CommandExec
		expectErr         error
	}

	tcs := []testCase{
		{
			input: "GET key1",
			expectCommandExec: compute.CommandExec{
				Command: compute.CommandTypeGet,
				Key:     "key1",
			},
		}, {
			input: "  GET   key1 ",
			expectCommandExec: compute.CommandExec{
				Command: compute.CommandTypeGet,
				Key:     "key1",
			},
		}, {
			input: "SET key1 val1",
			expectCommandExec: compute.CommandExec{
				Command: compute.CommandTypeSet,
				Key:     "key1",
				Val:     "val1",
			},
		}, {
			input: "  SET   key1    val1  ",
			expectCommandExec: compute.CommandExec{
				Command: compute.CommandTypeSet,
				Key:     "key1",
				Val:     "val1",
			},
		}, {
			input: "DEL key1",
			expectCommandExec: compute.CommandExec{
				Command: compute.CommandTypeDel,
				Key:     "key1",
			},
		}, {
			input:     "gGET key1",
			expectErr: invalidCommandTokenError,
		}, {
			input:     "GETk ey1",
			expectErr: noWhitespaceError,
		}, {
			input:     "GET key1 1",
			expectErr: unexpectedInputEndError,
		}, {
			input:     "GET ",
			expectErr: invalidKeyError,
		}, {
			input:     "SET key1 ",
			expectErr: invalidKeyValueError,
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("ParseCommand_case_%d", i),
			func(t *testing.T) {
				cmdParser := NewCommandExecParserPlain()
				cmdExec, err := cmdParser.Parse(tc.input)
				if tc.expectErr == nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectCommandExec, *cmdExec)
				} else {
					assert.ErrorIs(t, err, tc.expectErr)
				}
			},
		)
	}
}
