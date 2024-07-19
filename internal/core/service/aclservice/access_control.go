package aclservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
)

type ACLService struct {
	userClient port.UserClient
}

func New(userClient port.UserClient) *ACLService {
	return &ACLService{
		userClient: userClient,
	}
}

func (r ACLService) CheckAccess(
	ctx context.Context,
	uow port.AuthUnitOfWork,
	userUUID uuid.UUID,
	requiredPermissions ...domain.PermissionKeyType,
) (bool, uint64, error) {

	user, err := r.userClient.GetByUUID(ctx, userUUID.String())
	if err != nil {
		return false, 0, err
	}

	roleKeys, err := uow.RoleRepository().GetUserRoleKeys(user.Base.ID)
	if err != nil {
		return false, 0, err
	}

	for _, key := range roleKeys {
		if key == domain.RoleKeySuperAdmin {
			return true, user.Base.ID, nil
		}
	}

	permissionKeys, err := uow.PermissionRepository().GetUserPermissionKeys(user.Base.ID)
	if err != nil {
		return false, 0, err
	}

	for _, requiredPermission := range requiredPermissions {
		if requiredPermission == domain.PermissionKeyNone {
			return true, user.Base.ID, nil
		}
		for _, key := range permissionKeys {
			if requiredPermission == key {
				return true, user.Base.ID, nil
			}
		}
	}

	return false, 0, nil
}

func (r ACLService) AssignUserRoleToUser(uow port.AuthUnitOfWork, userID uint64) error {
	role, err := uow.RoleRepository().GetRoleUser()
	if err != nil {
		return err
	}

	roleIDs := make([]uint64, 1)
	roleIDs[0] = role.Base.ID

	return uow.ACLRepository().AssignRolesToUser(userID, roleIDs)
}
