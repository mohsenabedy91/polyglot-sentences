package tests

import (
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type UserRepositoryTestSuite struct {
	TestSuite
}

func (r *UserRepositoryTestSuite) TestUserRepository_SaveSuccess() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	user, err := repo.Save(&domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusUnverifiedStr,
	})

	require.NoError(r.T(), err)
	require.NotNil(r.T(), user.Base.ID)
	require.NotEqual(r.T(), uuid.Nil, user.Base.UUID)

	require.NoError(r.T(), tx.Commit())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_SaveInValidStatus() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	user, err := repo.Save(&domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
	})

	require.Error(r.T(), err)
	require.Nil(r.T(), user)

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByUUID_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByUUID(user.Base.UUID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), fetchedUser)
	require.Equal(r.T(), user.Base.UUID, fetchedUser.Base.UUID)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByUUID_UserInActive() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Warn", logger.Database, logger.DatabaseSelect, "The User is inactive", mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusUnverifiedStr,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByUUID(user.Base.UUID)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.UserInActive, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByUUID_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByUUID(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByUUID_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByUUID(uuid.New())

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_IsEmailUnique_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	isUnique, err := repo.IsEmailUnique("unique.email@example.com")

	require.NoError(r.T(), err)
	require.True(r.T(), isUnique)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_IsEmailUnique_NotUnique() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("Jane"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "existing.email@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	isUnique, err := repo.IsEmailUnique(user.Email)

	require.NoError(r.T(), err)
	require.False(r.T(), isUnique)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_IsEmailUnique_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	isUnique, err := repo.IsEmailUnique("email@example.com")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.False(r.T(), isUnique)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByID_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByID(user.Base.ID)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), fetchedUser)
	require.Equal(r.T(), user.Base.ID, fetchedUser.ID)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByID_UserInActive() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Warn", logger.Database, logger.DatabaseSelect, "The User is inactive", mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusUnverifiedStr,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByID(user.Base.ID)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.UserInActive, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByID_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByID(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.RecordNotFound, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByID_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByID(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByEmail_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByEmail(user.Email)

	require.NoError(r.T(), err)
	require.NotNil(r.T(), fetchedUser)
	require.Equal(r.T(), user.Email, fetchedUser.Email)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByEmail_RecordNotFound() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Warn", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByEmail("notfound.email@example.com")

	require.NoError(r.T(), err)
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Commit())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_GetByEmail_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	fetchedUser, err := repo.GetByEmail("email@example.com")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), fetchedUser)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_List_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("TRUNCATE users CASCADE;")
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("Jane"),
		LastName:  helper.StringPtr("Smith"),
		Email:     "jane.smith@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	users, err := repo.List()

	require.NoError(r.T(), err)
	require.NotNil(r.T(), users)
	require.Len(r.T(), users, 2)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_List_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseSelect, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	users, err := repo.List()

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())
	require.Nil(r.T(), users)

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_VerifiedEmail_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusUnverifiedStr,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.VerifiedEmail(user.Email)

	require.NoError(r.T(), err)

	fetchedUser, err := repo.GetByUUID(user.Base.UUID)

	require.NoError(r.T(), err)
	require.Equal(r.T(), domain.UserStatusActive, fetchedUser.Status)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_VerifiedEmail_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseUpdate, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.VerifiedEmail("email@example.com")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_VerifiedEmail_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.VerifiedEmail("invalid@example.com")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_MarkWelcomeMessageSent_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.MarkWelcomeMessageSent(user.Base.ID)

	require.NoError(r.T(), err)

	fetchedUser, err := repo.GetByEmail(user.Email)

	require.NoError(r.T(), err)
	require.True(r.T(), fetchedUser.WelcomeMessageSent)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_MarkWelcomeMessageSent_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseUpdate, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.MarkWelcomeMessageSent(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_MarkWelcomeMessageSent_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.MarkWelcomeMessageSent(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateGoogleID_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	googleID := "new-google-id"

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateGoogleID(user.Base.ID, googleID)

	require.NoError(r.T(), err)

	fetchedUser, err := repo.GetByEmail(user.Email)

	require.NoError(r.T(), err)
	require.Equal(r.T(), googleID, *fetchedUser.GoogleID)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateGoogleID_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseUpdate, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateGoogleID(100_000, "new-google-id")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateGoogleID_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateGoogleID(100_000, "new-google-id")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateLastLoginTime_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateLastLoginTime(user.Base.ID)

	require.NoError(r.T(), err)

	// Since we cannot accurately test the timestamp, we just check if the operation succeeded.
	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateLastLoginTime_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseUpdate, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateLastLoginTime(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdateLastLoginTime_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdateLastLoginTime(100_000)

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdatePassword_Success() {
	mockLogger := new(logger.MockLogger)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	user := insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	newPassword := "newHashedPassword"

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdatePassword(user.Base.ID, newPassword)

	require.NoError(r.T(), err)

	fetchedUser, err := repo.GetByEmail(user.Email)

	require.NoError(r.T(), err)
	require.Equal(r.T(), newPassword, *fetchedUser.Password)

	require.NoError(r.T(), tx.Commit())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdatePassword_DBError() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", logger.Database, logger.DatabaseUpdate, mock.Anything, mock.Anything).Return()

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	_, err = tx.Exec("DROP TABLE IF EXISTS users CASCADE")
	require.NoError(r.T(), err)

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdatePassword(100_000, "newHashedPassword")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}

func (r *UserRepositoryTestSuite) TestUserRepository_UpdatePassword_NoRowsAffected() {
	mockLogger := new(logger.MockLogger)
	mockLogger.On("Error", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	tx, err := r.GetDB().Begin()
	require.NoError(r.T(), err)

	insertUser(r.T(), tx, &domain.User{
		FirstName: helper.StringPtr("John"),
		LastName:  helper.StringPtr("Doe"),
		Email:     "john.doe@example.com",
		Password:  helper.StringPtr("hashedPassword"),
		Status:    domain.UserStatusActive,
	})

	repo := userrepository.NewUserRepository(mockLogger, tx)
	err = repo.UpdatePassword(100_000, "newHashedPassword")

	require.Error(r.T(), err)
	require.Equal(r.T(), serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

	require.NoError(r.T(), tx.Rollback())

	mockLogger.AssertExpectations(r.T())
}
