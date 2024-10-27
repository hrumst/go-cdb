package engine

import (
	"context"
	"errors"
)

var (
	ExecError   = errors.New("execution error")
	NotFoundKey = errors.New("not found key")
)

const defaultSectionsNumber = 1 << 5

type inMemoryEngine struct {
	dataStore []*inMemoryEngineSection
}

func NewInMemoryEngine() *inMemoryEngine {
	newEngine := inMemoryEngine{
		dataStore: make([]*inMemoryEngineSection, defaultSectionsNumber),
	}
	for i := 0; i < defaultSectionsNumber; i += 1 {
		newEngine.dataStore[i] = &inMemoryEngineSection{
			data: make(map[string]string),
		}
	}
	return &newEngine
}

func (ie *inMemoryEngine) sectionNum(key string) int {
	var sectionNum int
	for i := 0; i < len(key); i += 1 {
		sectionNum += int(key[i])
	}
	return sectionNum % len(ie.dataStore)
}

func (ie *inMemoryEngine) sectionByKey(key string) *inMemoryEngineSection {
	return ie.dataStore[ie.sectionNum(key)]
}

func (ie *inMemoryEngine) Set(ctx context.Context, key string, val string) error {
	return ie.sectionByKey(key).Set(ctx, key, val)
}

func (ie *inMemoryEngine) Del(ctx context.Context, key string) error {
	return ie.sectionByKey(key).Del(ctx, key)
}

func (ie *inMemoryEngine) Get(ctx context.Context, key string) (string, error) {
	return ie.sectionByKey(key).Get(ctx, key)
}
