package aclservice

import (
	"context"
	"github.com/google/uuid"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type ACLService struct {
	log        logger.Logger
	userClient port.UserClient
}

func New(log logger.Logger, userClient port.UserClient) *ACLService {
	return &ACLService{
		log:        log,
		userClient: userClient,
	}
}

func (r ACLService) CheckAccess(
	ctx context.Context,
	uow repository.UnitOfWork,
	userUUID uuid.UUID,
	requiredPermissions ...domain.PermissionKeyType,
) (bool, uint64, error) {

	user, err := r.userClient.GetByUUID(ctx, userUUID.String())
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.Error(), nil)
		return false, 0, err
	}

	roleKeys, err := uow.RoleRepository().GetRoleKeys(user.ID)
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.Error(), nil)
		return false, 0, err
	}

	for _, key := range roleKeys {
		if key == domain.RoleKeySuperAdmin {
			return true, user.ID, nil
		}
	}

	permissionKeys, err := uow.PermissionRepository().GetUserPermissionKeys(user.ID)
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.Error(), nil)
		return false, 0, err
	}

	for _, requiredPermission := range requiredPermissions {
		if requiredPermission == domain.PermissionKeyNone {
			return true, user.ID, nil
		}
		for _, key := range permissionKeys {
			if requiredPermission == key {
				return true, user.ID, nil
			}
		}
	}

	return false, 0, nil
}

func (r ACLService) AssignUserRoleToUser(uow repository.UnitOfWork, userID uint64) error {
	role, err := uow.RoleRepository().GetRoleUser()
	if err != nil {
		return err
	}

	roleIDs := make([]uint64, 1)
	roleIDs[0] = role.ID

	return uow.ACLRepository().AssignRolesToUser(userID, roleIDs)
}
