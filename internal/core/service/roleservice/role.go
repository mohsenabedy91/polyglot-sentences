package roleservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type Service struct {
	log       logger.Logger
	roleCache port.RoleCache
}

func New(log logger.Logger, roleCache port.RoleCache) *Service {
	return &Service{
		log:       log,
		roleCache: roleCache,
	}
}

func (r *Service) Create(uow port.AuthUnitOfWork, role domain.Role) error {
	role.SetKey(role.Title)

	if exists, err := uow.RoleRepository().ExistKey(role.Key); err != nil || !exists {
		return serviceerror.New(serviceerror.RoleExisted)
	}

	return uow.RoleRepository().Create(role)
}

func (r *Service) Get(uow port.AuthUnitOfWork, uuidStr string) (*domain.Role, error) {
	return uow.RoleRepository().GetByUUID(uuid.MustParse(uuidStr))
}

func (r *Service) List(ctx context.Context, uow port.AuthUnitOfWork) ([]*domain.Role, error) {
	roles, err := uow.RoleRepository().List()

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

func (r *Service) Update(ctx context.Context, uow port.AuthUnitOfWork, role domain.Role, uuidStr string) error {

	if cachedRoleKey, err := r.roleCache.Get(ctx, uuidStr); err != nil {
		return err
	} else {
		if cachedRoleKey != nil {
			role.Key = *cachedRoleKey
		} else {
			role.SetKey(role.Title)
		}
	}

	if exists, err := uow.RoleRepository().ExistKey(role.Key); err != nil || !exists {
		return serviceerror.New(serviceerror.RoleExisted)
	}

	return uow.RoleRepository().Update(role, uuid.MustParse(uuidStr))
}

func (r *Service) Delete(uow port.AuthUnitOfWork, uuidStr string, deletedBy uint64) error {
	return uow.RoleRepository().Delete(uuid.MustParse(uuidStr), deletedBy)
}

func (r *Service) GetPermissions(uow port.AuthUnitOfWork, uuidStr string) (*domain.Role, error) {
	return uow.RoleRepository().GetPermissions(uuid.MustParse(uuidStr))
}

func (r *Service) SyncPermissions(uow port.AuthUnitOfWork, uuidStr string, permissionUUIDStr []string) error {

	permissionUUIDs := make([]uuid.UUID, len(permissionUUIDStr))
	for i, p := range permissionUUIDStr {
		parsedUUID, err := uuid.Parse(p)
		if err != nil {
			return serviceerror.New(serviceerror.InvalidRequestBody)
		}
		permissionUUIDs[i] = parsedUUID
	}

	validPermissions, err := uow.PermissionRepository().FilterValidPermissions(permissionUUIDs)
	if err != nil {
		return err
	}

	roleUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		return serviceerror.New(serviceerror.InvalidRequestBody)
	}

	role, err := uow.RoleRepository().GetByUUID(roleUUID)
	if err != nil {
		return err
	}

	return uow.RoleRepository().SyncPermissions(role.ID, validPermissions)
}
