package port

import (
	"context"
	"github.com/google/uuid"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
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
	Create(uow repository.UnitOfWork, role domain.Role) error
	Get(uow repository.UnitOfWork, uuidStr string) (*domain.Role, error)
	List(ctx context.Context, uow repository.UnitOfWork) ([]*domain.Role, error)
	Update(ctx context.Context, uow repository.UnitOfWork, role domain.Role, uuidStr string) error
	Delete(uow repository.UnitOfWork, uuidStr string) error

	GetPermissions(uow repository.UnitOfWork, uuidStr string) (*domain.Role, error)
	SyncPermissions(uow repository.UnitOfWork, uuidStr string, permissionUUIDStr []string) error
}

type RoleCache interface {
	Get(ctx context.Context, key string) (*domain.RoleKeyType, error)
	SetBulk(ctx context.Context, items map[string]domain.RoleKeyType) error
}
