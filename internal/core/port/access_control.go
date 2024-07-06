package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type ACLRepository interface {
	AssignRolesToUser(userID uint64, roleIDs []uint64) error
}

type ACLService interface {
	CheckAccess(
		ctx context.Context,
		uow AuthUnitOfWork,
		userUUID uuid.UUID,
		requiredPermissions ...domain.PermissionKeyType,
	) (bool, uint64, error)
	AssignUserRoleToUser(uow AuthUnitOfWork, userID uint64) error
}
