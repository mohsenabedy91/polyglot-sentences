package authrepository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type RoleCache struct {
	log    logger.Logger
	conf   config.Redis
	client *redis.Client
}

func NewRoleCache(log logger.Logger, conf config.Redis, driver *redis.Client) *RoleCache {
	return &RoleCache{
		log:    log,
		conf:   conf,
		client: driver,
	}
}

func (r RoleCache) Set(ctx context.Context, key string, value *domain.RoleKeyType, expiration time.Duration) error {
	key = fmt.Sprintf("%s:%s", r.conf.Prefix, key)

	extra := map[logger.ExtraKey]interface{}{
		logger.CacheKey:    key,
		logger.CacheSetArg: value,
	}

	data, marshalErr := json.Marshal(value)
	if marshalErr != nil {
		r.log.Error(logger.Cache, logger.RedisSet, fmt.Sprintf("Error marshalling value: %v", marshalErr), extra)
		return serviceerror.NewServerError()
	}

	if err := r.client.WithContext(ctx).Set(key, data, expiration).Err(); err != nil {
		r.log.Error(logger.Cache, logger.RedisSet, fmt.Sprintf("Error Set value: %v", err), extra)
		return serviceerror.NewServerError()
	}

	r.log.Info(logger.Cache, logger.RedisSet, "The Role set successfully.", extra)

	return nil
}

func (r RoleCache) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	key = fmt.Sprintf("%s:%s", r.conf.Prefix, key)

	result, err := r.client.WithContext(ctx).Get(key).Result()
	if err != nil {
		extra := map[logger.ExtraKey]interface{}{
			logger.CacheKey: key,
		}
		if errors.Is(err, redis.Nil) {
			r.log.Warn(logger.Cache, logger.RedisGet, fmt.Sprintf("Warn Get value: %v", err), extra)
			return nil, nil
		}

		r.log.Error(logger.Cache, logger.RedisGet, fmt.Sprintf("Error Get value: %v", err), extra)
		return nil, serviceerror.NewServerError()
	}

	var roleKey *domain.RoleKeyType
	if err = json.Unmarshal([]byte(result), &roleKey); err != nil {
		r.log.Error(logger.Cache, logger.RedisGet, fmt.Sprintf("Error Get value: %v", err), map[logger.ExtraKey]interface{}{
			logger.CacheKey: key,
		})
		return nil, serviceerror.NewServerError()
	}

	return roleKey, nil
}
