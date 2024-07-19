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

type RoleRepositoryTestSuite struct {
	TestSuite
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Create_Success() {
	mockLogger := new(logger.MockLogger)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.Create(domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
		Modifier: domain.Modifier{
			CreatedBy: &user.Base.ID,
		},
	})
	require.NoError(r.T(), err)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Create_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err = repo.Create(domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
		Modifier: domain.Modifier{
			CreatedBy: new(uint64),
		},
	})

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Create_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	// Create the trigger to prevent inserts
	_, err := r.GetTx().Exec(`
		CREATE OR REPLACE FUNCTION prevent_insert() RETURNS TRIGGER AS $$
		BEGIN
		    RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;

		CREATE TRIGGER prevent_insert_trigger
		BEFORE INSERT ON roles
		FOR EACH ROW
		EXECUTE FUNCTION prevent_insert();
	`)
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err = repo.Create(domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	// Clean up: Drop the trigger and function
	_, err = r.GetTx().Exec(`
		DROP TRIGGER IF EXISTS prevent_insert_trigger ON roles;
		DROP FUNCTION IF EXISTS prevent_insert();
	`)
	require.NoError(r.T(), err)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetByUUID_Success() {
	mockLogger := new(logger.MockLogger)

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetByUUID(newRole.Base.UUID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), fetchedRole)
	require.Equal(r.T(), newRole.Base.UUID, fetchedRole.Base.UUID)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetByUUID_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetByUUID(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedRole)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetByUUID_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE;")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetByUUID(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedRole)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_List_Success() {
	mockLogger := new(logger.MockLogger)

	_, err := r.GetTx().Exec("TRUNCATE roles CASCADE")
	require.NoError(r.T(), err)

	insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})
	insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "User",
		Key:         "user",
		Description: "User Role",
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	roles, err := repo.List()

	require.NoError(r.T(), err)
	require.NotNil(r.T(), roles)
	require.Len(r.T(), roles, 2)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_List_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	roles, err := repo.List()

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), roles)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Update_Success() {
	mockLogger := new(logger.MockLogger)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
		Modifier: domain.Modifier{
			UpdatedBy: user.Base.ID,
		},
	})

	newRole.Title = "Super Admin"
	newRole.Description = "Super Administrator Role"

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.Update(*newRole, newRole.Base.UUID)

	require.NoError(r.T(), err)

	updatedRole, err := repo.GetByUUID(newRole.Base.UUID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), updatedRole)
	require.Equal(r.T(), "Super Admin", updatedRole.Title)
	require.Equal(r.T(), "Super Administrator Role", updatedRole.Description)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Update_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	role := domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	}

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err = repo.Update(role, uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Update_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	role := domain.Role{
		Title:       "Nonexistent Role",
		Key:         "nonexistent_key",
		Description: "This role does not exist",
	}

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.Update(role, uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.NoRowsEffected, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Delete_Success() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Warn", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "User",
		Key:         "user",
		Description: "User Role",
		IsDefault:   false,
		Modifier: domain.Modifier{
			CreatedBy: &user.Base.ID,
		},
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.Delete(newRole.Base.UUID, user.Base.ID)

	require.NoError(r.T(), err)

	deletedRole, err := repo.GetByUUID(newRole.Base.UUID)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), deletedRole)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Delete_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err = repo.Delete(uuid.New(), 1)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_Delete_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.Delete(uuid.New(), 100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.IsNotDeletable, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_ExistKey_Success() {
	mockLogger := new(logger.MockLogger)

	insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	exists, err := repo.ExistKey("admin")

	require.NoError(r.T(), err)
	require.True(r.T(), exists)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_ExistKey_NotExist() {
	mockLogger := new(logger.MockLogger)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	exists, err := repo.ExistKey("nonexistent")

	require.NoError(r.T(), err)
	require.False(r.T(), exists)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_ExistKey_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	exists, err := repo.ExistKey("admin")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.True(r.T(), exists)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetRoleUser_Success() {
	mockLogger := new(logger.MockLogger)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetRoleUser()

	require.NoError(r.T(), err)
	require.NotNil(r.T(), fetchedRole)
	require.Equal(r.T(), domain.RoleKeyUser, fetchedRole.Key)

}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetRoleUser_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("TRUNCATE roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetRoleUser()

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Equal(r.T(), uuid.Nil, fetchedRole.Base.UUID)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetRoleUser_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	fetchedRole, err := repo.GetRoleUser()

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Equal(r.T(), uuid.Nil, fetchedRole.Base.UUID)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetPermissions_Success() {
	mockLogger := new(logger.MockLogger)

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
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
	deleteUserPermission := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Delete_User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("delete_user")),
		Description: helper.StringPtr("Delete User Permission"),
	})

	insertRolePermission(r.T(), r.GetTx(), newRole.Base.ID, readUserPermission.Base.ID)
	insertRolePermission(r.T(), r.GetTx(), newRole.Base.ID, createUserPermission.Base.ID)
	insertRolePermission(r.T(), r.GetTx(), newRole.Base.ID, deleteUserPermission.Base.ID)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	role, err := repo.GetPermissions(newRole.Base.UUID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), role)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetPermissions_RoleWithoutPermission_Success() {
	mockLogger := new(logger.MockLogger)

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	role, err := repo.GetPermissions(newRole.Base.UUID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), role)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetPermissions_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	role, err := repo.GetPermissions(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), role)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetPermissions_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	role, err := repo.GetPermissions(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), role)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetUserRoleKeys_Success() {
	mockLogger := new(logger.MockLogger)

	user := insertUser(r.T(), r.GetTx(), &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	addRoleToUser(r.T(), r.GetTx(), user.Base.ID, newRole.Base.ID)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	keys, err := repo.GetUserRoleKeys(user.Base.ID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), keys)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_GetUserRoleKeys_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS roles CASCADE")
	require.NoError(r.T(), err)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	keys, err := repo.GetUserRoleKeys(1)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), keys)

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_SyncPermissions_Success() {
	mockLogger := new(logger.MockLogger)

	newRole := insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
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
	deleteUserPermission := insertPermission(r.T(), r.GetTx(), &domain.Permission{
		Title:       helper.StringPtr("Delete_User"),
		Key:         (*domain.PermissionKeyType)(helper.StringPtr("delete_user")),
		Description: helper.StringPtr("Delete User Permission"),
	})

	permissionIDs := []uint64{readUserPermission.Base.ID, createUserPermission.Base.ID, deleteUserPermission.Base.ID}

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.SyncPermissions(newRole.Base.ID, permissionIDs)

	require.NoError(r.T(), err)
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_SyncPermissions_DBError_ExecDeleteRolePermissions() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	_, err := r.GetTx().Exec("DROP TABLE IF EXISTS role_permissions CASCADE")
	require.NoError(r.T(), err)

	insertRole(r.T(), r.GetTx(), &domain.Role{
		Title:       "Admin",
		Key:         "admin",
		Description: "Administrator Role",
	})

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err = repo.SyncPermissions(100_000, []uint64{1, 2, 3})

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}

func (r *RoleRepositoryTestSuite) TestRoleRepository_SyncPermissions_DBError_ExecInsertRolePermissions() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	repo := authrepository.NewRoleRepository(mockLogger, r.GetTx())
	err := repo.SyncPermissions(100_000, []uint64{1, 2, 3})

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	mockLogger.AssertExpectations(r.T())
}
