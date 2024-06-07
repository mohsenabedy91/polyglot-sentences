package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type PermissionRepository interface {
	GetUserPermissionKeys(ctx context.Context, userID uint64) ([]domain.PermissionKeyType, error)
}

type ACLRepository interface {
	AddUserRole(ctx context.Context, userID, roleID uint64) error
}

type ACLService interface {
	CheckAccess(ctx context.Context, userUUID uuid.UUID, permission ...domain.PermissionKeyType) (bool, error)
	AddUserRole(ctx context.Context, userID uint64) error
}
