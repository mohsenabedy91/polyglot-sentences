package permissionservice

import (
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
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

func (r *Service) List(uow repository.UnitOfWork) ([]*domain.Permission, error) {
	return uow.PermissionRepository().List()
}
