package permissionservice

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type Service struct {
	log logger.Logger
}

func New(log logger.Logger) *Service {
	return &Service{
		log: log,
	}
}

func (r *Service) List(uow port.AuthUnitOfWork) ([]*domain.Permission, error) {
	return uow.PermissionRepository().List()
}
