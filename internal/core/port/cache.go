package port

import (
	"context"
	"time"
)

type CacheDriver[T any] interface {
	Close()
	Get(ctx context.Context, key string) (*T, error)
	Set(ctx context.Context, key string, value *T, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	FlushAll(ctx context.Context)
	Remember(ctx context.Context, key string, expiration time.Duration, callback func() (*T, error)) (*T, error)
}
