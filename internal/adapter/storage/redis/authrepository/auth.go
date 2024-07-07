package authrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type AuthCache struct {
	log    logger.Logger
	conf   config.Redis
	client *redis.Client
}

func NewAuthCache(log logger.Logger, conf config.Redis, driver *redis.Client) *AuthCache {
	return &AuthCache{
		log:    log,
		conf:   conf,
		client: driver,
	}
}

func (r AuthCache) SetTokenState(ctx context.Context, key string, value string, expiration time.Duration) error {
	key = fmt.Sprintf("%s:%s", r.conf.Prefix, key)

	extra := map[logger.ExtraKey]interface{}{
		logger.CacheKey:    key,
		logger.CacheSetArg: value,
	}

	if err := r.client.WithContext(ctx).Set(key, value, expiration).Err(); err != nil {
		r.log.Error(logger.Cache, logger.RedisSet, fmt.Sprintf("Error Set value: %v", err), extra)
		return serviceerror.NewServerError()
	}

	return nil
}

func (r AuthCache) GetTokenState(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("%s:%s", r.conf.Prefix, key)

	result, err := r.client.WithContext(ctx).Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.log.Warn(logger.Cache, logger.RedisGet, fmt.Sprintf("Warn Get value: %v", err), nil)
			return "", nil
		}

		r.log.Error(logger.Cache, logger.RedisGet, fmt.Sprintf("Error Get value: %v", err), nil)
		return "", serviceerror.NewServerError()
	}

	return result, nil
}
