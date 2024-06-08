package port

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type PermissionRepository interface {
	GetUserPermissionKeys(ctx context.Context, userID uint64) ([]domain.PermissionKeyType, error)
	List(ctx context.Context) ([]*domain.Permission, error)
}

type PermissionService interface {
	List(ctx context.Context) ([]*domain.Permission, error)
}
