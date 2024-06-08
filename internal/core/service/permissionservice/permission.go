package permissionservice

import (
	"context"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
)

type Service struct {
	log            logger.Logger
	permissionRepo port.PermissionRepository
}

func New(log logger.Logger, permissionRepo port.PermissionRepository) *Service {
	return &Service{
		log:            log,
		permissionRepo: permissionRepo,
	}
}

func (r *Service) List(ctx context.Context) ([]*domain.Permission, error) {
	return r.permissionRepo.List(ctx)
}
