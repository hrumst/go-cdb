package storage

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/hrumst/go-cdb/internal/database/storage/engine"
)

func TestStorage_InMemoryEngine(t *testing.T) {
	const (
		testKey = "key1"
		testVal = "val1"
	)

	testStorage := NewStorage(engine.NewInMemoryEngine())
	result, err := testStorage.Get(context.Background(), testKey)
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, engine.NotFoundKey)

	err = testStorage.Del(context.Background(), testKey)
	assert.NoError(t, err, engine.NotFoundKey)

	err = testStorage.Set(context.Background(), testKey, testVal)
	assert.NoError(t, err)

	result, err = testStorage.Get(context.Background(), testKey)
	assert.Equal(t, testVal, result)
	assert.NoError(t, err)

	err = testStorage.Del(context.Background(), testKey)
	assert.NoError(t, err)

	result, err = testStorage.Get(context.Background(), testKey)
	assert.Equal(t, "", result)
	assert.ErrorIs(t, err, engine.NotFoundKey)
}
