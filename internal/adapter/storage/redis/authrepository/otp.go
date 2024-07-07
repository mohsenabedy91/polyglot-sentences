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

type OTPCache struct {
	log    logger.Logger
	conf   config.Redis
	client *redis.Client
}

func NewOTPCache(log logger.Logger, conf config.Redis, driver *redis.Client) *OTPCache {
	return &OTPCache{
		log:    log,
		conf:   conf,
		client: driver,
	}
}

func (r OTPCache) Set(ctx context.Context, key string, value *domain.OTP, expiration time.Duration) error {
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

	r.log.Info(logger.Cache, logger.RedisSet, "The OTP state set successfully.", extra)

	return nil
}

func (r OTPCache) Get(ctx context.Context, key string) (*domain.OTP, error) {
	key = fmt.Sprintf("%s:%s", r.conf.Prefix, key)

	result, err := r.client.WithContext(ctx).Get(key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			r.log.Warn(logger.Cache, logger.RedisGet, fmt.Sprintf("Warn Get value: %v", err), nil)
			return nil, nil
		}

		r.log.Error(logger.Cache, logger.RedisGet, fmt.Sprintf("Error Get value: %v", err), nil)
		return nil, serviceerror.NewServerError()
	}

	var otp *domain.OTP
	if err = json.Unmarshal([]byte(result), &otp); err != nil {
		r.log.Error(logger.Cache, logger.RedisGet, fmt.Sprintf("Error Get value: %v", err), nil)
		return nil, serviceerror.NewServerError()
	}

	return otp, nil
}
