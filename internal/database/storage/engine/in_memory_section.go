package engine

import (
	"context"
	"fmt"
	"sync"
)

type inMemoryEngineSection struct {
	mu   sync.RWMutex
	data map[string]string
}

func (ies *inMemoryEngineSection) Set(_ context.Context, key string, val string) error {
	ies.mu.Lock()
	defer ies.mu.Unlock()

	ies.data[key] = val
	return nil
}

func (ies *inMemoryEngineSection) Del(_ context.Context, key string) error {
	ies.mu.Lock()
	defer ies.mu.Unlock()

	_, exists := ies.data[key]
	if !exists {
		return nil
	}

	delete(ies.data, key)
	return nil
}

func (ies *inMemoryEngineSection) Get(_ context.Context, key string) (string, error) {
	ies.mu.RLock()
	defer ies.mu.RUnlock()

	val, exists := ies.data[key]
	if !exists {
		return "", fmt.Errorf("%w: inMemoryEngine.Get error: %w: %s", ExecError, NotFoundKey, key)
	}
	return val, nil
}
