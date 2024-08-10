package userrepository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
)

type TOTPCacheRepository struct {
	log    logger.Logger
	conf   config.Config
	client *redis.Client
}

func NewTOTPCache(log logger.Logger, conf config.Config, driver *redis.Client) *TOTPCacheRepository {
	return &TOTPCacheRepository{
		log:    log,
		conf:   conf,
		client: driver,
	}
}

func (r TOTPCacheRepository) Set(ctx context.Context, key string, value string) error {
	key = fmt.Sprintf("%s:%s:%s", r.conf.Redis.Prefix, constant.RedisTOTPPrefix, strings.ToLower(key))

	extra := map[logger.ExtraKey]interface{}{
		logger.CacheKey:    key,
		logger.CacheSetArg: value,
	}

	if err := r.client.WithContext(ctx).Set(key, value, r.conf.TOTP.ExpireSecond).Err(); err != nil {
		r.log.Error(logger.Cache, logger.RedisSet, fmt.Sprintf("Error Set value: %v", err), extra)
		return serviceerror.NewServerError()
	}

	r.log.Info(logger.Cache, logger.RedisSet, "The TOTP state set successfully.", extra)

	return nil
}

func (r TOTPCacheRepository) Get(ctx context.Context, key string) (string, error) {
	key = fmt.Sprintf("%s:%s:%s", r.conf.Redis.Prefix, constant.RedisTOTPPrefix, strings.ToLower(key))

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
