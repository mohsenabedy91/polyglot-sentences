package roleservice

import (
	"context"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type RoleCacheService struct {
	cache port.RoleCache
}

func NewRoleCacheService(cache port.RoleCache) RoleCacheService {
	return RoleCacheService{
		cache: cache,
	}
}

func (r RoleCacheService) Get(ctx context.Context, key string) (*domain.RoleKeyType, error) {
	key = fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, key)
	return r.cache.Get(ctx, key)
}

func (r RoleCacheService) SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error {
	for key, value := range items {
		key = fmt.Sprintf("%s:%s", constant.RoleKeyPrefix, key)
		if err := r.cache.Set(ctx, key, &value, time.Duration(0)); err != nil {
			return serviceerror.NewServerError()
		}
	}

	return nil
}
