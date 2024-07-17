package userservice_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_GetByUUID(t *testing.T) {
	mockLog := new(logger.MockLogger)

	userID := uuid.New()
	expectedUser := &domain.User{
		Base: domain.Base{
			UUID: userID,
		},
	}

	t.Run("GetByUUID success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("GetByUUID", userID).Return(expectedUser, nil)

		service := userservice.New(mockLog)
		user, err := service.GetByUUID(mockUow, userID.String())

		require.NoError(t, err)
		require.Equal(t, expectedUser, user)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByUUID repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("GetByUUID", userID).Return(&domain.User{}, serviceerror.NewServerError())

		service := userservice.New(mockLog)
		user, err := service.GetByUUID(mockUow, userID.String())

		require.Error(t, err)
		require.Equal(t, &domain.User{}, user)
		require.IsType(t, &serviceerror.ServiceError{}, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetByID(t *testing.T) {
	mockLog := new(logger.MockLogger)

	userID := uint64(1)
	userUUID := uuid.New()
	expectedUser := &domain.User{
		Base: domain.Base{
			ID:   userID,
			UUID: userUUID,
		},
	}

	t.Run("GetByID success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("GetByID", userID).Return(expectedUser, nil)

		service := userservice.New(mockLog)
		user, err := service.GetByID(mockUow, userID)

		require.NoError(t, err)
		require.Equal(t, expectedUser, user)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByID repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("GetByID", userID).Return(&domain.User{}, serviceerror.NewServerError())

		service := userservice.New(mockLog)
		user, err := service.GetByID(mockUow, userID)

		require.Error(t, err)
		require.Equal(t, &domain.User{}, user)
		require.IsType(t, &serviceerror.ServiceError{}, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_IsEmailUnique(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	email := "test@example.com"

	t.Run("IsEmailUnique success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("IsEmailUnique", email).Return(true, nil)

		service := userservice.New(mockLogger)
		err := service.IsEmailUnique(mockUow, email)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Email is not unique", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("IsEmailUnique", email).Return(false, nil)

		service := userservice.New(mockLogger)
		err := service.IsEmailUnique(mockUow, email)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.EmailRegistered, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("IsEmailUnique repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("IsEmailUnique", email).Return(false, serviceerror.NewServerError())

		mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

		service := userservice.New(mockLogger)
		err := service.IsEmailUnique(mockUow, email)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})
}

func TestUserService_GetByEmail(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	email := "test@example.com"

	t.Run("GetByEmail success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		expectedUser := &domain.User{Email: email}
		mockRepo.On("GetByEmail", email).Return(expectedUser, nil)

		service := userservice.New(mockLogger)
		user, err := service.GetByEmail(mockUow, email)

		require.NoError(t, err)
		require.Equal(t, expectedUser, user)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("GetByEmail repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("GetByEmail", email).Return(&domain.User{}, serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		user, err := service.GetByEmail(mockUow, email)

		require.Error(t, err)
		require.Equal(t, &domain.User{}, user)
		require.IsType(t, &serviceerror.ServiceError{}, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_List(t *testing.T) {
	mockLogger := new(logger.MockLogger)

	expectedUsers := []*domain.User{
		{Email: "user1@example.com"},
		{Email: "user2@example.com"},
	}

	t.Run("List success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("List").Return(expectedUsers, nil)

		service := userservice.New(mockLogger)
		users, err := service.List(mockUow)

		require.NoError(t, err)
		require.Equal(t, expectedUsers, users)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("List repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("List").Return([]*domain.User{}, serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		users, err := service.List(mockUow)

		require.Error(t, err)
		require.Equal(t, []*domain.User{}, users)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Create(t *testing.T) {
	mockLogger := new(logger.MockLogger)

	userID := uuid.New()
	expectedUser := &domain.User{
		Base: domain.Base{
			UUID: userID,
		},
		Email:  "test@example.com",
		Status: domain.UserStatusUnverified,
	}

	newUser := domain.User{Email: "test@example.com"}
	newUser.Status = domain.UserStatusUnverified

	t.Run("Create success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("Save", &newUser).Return(expectedUser, nil)

		service := userservice.New(mockLogger)
		user, err := service.Create(mockUow, newUser)

		require.NoError(t, err)
		require.Equal(t, expectedUser.Email, user.Email)
		require.Equal(t, expectedUser.Status, user.Status)
		require.Equal(t, userID, user.Base.UUID)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("Save", &newUser).Return(&domain.User{}, serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		user, err := service.Create(mockUow, newUser)

		require.Error(t, err)
		require.Equal(t, &domain.User{}, user)
		require.IsType(t, &serviceerror.ServiceError{}, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_VerifiedEmail(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	email := "test@example.com"

	t.Run("VerifiedEmail success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("VerifiedEmail", email).Return(nil)

		service := userservice.New(mockLogger)
		err := service.VerifiedEmail(mockUow, email)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("VerifiedEmail repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("VerifiedEmail", email).Return(serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		err := service.VerifiedEmail(mockUow, email)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_MarkWelcomeMessageSent(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	id := uint64(1)

	t.Run("MarkWelcomeMessageSent success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)
		mockUow := new(userrepository.MockUnitOfWork)

		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("MarkWelcomeMessageSent", id).Return(nil)

		service := userservice.New(mockLogger)
		err := service.MarkWelcomeMessageSent(mockUow, id)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MarkWelcomeMessageSent repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("MarkWelcomeMessageSent", id).Return(serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		err := service.MarkWelcomeMessageSent(mockUow, id)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateGoogleID(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	id := uint64(1)
	googleID := "google123"

	t.Run("UpdateGoogleID success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)
		mockUow := new(userrepository.MockUnitOfWork)

		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdateGoogleID", id, googleID).Return(nil)

		service := userservice.New(mockLogger)
		err := service.UpdateGoogleID(mockUow, id, googleID)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateGoogleID repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)
		mockUow := new(userrepository.MockUnitOfWork)

		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdateGoogleID", id, googleID).Return(serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		err := service.UpdateGoogleID(mockUow, id, googleID)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateLastLoginTime(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	id := uint64(1)

	t.Run("UpdateLastLoginTime success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)
		mockUow := new(userrepository.MockUnitOfWork)

		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdateLastLoginTime", id).Return(nil)

		service := userservice.New(mockLogger)
		err := service.UpdateLastLoginTime(mockUow, id)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdateLastLoginTime repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdateLastLoginTime", id).Return(serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		err := service.UpdateLastLoginTime(mockUow, id)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdatePassword(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	id := uint64(1)
	password := "new-password123"

	t.Run("UpdatePassword success", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)
		mockUow := new(userrepository.MockUnitOfWork)

		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdatePassword", id, password).Return(nil)

		service := userservice.New(mockLogger)
		err := service.UpdatePassword(mockUow, id, password)

		require.NoError(t, err)

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("UpdatePassword repository error", func(t *testing.T) {
		mockRepo := new(userrepository.MockUserRepository)

		mockUow := new(userrepository.MockUnitOfWork)
		mockUow.On("UserRepository").Return(mockRepo)

		mockRepo.On("UpdatePassword", id, password).Return(serviceerror.NewServerError())

		service := userservice.New(mockLogger)
		err := service.UpdatePassword(mockUow, id, password)

		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
