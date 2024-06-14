package port

import (
	"github.com/google/uuid"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type PermissionRepository interface {
	GetUserPermissionKeys(ctx context.Context, userID uint64) ([]domain.PermissionKeyType, error)
	List(ctx context.Context) ([]*domain.Permission, error)
}

type PermissionService interface {
	List(uow repository.UnitOfWork) ([]*domain.Permission, error)
}
