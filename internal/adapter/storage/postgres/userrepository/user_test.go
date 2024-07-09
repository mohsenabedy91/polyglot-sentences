package userrepository_test

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	mocklogger "github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func (r *UserRepositoryTestSuite) TestUserRepository_SaveSuccess() {
	mockLogger := new(mocklogger.MockLogger)

	tx, err := r.db.Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)

	var firstName = "John"
	var lastName = "Doe"
	var email = "john.doe@google.com"
	var password = "hashedPassword"

	user := &domain.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     email,
		Password:  &password,
		Status:    domain.UserStatusUnverifiedStr,
	}

	user, err = repo.Save(user)
	require.NoError(r.T(), err)

	require.NoError(r.T(), tx.Commit())

	require.NotNil(r.T(), user.ID)
	require.NotEqual(r.T(), uuid.Nil, user.UUID)

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_SaveInValidStatus() {
	mockLogger := new(mocklogger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.db.Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)

	var firstName = "John"
	var lastName = "Doe"
	var email = "john.doe@google.com"
	var password = "hashedPassword"
	user := &domain.User{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     email,
		Password:  &password,
	}

	user, err = repo.Save(user)
	require.Error(r.T(), err)
	require.Nil(r.T(), user)

	mockLogger.AssertExpectations(r.T())
}
