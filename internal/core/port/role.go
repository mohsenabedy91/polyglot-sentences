package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type RoleRepository interface {
	Create(ctx context.Context, role domain.Role) error
	GetByUUID(ctx context.Context, uuid uuid.UUID) (*domain.Role, error)
	List(ctx context.Context) ([]*domain.Role, error)
	Update(ctx context.Context, role domain.Role, uuid uuid.UUID) error
	Delete(ctx context.Context, uuid uuid.UUID) error

	ExistKey(ctx context.Context, key string) (bool, error)
	GetRoleUser(ctx context.Context) (domain.Role, error)

	GetPermissions(ctx context.Context, uuid uuid.UUID) (*domain.Role, error)

	GetRoleKeys(ctx context.Context, userID uint64) ([]domain.RoleKeyType, error)
}

type RoleService interface {
	Create(ctx context.Context, role domain.Role) error
	Get(ctx context.Context, uuidStr string) (*domain.Role, error)
	List(ctx context.Context) ([]*domain.Role, error)
	Update(ctx context.Context, role domain.Role, uuidStr string) error
	Delete(ctx context.Context, uuidStr string) error

	GetPermissions(ctx context.Context, uuidStr string) (*domain.Role, error)
}

type RoleCache interface {
	Get(ctx context.Context, key string) (*domain.RoleKeyType, error)
	SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error
}
