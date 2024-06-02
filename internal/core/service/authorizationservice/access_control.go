package authorizationservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type ACLService struct {
	log            logger.Logger
	permissionRepo port.PermissionRepository
	roleRepo       port.RoleRepository
	aclRepo        port.ACLRepository
	userClient     port.UserClient
}

func New(
	log logger.Logger,
	permissionRepo port.PermissionRepository,
	roleRepo port.RoleRepository,
	aclRepo port.ACLRepository,
	userClient port.UserClient,
) *ACLService {
	return &ACLService{
		log:            log,
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
		aclRepo:        aclRepo,
		userClient:     userClient,
	}
}

func (r ACLService) CheckAccess(
	ctx context.Context,
	userUUID uuid.UUID,
	permissions ...domain.PermissionKeyType,
) (bool, error) {

	user, err := r.userClient.GetByUUID(ctx, userUUID.String())
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.Error(), nil)
		return false, err
	}

	permissionKeys, err := r.permissionRepo.GetUserPermissionKeys(ctx, user.ID)
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.Error(), nil)
		return false, err
	}

	for _, key := range permissionKeys {
		for _, permission := range permissions {
			if permission == key {
				return true, nil
			}
		}
	}

	return false, nil
}

func (r ACLService) AddUserRole(ctx context.Context, userID uint64) error {
	role, err := r.roleRepo.GetUserRole(ctx)
	if err != nil {
		return err
	}

	return r.aclRepo.AddUserRole(ctx, userID, role.ID)
}
