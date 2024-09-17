package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ExecError   = errors.New("execution error")
	NotFoundKey = errors.New("not found key")
)

type inMemoryEngine struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewInMemoryEngine() *inMemoryEngine {
	return &inMemoryEngine{
		data: make(map[string]string),
	}
}

func (ime *inMemoryEngine) Set(_ context.Context, key string, val string) error {
	ime.mu.Lock()
	defer ime.mu.Unlock()

	ime.data[key] = val
	return nil
}

func (ime *inMemoryEngine) Del(_ context.Context, key string) error {
	ime.mu.Lock()
	defer ime.mu.Unlock()

	_, exists := ime.data[key]
	if !exists {
		return fmt.Errorf("%w: inMemoryEngine.Del error: %w: %s", ExecError, NotFoundKey, key)
	}

	delete(ime.data, key)
	return nil
}

func (ime *inMemoryEngine) Get(_ context.Context, key string) (string, error) {
	ime.mu.RLock()
	defer ime.mu.RUnlock()

	val, exists := ime.data[key]
	if !exists {
		return "", fmt.Errorf("%w: inMemoryEngine.Get error: %w: %s", ExecError, NotFoundKey, key)
	}
	return val, nil
}
