package userservice

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

// UserService implements port.UserService interface and provides access to the user repository and cache service
type UserService struct {
	log logger.Logger
}

// New creates a new user service instance
func New(log logger.Logger) *UserService {
	return &UserService{
		log: log,
	}
}

func (r *UserService) GetByUUID(uow port.UserUnitOfWork, uuidStr string) (user *domain.User, err error) {
	return uow.UserRepository().GetByUUID(uuid.MustParse(uuidStr))
}

func (r *UserService) GetByID(uow port.UserUnitOfWork, id uint64) (user *domain.User, err error) {
	return uow.UserRepository().GetByID(id)
}

func (r *UserService) IsEmailUnique(uow port.UserUnitOfWork, email string) error {
	isUniqueEmail, err := uow.UserRepository().IsEmailUnique(email)
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

func (r *UserService) GetByEmail(uow port.UserUnitOfWork, email string) (*domain.User, error) {
	return uow.UserRepository().GetByEmail(email)
}

func (r *UserService) List(uow port.UserUnitOfWork) ([]*domain.User, error) {
	return uow.UserRepository().List()
}

func (r *UserService) Create(uow port.UserUnitOfWork, user domain.User) (*domain.User, error) {
	if user.Status != domain.UserStatusActive {
		user.Status = domain.UserStatusUnverified
	}
	return uow.UserRepository().Save(&user)
}

func (r *UserService) VerifiedEmail(uow port.UserUnitOfWork, email string) error {
	return uow.UserRepository().VerifiedEmail(email)
}

func (r *UserService) MarkWelcomeMessageSent(uow port.UserUnitOfWork, id uint64) error {
	return uow.UserRepository().MarkWelcomeMessageSent(id)
}

func (r *UserService) UpdateGoogleID(uow port.UserUnitOfWork, id uint64, googleID string) error {
	return uow.UserRepository().UpdateGoogleID(id, googleID)
}

func (r *UserService) UpdateLastLoginTime(uow port.UserUnitOfWork, id uint64) error {
	return uow.UserRepository().UpdateLastLoginTime(id)
}

func (r *UserService) UpdatePassword(uow port.UserUnitOfWork, id uint64, password string) error {
	return uow.UserRepository().UpdatePassword(id, password)
}
