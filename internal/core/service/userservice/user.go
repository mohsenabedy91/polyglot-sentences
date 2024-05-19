package userservice

import (
	"context"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

// UserService implements port.UserService interface and provides access to the user repository and cache service
type UserService struct {
	log      logger.Logger
	userRepo port.UserRepository
}

// New creates a new user service instance
func New(log logger.Logger, userRepo port.UserRepository) *UserService {
	return &UserService{
		log:      log,
		userRepo: userRepo,
	}
}

func (r UserService) RegisterUser(ctx context.Context, user domain.User) serviceerror.Error {
	return r.userRepo.Save(ctx, &user)
}

func (r UserService) GetUser(ctx context.Context, uuid uuid.UUID) (user *domain.User, err serviceerror.Error) {
	return r.userRepo.GetByUUID(ctx, uuid)
}

func (r UserService) IsEmailUnique(ctx context.Context, email string) serviceerror.Error {
	isUniqueEmail, err := r.userRepo.IsEmailUnique(ctx, email)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.String(), nil)
		return err
	}

	if !isUniqueEmail {
		return serviceerror.NewServiceError(
			serviceerror.EmailRegistered,
			map[string]interface{}{
				"email": email,
			},
		)
	}

	return nil
}

func (r UserService) GetByEmail(ctx context.Context, email string) (*domain.User, serviceerror.Error) {
	return r.userRepo.GetByEmail(ctx, email)
}
