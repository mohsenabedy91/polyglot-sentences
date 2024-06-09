package roleservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type Service struct {
	log       logger.Logger
	roleRepo  port.RoleRepository
	roleCache port.RoleCache
}

func New(log logger.Logger, roleRepo port.RoleRepository, roleCache port.RoleCache) *Service {
	return &Service{
		log:       log,
		roleRepo:  roleRepo,
		roleCache: roleCache,
	}
}

func (r *Service) Create(ctx context.Context, role domain.Role) error {

	key := helper.ConvertToUpperCase(role.Title)
	if exists, err := r.roleRepo.ExistKey(ctx, key); err != nil || !exists {
		return serviceerror.New(serviceerror.RoleExisted)
	}

	role.Key = domain.RoleKeyType(key)
	return r.roleRepo.Create(ctx, role)
}

func (r *Service) Get(ctx context.Context, uuidStr string) (*domain.Role, error) {
	return r.roleRepo.GetByUUID(ctx, uuid.MustParse(uuidStr))
}

func (r *Service) List(ctx context.Context) ([]*domain.Role, error) {
	roles, err := r.roleRepo.List(ctx)

	go func() {
		cacheRoles := make(map[string]domain.RoleKeyType)
		for _, role := range roles {
			if role.IsDefault == true {
				cacheRoles[role.UUID.String()] = role.Key
			}
		}

		_ = r.roleCache.SetBulk(ctx, cacheRoles)
	}()

	return roles, err
}

func (r *Service) Update(ctx context.Context, role domain.Role, uuidStr string) error {

	cachedRoleKey, err := r.roleCache.Get(ctx, uuidStr)
	if err != nil {
		return err
	}

	if cachedRoleKey != nil {
		role.Key = *cachedRoleKey
	} else {
		role.SetKey(role.Title)
	}

	return r.roleRepo.Update(ctx, role, uuid.MustParse(uuidStr))
}

func (r *Service) Delete(ctx context.Context, uuidStr string) error {
	return r.roleRepo.Delete(ctx, uuid.MustParse(uuidStr))
}

func (r *Service) GetPermissions(ctx context.Context, uuidStr string) (*domain.Role, error) {
	return r.roleRepo.GetPermissions(ctx, uuid.MustParse(uuidStr))
}
