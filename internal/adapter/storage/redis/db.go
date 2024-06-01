package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/base64x"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type Interface[T any] interface {
	Get(key string) (destination T, err error)
	Set(key string, value T, expiration time.Duration) (err error)
	Delete(key string) (err error)
	FlushAll()
	Remember(key string, expiration time.Duration, callback func() (T, error)) (T, error)
}

type CacheDriver[T any] struct {
	log    logger.Logger
	cfg    config.Config
	client *redis.Client
}

func NewCacheDriver[T any](log logger.Logger, cfg config.Config) (*CacheDriver[T], error) {

	client := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		Password:           cfg.Redis.Password,
		DB:                 cfg.Redis.DB,
		DialTimeout:        cfg.Redis.DialTimeout * time.Second,
		ReadTimeout:        cfg.Redis.ReadTimeout * time.Second,
		WriteTimeout:       cfg.Redis.WriteTimeout * time.Second,
		PoolSize:           cfg.Redis.PoolSize,
		PoolTimeout:        cfg.Redis.PoolTimeout * time.Second,
		IdleTimeout:        cfg.Redis.IdleTimeout * time.Millisecond,
		IdleCheckFrequency: cfg.Redis.IdleCheckFrequency * time.Millisecond,
	})

	if val, err := client.Ping().Result(); err != nil {
		log.Error(
			logger.Redis,
			logger.RedisPing,
			fmt.Sprintf("Redis client doesn't connected val: %s, error: %v", val, err.Error()),
			nil,
		)
		return nil, serviceerror.New(serviceerror.ServiceUnavailable)
	}

	log.Info(logger.Redis, logger.Startup, "Redis client initialized", nil)

	return &CacheDriver[T]{
		log:    log,
		client: client,
		cfg:    cfg,
	}, nil
}

func (r *CacheDriver[T]) Get(ctx context.Context, key string) (T, error) {
	var destination T

	key = fmt.Sprintf("%s:%s", r.cfg.Redis.Prefix, key)
	v, err := r.client.WithContext(ctx).Get(key).Result()
	if err != nil {
		return destination, serviceerror.NewServerError()
	}

	err = json.Unmarshal([]byte(v), &destination)
	if err != nil {
		return destination, serviceerror.NewServerError()
	}

	return destination, nil
}

func (r *CacheDriver[T]) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	key = fmt.Sprintf("%s:%s", r.cfg.Redis.Prefix, key)

	data, err := json.Marshal(value)
	if err != nil {
		r.log.Error(logger.Redis, logger.RedisSet, fmt.Sprintf("Error marshalling value, error: %v", err), nil)
		return serviceerror.NewServerError()
	}

	if err = r.client.WithContext(ctx).Set(key, data, expiration).Err(); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}

func (r *CacheDriver[T]) Delete(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s:%s", r.cfg.Redis.Prefix, key)
	if err := r.client.WithContext(ctx).Del(key).Err(); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}

func (r *CacheDriver[T]) FlushAll(ctx context.Context) {
	if r.cfg.App.Debug {
		r.client.WithContext(ctx).FlushAll()
	}
}

func (r *CacheDriver[T]) Remember(ctx context.Context, key string, expiration time.Duration, callback func() (T, error)) (T, error) {
	var destination T

	key = fmt.Sprintf("%s:%s", r.cfg.Redis.Prefix, key)
	destination, err := r.Get(ctx, key)
	if err != nil {

		if destination, err = callback(); err != nil {
			return destination, serviceerror.NewServerError()
		}

		if err = r.Set(ctx, key, destination, expiration); err != nil {
			return destination, serviceerror.NewServerError()
		}
	}
	return destination, nil
}

func (r *CacheDriver[T]) Subscribe(ctx context.Context, key string) (T, error) {
	var destination T

	subscriber := r.client.WithContext(ctx).Subscribe(key)
	for {
		message, err := subscriber.ReceiveMessage()
		if err != nil {
			return destination, serviceerror.NewServerError()
		}

		payloadStr, err := base64x.StdEncoding.DecodeString(message.Payload)
		if err != nil {
			return destination, serviceerror.NewServerError()
		}

		if err = json.Unmarshal(payloadStr, &destination); err != nil {
			return destination, serviceerror.NewServerError()
		}

		return destination, nil
	}
}

func (r *CacheDriver[T]) Publish(ctx context.Context, key string, value []byte) error {
	payloadStr := base64x.StdEncoding.EncodeToString(value)
	if err := r.client.WithContext(ctx).Publish(key, payloadStr).Err(); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}
