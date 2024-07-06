package port

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type RoleRepository interface {
	Create(role domain.Role) error
	GetByUUID(uuid uuid.UUID) (*domain.Role, error)
	List() ([]*domain.Role, error)
	Update(role domain.Role, uuid uuid.UUID) error
	Delete(uuid uuid.UUID, deletedBy uint64) error

	ExistKey(key domain.RoleKeyType) (bool, error)
	GetRoleUser() (domain.Role, error)

	GetPermissions(uuid uuid.UUID) (*domain.Role, error)
	SyncPermissions(roleID uint64, permissionIDs []uint64) error

	GetRoleKeys(userID uint64) ([]domain.RoleKeyType, error)
}

type RoleService interface {
	Create(uow AuthUnitOfWork, role domain.Role) error
	Get(uow AuthUnitOfWork, uuidStr string) (*domain.Role, error)
	List(ctx context.Context, uow AuthUnitOfWork) ([]*domain.Role, error)
	Update(ctx context.Context, uow AuthUnitOfWork, role domain.Role, uuidStr string) error
	Delete(uow AuthUnitOfWork, uuidStr string, deletedBy uint64) error

	GetPermissions(uow AuthUnitOfWork, uuidStr string) (*domain.Role, error)
	SyncPermissions(uow AuthUnitOfWork, uuidStr string, permissionUUIDStr []string) error
}

type RoleCache interface {
	Get(ctx context.Context, key string) (*domain.RoleKeyType, error)
	SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error
}
