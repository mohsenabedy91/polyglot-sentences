package permissionservice

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (r *Service) List(uow port.AuthUnitOfWork) ([]*domain.Permission, error) {
	return uow.PermissionRepository().List()
}
