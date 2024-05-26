package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/base64x"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
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
	client *redis.Client
	log    logger.Logger
	config config.Config
}

func NewRedisCacheDriver[T any](config config.Config, log logger.Logger) (*CacheDriver[T], error) {

	client := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port),
		Password:           config.Redis.Password,
		DB:                 config.Redis.DB,
		DialTimeout:        config.Redis.DialTimeout * time.Second,
		ReadTimeout:        config.Redis.ReadTimeout * time.Second,
		WriteTimeout:       config.Redis.WriteTimeout * time.Second,
		PoolSize:           config.Redis.PoolSize,
		PoolTimeout:        config.Redis.PoolTimeout * time.Second,
		IdleTimeout:        config.Redis.IdleTimeout * time.Millisecond,
		IdleCheckFrequency: config.Redis.IdleCheckFrequency * time.Millisecond,
	})

	if _, err := client.Ping().Result(); err != nil {
		log.Error(logger.Redis, logger.RedisPing, fmt.Sprintf("Redis client doesn't connected %v", err.Error()), nil)
		return nil, err
	}

	log.Info(logger.Redis, logger.Startup, "Redis client initialized", nil)

	return &CacheDriver[T]{
		log:    log,
		client: client,
		config: config,
	}, nil
}

func (r *CacheDriver[T]) Get(key string) (destination T, err error) {
	key = fmt.Sprintf("%s:%s", r.config.Redis.Prefix, key)
	v, err := r.client.Get(key).Result()
	if err != nil {
		return destination, err
	}

	err = json.Unmarshal([]byte(v), &destination)
	if err != nil {
		return destination, err
	}

	return destination, nil
}

func (r *CacheDriver[T]) Set(key string, value interface{}, expiration time.Duration) (err error) {
	key = fmt.Sprintf("%s:%s", r.config.Redis.Prefix, key)

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(key, data, expiration).Err()
}

func (r *CacheDriver[T]) Delete(key string) (err error) {
	key = fmt.Sprintf("%s:%s", r.config.Redis.Prefix, key)
	return r.client.Del(key).Err()
}

func (r *CacheDriver[T]) FlushAll() {
	if r.config.App.Debug {
		r.client.FlushAll()
	}
}

func (r *CacheDriver[T]) Remember(key string, expiration time.Duration, callback func() (T, error)) (destination T, err error) {
	key = fmt.Sprintf("%s:%s", r.config.Redis.Prefix, key)
	destination, err = r.Get(key)
	if err != nil {
		destination, err = callback()
		if err != nil {
			return destination, err
		}
		err = r.Set(key, destination, expiration)
		if err != nil {
			return destination, err
		}
	}
	return destination, nil
}

func (r *CacheDriver[T]) Subscribe(ctx context.Context, key string) (destination T, err error) {
	subscriber := r.client.WithContext(ctx).Subscribe(key)

	for {
		message, err := subscriber.ReceiveMessage()
		if err != nil {
			return destination, err
		}

		payloadStr, _ := base64x.StdEncoding.DecodeString(message.Payload)

		err = json.Unmarshal(payloadStr, &destination)
		if err != nil {
			return destination, err
		}

		return destination, nil
	}
}

func (r *CacheDriver[T]) Publish(ctx context.Context, key string, value []byte) (err error) {
	payloadStr := base64x.StdEncoding.EncodeToString(value)
	return r.client.WithContext(ctx).Publish(key, payloadStr).Err()
}
