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

func (r UserService) GetByUUID(ctx context.Context, uuidStr string) (user *domain.User, err error) {
	return r.userRepo.GetByUUID(ctx, uuid.MustParse(uuidStr))
}

func (r UserService) IsEmailUnique(ctx context.Context, email string) error {
	isUniqueEmail, err := r.userRepo.IsEmailUnique(ctx, email)
	if err != nil {
		r.log.Error(logger.Database, logger.DatabaseSelect, err.Error(), nil)
		return err
	}

	if !isUniqueEmail {
		return serviceerror.New(
			serviceerror.EmailRegistered,
			map[string]interface{}{
				"email": email,
			},
		)
	}

	return nil
}

func (r UserService) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.userRepo.GetByEmail(ctx, email)
}

func (r UserService) List(ctx context.Context) ([]domain.User, error) {
	return r.userRepo.List(ctx)
}

func (r UserService) Create(ctx context.Context, user domain.User) error {
	return r.userRepo.Save(ctx, &user)
}

func (r UserService) VerifiedEmail(ctx context.Context, email string) error {
	return r.userRepo.VerifiedEmail(ctx, email)
}

func (r UserService) UpdateWelcomeMessageToSentFlag(ctx context.Context, id uint64) error {
	return r.userRepo.UpdateWelcomeMessageToSentFlag(ctx, id)
}
