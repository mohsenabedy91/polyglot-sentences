package tests

import (
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type ACLRepositoryTestSuite struct {
	TestSuite
}

func (r *ACLRepositoryTestSuite) TestACLRepository_AssignRolesToUser_Success() {
	mockLogger := new(logger.MockLogger)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	adminRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})
	userRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "User",
		Key:         "user",
		Description: "User Role",
	})

	repo := authrepository.NewACLRepository(mockLogger, r.GetTx())
	err := repo.AssignRolesToUser(user.Base.ID, []uint64{adminRole.Base.ID, userRole.Base.ID})
	require.NoError(r.T(), err)
}

func (r *ACLRepositoryTestSuite) TestACLRepository_AssignRolesToUser_DeleteError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS access_controls CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewACLRepository(mockLogger, r.GetTx())
	err = repo.AssignRolesToUser(1, []uint64{1, 2, 3})
	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *ACLRepositoryTestSuite) TestACLRepository_AssignRolesToUser_PrepareError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Simulate prepare error by dropping the table
	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS access_controls CASCADE")
	require.NoError(r.T(), err)

	_, err = r.GetTx().Exec("CREATE TABLE access_controls (user_id BIGINT NOT NULL, role_i BIGINT NOT NULL)")
	require.NoError(r.T(), err)

	repo := authrepository.NewACLRepository(mockLogger, r.GetTx())
	err = repo.AssignRolesToUser(1, []uint64{1, 2, 3})
	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *ACLRepositoryTestSuite) TestACLRepository_AssignRolesToUser_InsertError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS access_controls CASCADE")
	require.NoError(r.T(), err)

	_, err = r.GetTx().Exec("CREATE TABLE access_controls (user_id BIGINT NOT NULL, role_id BOOLEAN)")
	require.NoError(r.T(), err)

	repo := authrepository.NewACLRepository(mockLogger, r.GetTx())
	err = repo.AssignRolesToUser(1, []uint64{1, 2, 3})

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}
