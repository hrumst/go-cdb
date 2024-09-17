package storage

import "context"

type storageEngine interface {
	Set(ctx context.Context, key string, val string) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string) (string, error)
}

type storage struct {
	engine storageEngine
}

func NewStorage(engine storageEngine) *storage {
	return &storage{
		engine: engine,
	}
}

func (s storage) Set(ctx context.Context, key string, val string) error {
	return s.engine.Set(ctx, key, val)
}

func (s storage) Del(ctx context.Context, key string) error {
	return s.engine.Del(ctx, key)
}

func (s storage) Get(ctx context.Context, key string) (string, error) {
	return s.engine.Get(ctx, key)
}
