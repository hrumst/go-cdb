package config

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestParserSize(t *testing.T) {
	type tc struct {
		input  string
		result int64
	}

	tcs := []tc{
		{
			input:  "100Gb",
			result: 100 << 30,
		}, {
			input:  "3mb",
			result: 3 << 20,
		}, {
			input:  "76kB",
			result: 76 << 10,
		}, {
			input:  "96b",
			result: 96,
		}, {
			input:  "0Gb",
			result: 0,
		}, {
			input:  "test",
			result: 0,
		},
	}

	for tci := 0; tci < len(tcs); tci += 1 {
		t.Run(
			"test_case_"+strconv.Itoa(tci),
			func(t *testing.T) {
				assert.Equal(t, tcs[tci].result, ParserSize(tcs[tci].input))
			},
		)
	}

}
