package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type ACLRepository interface {
	AssignRolesToUser(ctx context.Context, userID uint64, roleIDs []uint64) error
}

type ACLService interface {
	CheckAccess(ctx context.Context, userUUID uuid.UUID, permission ...domain.PermissionKeyType) (bool, error)
	AssignUserRoleToUser(ctx context.Context, userID uint64) error
}
