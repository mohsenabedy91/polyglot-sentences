package authorizationservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

type AccessControlService struct {
	log            logger.Logger
	userRepo       port.UserRepository
	permissionRepo port.PermissionRepository
}

func New(
	log logger.Logger,
	userRepo port.UserRepository,
	permissionRepo port.PermissionRepository,
) *AccessControlService {
	return &AccessControlService{
		log:            log,
		userRepo:       userRepo,
		permissionRepo: permissionRepo,
	}
}

func (r AccessControlService) CheckAccess(
	ctx context.Context,
	userUUID uuid.UUID,
	permissions ...domain.PermissionKeyType,
) (bool, serviceerror.Error) {

	user, err := r.userRepo.GetByUUID(ctx, userUUID)
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.String(), nil)
		return false, err
	}

	permissionKeys, err := r.permissionRepo.GetUserPermissionKeys(ctx, user.ID)
	if err != nil {
		r.log.Error(logger.Authorization, logger.DatabaseSelect, err.String(), nil)
		return false, err
	}

	for _, key := range permissionKeys {
		for _, permission := range permissions {
			if permission == key {
				return true, nil
			}
		}
	}

	return false, nil
}
