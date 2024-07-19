package tests

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type PermissionRepositoryTestSuite struct {
	TestSuite
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_GetUserPermissionKeys_Success() {
	mockLogger := new(logger.MockLogger)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	userRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "User",
		Key:         "user",
		Description: "User Role",
	})
	adminRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Admin Role",
	})

	readUserPermission := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Read User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("read_user")),
		Description: helper.StringPtr("Read User Permission"),
	})
	createUserPermission := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Create User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("create_user")),
		Description: helper.StringPtr("Create User Permission"),
	})

	insertRolePermission(r.T(), r.GetTx(), userRole.Base.ID, readUserPermission.Base.ID)
	insertRolePermission(r.T(), r.GetTx(), userRole.Base.ID, createUserPermission.Base.ID)

	addRoleToUser(r.T(), r.GetTx(), user.Base.ID, adminRole.Base.ID)
	addRoleToUser(r.T(), r.GetTx(), user.Base.ID, userRole.Base.ID)

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	keys, err := repo.GetUserPermissionKeys(user.Base.ID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), keys)
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_GetUserPermissionKeys_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS permissions CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	keys, err := repo.GetUserPermissionKeys(1)

	require.Error(r.T(), err)
	require.Nil(r.T(), keys)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_List_Success() {
	mockLogger := new(logger.MockLogger)

	_, err := r.GetTx().Exec("TRUNCATE permissions CASCADE")
	require.NoError(r.T(), err)

	insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Read User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("read_user")),
		Description: helper.StringPtr("Read User Permission"),
		Group:       helper.StringPtr("user"),
	})
	insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Write User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("write_user")),
		Description: helper.StringPtr("Write User Permission"),
		Group:       helper.StringPtr("user"),
	})

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	permissions, err := repo.List()

	require.NoError(r.T(), err)
	require.NotNil(r.T(), permissions)
	require.Len(r.T(), permissions, 2)
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_List_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS permissions CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	permissions, err := repo.List()

	require.Error(r.T(), err)
	require.Nil(r.T(), permissions)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_FilterValidPermissions_Success() {
	mockLogger := new(logger.MockLogger)

	var permissionUUIDs []uuid.UUID

	permission1 := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Permission 1"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("key 1")),
		Description: helper.StringPtr("Description"),
		Group:       helper.StringPtr("group"),
	})
	permission2 := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Permission 2"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("key 2")),
		Description: helper.StringPtr("Description"),
		Group:       helper.StringPtr("group"),
	})

	permissionUUIDs = append(permissionUUIDs, permission1.UUID)
	permissionUUIDs = append(permissionUUIDs, permission2.UUID)

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	validPermissions, err := repo.FilterValidPermissions(permissionUUIDs)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), validPermissions)
	require.Len(r.T(), validPermissions, 2)
}

func (r *PermissionRepositoryTestSuite) TestPermissionRepository_FilterValidPermissions_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	permissionUUIDs := []uuid.UUID{uuid.New(), uuid.New()}

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS permissions CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewPermissionRepository(mockLogger, r.GetTx())
	validPermissions, err := repo.FilterValidPermissions(permissionUUIDs)

	require.Error(r.T(), err)
	require.Nil(r.T(), validPermissions)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}
