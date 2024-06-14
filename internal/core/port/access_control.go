package port

import (
	"context"
	"github.com/google/uuid"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type ACLRepository interface {
	AssignRolesToUser(userID uint64, roleIDs []uint64) error
}

type ACLService interface {
	CheckAccess(
		ctx context.Context,
		uow repository.UnitOfWork,
		userUUID uuid.UUID,
		requiredPermissions ...domain.PermissionKeyType,
	) (bool, uint64, error)
	AssignUserRoleToUser(uow repository.UnitOfWork, userID uint64) error
}
