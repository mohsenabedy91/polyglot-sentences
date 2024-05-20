package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type PermissionRepository interface {
	GetUserPermissionKeys(ctx context.Context, userID uint) ([]domain.PermissionKeyType, error)
}

type AccessControlService interface {
	CheckAccess(ctx context.Context, userUUID uuid.UUID, permission ...domain.PermissionKeyType) (bool, error)
}
