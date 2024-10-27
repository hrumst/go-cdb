package engine

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRandString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestInMemory(t *testing.T) {
	type testCase struct {
		sectNum int
		key     string
		value   string
	}
	tcs := make(map[int][]testCase)

	testInMemoryStore := NewInMemoryEngine()
	for i := 0; i < 20; i += 1 {
		testKey := testRandString(20)
		testVal := "test-val-" + strconv.Itoa(i)
		testSectNum := testInMemoryStore.sectionNum(testKey)
		if _, e := tcs[testSectNum]; !e {
			tcs[testSectNum] = make([]testCase, 0)
		}
		tcs[testSectNum] = append(
			tcs[testSectNum],
			testCase{
				sectNum: testInMemoryStore.sectionNum(testKey),
				key:     testKey,
				value:   testVal,
			},
		)

		err := testInMemoryStore.Set(context.Background(), testKey, testVal)
		assert.NoError(t, err)
	}

	for i := 0; i < len(testInMemoryStore.dataStore); i += 1 {
		if _, e := tcs[i]; !e {
			assert.Empty(t, testInMemoryStore.dataStore[i].data)
		}
		assert.Equal(t, len(tcs[i]), len(testInMemoryStore.dataStore[i].data))
		for _, tc := range tcs[i] {
			res, err := testInMemoryStore.Get(context.Background(), tc.key)
			assert.Equal(t, tc.value, res)
			assert.NoError(t, err)

			res2, err2 := testInMemoryStore.dataStore[i].Get(context.Background(), tc.key)
			assert.Equal(t, tc.value, res2)
			assert.NoError(t, err2)
		}
	}

	for tci := range tcs {
		for _, tc := range tcs[tci] {
			err := testInMemoryStore.Del(context.Background(), tc.key)
			assert.NoError(t, err)
		}
	}

	for i := 0; i < len(testInMemoryStore.dataStore); i += 1 {
		assert.Empty(t, testInMemoryStore.dataStore[i].data)
	}
}
