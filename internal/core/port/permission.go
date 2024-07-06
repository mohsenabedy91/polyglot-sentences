package port

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
)

type PermissionRepository interface {
	GetUserPermissionKeys(userID uint64) ([]domain.PermissionKeyType, error)
	List() ([]*domain.Permission, error)
	FilterValidPermissions(uuids []uuid.UUID) ([]uint64, error)
}

type PermissionService interface {
	List(uow AuthUnitOfWork) ([]*domain.Permission, error)
}
