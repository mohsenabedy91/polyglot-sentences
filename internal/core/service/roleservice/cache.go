package roleservice

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type RoleCacheService struct {
	log       logger.Logger
	roleCache *redis.CacheDriver[domain.RoleKeyType]
}

func NewRoleCache(log logger.Logger, roleCache *redis.CacheDriver[any]) *RoleCacheService {
	return &RoleCacheService{
		log:       log,
		roleCache: (*redis.CacheDriver[domain.RoleKeyType])(roleCache),
	}
}

func (r RoleCacheService) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	key = fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, key)

	value, err := r.roleCache.Get(ctx, key)
	if err != nil {
		r.log.Error(logger.Cache, logger.RedisGet, err.Error(), map[logger.ExtraKey]interface{}{
			logger.CacheKey: key,
		})
		return nil, err
	}

	return value, nil
}

func (r RoleCacheService) SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error {
	for key, value := range items {
		key = fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, key)

		extra := map[logger.ExtraKey]interface{}{
			logger.CacheKey:    key,
			logger.CacheSetArg: value,
		}

		if err := r.roleCache.Set(ctx, key, &value, time.Duration(0)); err != nil {
			r.log.Error(logger.Cache, logger.RedisSet, fmt.Sprintf("Error setting value: %v", err), extra)
			return serviceerror.NewServerError()
		}
	}

	return nil
}
