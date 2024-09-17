package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hrumst/go-cdb/internal/database/compute"
)

func TestInputParserPlain_ParseStart(t *testing.T) {
	ip := NewInputParserPlain("   ")
	ip.ParseStart()
	assert.Equal(t, ip.posi, 2)

	ip2 := NewInputParserPlain(" test")
	ip2.ParseStart()
	assert.Equal(t, ip2.posi, 1)

	ip3 := NewInputParserPlain("test")
	ip3.ParseStart()
	assert.Equal(t, ip3.posi, 0)

	ip4 := NewInputParserPlain("")
	ip4.ParseStart()
	assert.Equal(t, ip4.posi, 0)
}

func TestInputParserPlain_ParseWhitespace(t *testing.T) {
	ip := NewInputParserPlain("   ")
	err := ip.ParseWhitespace()
	assert.NoError(t, err)
	assert.Equal(t, ip.posi, 2)

	ip2 := NewInputParserPlain(" test")
	err2 := ip2.ParseWhitespace()
	assert.NoError(t, err2)
	assert.Equal(t, ip2.posi, 1)

	ip3 := NewInputParserPlain("test test")
	ip3.posi = 4
	err3 := ip3.ParseWhitespace()
	assert.NoError(t, err3)
	assert.Equal(t, ip3.posi, 5)

	ip4 := NewInputParserPlain("")
	err4 := ip4.ParseWhitespace()
	assert.ErrorIs(t, err4, ParseError)
}

func TestInputParserPlain_ParseEnd(t *testing.T) {
	ip := NewInputParserPlain("   ")
	err := ip.ParseEnd()
	assert.NoError(t, err)
	assert.Equal(t, ip.posi, 2)

	ip2 := NewInputParserPlain("")
	err2 := ip2.ParseEnd()
	assert.NoError(t, err2)
	assert.Equal(t, ip2.posi, 0)

	ip3 := NewInputParserPlain(" garbage")
	err3 := ip3.ParseEnd()
	assert.ErrorIs(t, err3, ParseError)
}

func TestInputParserPlain_ParseCommand(t *testing.T) {
	type testCase struct {
		input      string
		initPosi   int
		expectPosi int
		expectCt   compute.CommandType
		expectErr  error
	}

	tcs := []testCase{
		{
			input:      "GET",
			expectPosi: 2,
			expectCt:   compute.CommandTypeGet,
		}, {
			input:      "GET key",
			expectPosi: 3,
			expectCt:   compute.CommandTypeGet,
		}, {
			input:      "   GET key",
			initPosi:   3,
			expectPosi: 6,
			expectCt:   compute.CommandTypeGet,
		}, {
			input:      "SET key",
			expectPosi: 3,
			expectCt:   compute.CommandTypeSet,
		}, {
			input:      "DEL key",
			expectPosi: 3,
			expectCt:   compute.CommandTypeDel,
		}, {
			input:     "GEET key",
			expectErr: invalidCommandTokenError,
		}, {
			input:     "GE",
			expectErr: invalidCommandTokenError,
		}, {
			input:     "",
			expectErr: invalidCommandTokenError,
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("ParseCommand_case_%d", i),
			func(t *testing.T) {
				ip := NewInputParserPlain(tc.input)
				ip.posi = tc.initPosi
				ct, err := ip.ParseCommandToken()
				if tc.expectErr == nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectPosi, ip.posi)
					assert.Equal(t, tc.expectCt, ct)
				} else {
					assert.ErrorIs(t, err, tc.expectErr)
				}
			},
		)
	}
}

func TestInputParserPlain_ParseKey(t *testing.T) {
	type testCase struct {
		input      string
		initPosi   int
		expectPosi int
		expectKey  string
		expectErr  error
	}

	tcs := []testCase{
		{
			input:      "*.667ky",
			expectPosi: 6,
			expectKey:  "*.667ky",
		}, {
			input:      "  k",
			initPosi:   2,
			expectPosi: 2,
			expectKey:  "k",
		}, {
			input:      "key k",
			expectPosi: 3,
			expectKey:  "key",
		}, {
			input:     " k",
			expectErr: invalidKeyError,
		}, {
			input:     " ",
			expectErr: invalidKeyError,
		}, {
			input:     "",
			expectErr: invalidKeyError,
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("ParseKey_case_%d", i),
			func(t *testing.T) {
				ip := NewInputParserPlain(tc.input)
				ip.posi = tc.initPosi
				key, err := ip.ParseKey()
				if tc.expectErr == nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectPosi, ip.posi)
					assert.Equal(t, tc.expectKey, key)
				} else {
					assert.ErrorIs(t, err, tc.expectErr)
				}
			},
		)
	}
}

func TestInputParserPlain_ParseVal(t *testing.T) {
	type testCase struct {
		input      string
		initPosi   int
		expectPosi int
		expectVal  string
		expectErr  error
	}

	tcs := []testCase{
		{
			input:      "*.667vy",
			expectPosi: 6,
			expectVal:  "*.667vy",
		}, {
			input:      "  v",
			initPosi:   2,
			expectPosi: 2,
			expectVal:  "v",
		}, {
			input:      "val v",
			expectPosi: 3,
			expectVal:  "val",
		}, {
			input:     " v",
			expectErr: invalidKeyValueError,
		}, {
			input:     " ",
			expectErr: invalidKeyValueError,
		}, {
			input:     "",
			expectErr: invalidKeyValueError,
		},
	}

	for i, tc := range tcs {
		t.Run(
			fmt.Sprintf("ParseVal_case_%d", i),
			func(t *testing.T) {
				ip := NewInputParserPlain(tc.input)
				ip.posi = tc.initPosi
				key, err := ip.ParseValue()
				if tc.expectErr == nil {
					assert.NoError(t, err)
					assert.Equal(t, tc.expectPosi, ip.posi)
					assert.Equal(t, tc.expectVal, key)
				} else {
					assert.ErrorIs(t, err, tc.expectErr)
				}
			},
		)
	}
}
