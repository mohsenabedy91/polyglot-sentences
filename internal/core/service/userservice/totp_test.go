package userservice_test

import (
	"context"
	"errors"
	"fmt"
	repository "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/postgres/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis/userrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/userservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func generateTOTPKey(t *testing.T, conf config.Config, email string) *otp.Key {
	t.Helper()

	otpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      conf.App.Name,
		AccountName: email,
		Digits:      otp.DigitsSix,
	})
	require.NoError(t, err, "failed to generate TOTP key")
	return otpKey
}

func generateValidCode(t *testing.T, otpKey *otp.Key) string {
	t.Helper()

	code, err := totp.GenerateCode(otpKey.Secret(), time.Now().UTC())
	require.NoError(t, err, "failed to generate valid TOTP code")
	return code
}

func TestTOTPService_Enroll(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Config{
		App: config.App{
			Name: "TestApp",
		},
		Auth: config.Auth{
			EncryptionKey: "encryption-key",
		},
	}

	email := "test@example.com"
	ctx := context.Background()

	t.Run("Enroll success", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Set", ctx, email, mock.AnythingOfType("string")).Return(nil)
		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)

		totpKey, err := service.Enroll(ctx, email)
		require.NoError(t, err, "expected no error on successful enroll")
		require.NotNil(t, totpKey, "expected non-nil TOTP key")
		require.NotEmpty(t, totpKey.Secret, "expected non-empty TOTP secret")
		require.NotEmpty(t, totpKey.URL, "expected non-empty TOTP URL")

		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enroll TOTP generation error", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockLogger.On("Error", logger.TOTP, logger.EnrollTOTP, mock.Anything, mock.Anything).Return()

		invalidConf := conf
		invalidConf.App.Name = ""
		service := userservice.NewTOTPService(mockLogger, invalidConf, mockTOTPCache)

		totpKey, err := service.Enroll(ctx, email)
		require.Error(t, err, "expected an error due to invalid configuration")
		require.Nil(t, totpKey, "expected nil TOTP key on failure")

		mockLogger.AssertExpectations(t)
	})

	t.Run("Enroll cache error", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Set", ctx, email, mock.AnythingOfType("string")).Return(errors.New("cache error"))
		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)

		totpKey, err := service.Enroll(ctx, email)
		require.Error(t, err, "expected an error due to cache failure")
		require.Nil(t, totpKey, "expected nil TOTP key on cache failure")

		mockTOTPCache.AssertExpectations(t)
	})
}

func TestTOTPService_Enable(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Config{
		App: config.App{
			Name: "TestApp",
		},
		Auth: config.Auth{
			EncryptionKey: "encryption-key",
		},
	}

	ctx := context.Background()
	userID := uint64(1)
	email := "test@example.com"
	otpKey := generateTOTPKey(t, conf, email)
	validCode := generateValidCode(t, otpKey)

	t.Run("Enable success", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("UpdateTOTPSecret", userID, &encryptedSecret).Return(nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		require.NoError(t, service.Enable(ctx, mockUow, userID, email, validCode), "expected no error on successful TOTP enable")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enable cache get error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Get", ctx, email).Return("", errors.New("cache error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Enable(ctx, mockUow, userID, email, validCode)
		require.Error(t, err, "expected an error due to cache get failure")

		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enable verification error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))
		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		invalidCode := "654321"
		mockLogger.On("Warn", logger.TOTP, logger.VerifyTOTP, fmt.Sprintf("The code «%s» is not valid ", invalidCode), mock.Anything).Return()

		err := service.Enable(ctx, mockUow, userID, email, invalidCode)
		require.Error(t, err, "expected an error due to invalid TOTP code")
		require.IsType(t, &serviceerror.ServiceError{}, err, "expected ServiceError type")
		require.Equal(t, serviceerror.InvalidTOTPCode, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockTOTPCache.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("Enable repository error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("UpdateTOTPSecret", userID, &encryptedSecret).Return(errors.New("repository error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Enable(ctx, mockUow, userID, email, validCode)
		require.Error(t, err, "expected an error due to repository failure")
		require.Equal(t, "repository error", err.Error(), "expected specific repository error message")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTOTPCache.AssertExpectations(t)
	})
}

func TestTOTPService_Disable(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Config{
		App: config.App{
			Name: "TestApp",
		},
		Auth: config.Auth{
			EncryptionKey: "encryption-key",
		},
	}

	ctx := context.Background()
	userID := uint64(1)

	email := "test@example.com"
	otpKey := generateTOTPKey(t, conf, email)
	validCode := generateValidCode(t, otpKey)

	t.Run("Disable success", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))
		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)

		// Enable TOTP first
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("UpdateTOTPSecret", userID, &encryptedSecret).Return(nil)
		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		require.NoError(t, service.Enable(ctx, mockUow, userID, email, validCode), "expected no error on successful TOTP enable")

		// Now disable TOTP
		mockUow = new(repository.MockUnitOfWork)
		mockRepo = new(repository.MockUserRepository)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)
		mockRepo.On("UpdateTOTPSecret", userID, (*string)(nil)).Return(nil)
		require.NoError(t, service.Disable(ctx, mockUow, userID, validCode), "expected no error on successful TOTP disable")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Disable repository GetTOTPSecret error", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return((*string)(nil), serviceerror.NewServerError())

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Disable(ctx, mockUow, userID, validCode)
		require.Error(t, err, "expected an error due to repository GetTOTPSecret failure")
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Disable repository GetTOTPSecret return nil", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return((*string)(nil), nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Disable(ctx, mockUow, userID, validCode)
		require.Error(t, err, "expected an error due to TOTP not being enrolled")
		require.Equal(t, serviceerror.TOTPNotEnrolled, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Disable verification error", func(t *testing.T) {
		t.Parallel()

		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		invalidCode := "654321"
		mockLogger.On("Warn", logger.TOTP, logger.VerifyTOTP, fmt.Sprintf("The code «%s» is not valid ", invalidCode), mock.Anything).Return()

		err := service.Disable(ctx, mockUow, userID, invalidCode)
		require.Error(t, err, "expected an error due to invalid TOTP code")
		require.IsType(t, &serviceerror.ServiceError{}, err, "expected ServiceError type")
		require.Equal(t, serviceerror.InvalidTOTPCode, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("Disable repository UpdateTOTPSecret error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)
		mockRepo.On("UpdateTOTPSecret", userID, (*string)(nil)).Return(errors.New("repository error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Disable(ctx, mockUow, userID, validCode)
		require.Error(t, err, "expected an error due to repository UpdateTOTPSecret failure")
		require.Equal(t, "repository error", err.Error(), "expected specific repository error message")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}

func TestTOTPService_Verify(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Config{
		App: config.App{
			Name: "TestApp",
		},
	}
	service := userservice.NewTOTPService(mockLogger, conf, nil)

	email := "test@example.com"
	otpKey := generateTOTPKey(t, conf, email)

	t.Run("Verify success", func(t *testing.T) {
		t.Parallel()

		validCode := generateValidCode(t, otpKey)

		valid, err := service.Verify(validCode, otpKey.Secret())
		require.NoError(t, err, "expected no error on successful verification")
		require.True(t, valid, "expected verification to succeed")
	})

	t.Run("Verify invalid code", func(t *testing.T) {
		t.Parallel()

		invalidCode := "123456"
		mockLogger.On("Warn", logger.TOTP, logger.VerifyTOTP, fmt.Sprintf("The code «%s» is not valid ", invalidCode), mock.Anything).Return()

		valid, err := service.Verify(invalidCode, otpKey.Secret())
		require.Error(t, err, "expected an error due to invalid TOTP code")
		require.False(t, valid, "expected verification to fail")
		require.IsType(t, &serviceerror.ServiceError{}, err, "expected ServiceError type")
		require.Equal(t, serviceerror.InvalidTOTPCode, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockLogger.AssertExpectations(t)
	})
}

func TestTOTPService_Get(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Config{
		App: config.App{
			Name: "TestApp",
		},
		Auth: config.Auth{
			EncryptionKey: "encryption-key",
		},
	}

	ctx := context.Background()
	userID := uint64(1)
	email := "test@example.com"

	otpKey := generateTOTPKey(t, conf, email)

	t.Run("Get success", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		service := userservice.NewTOTPService(mockLogger, conf, nil)

		key := helper.ToAESKey(conf.Auth.EncryptionKey)
		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), key)

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)

		totpKey, err := service.Get(ctx, mockUow, userID, email)
		require.NoError(t, err, "expected no error on successful Get")
		require.NotNil(t, totpKey, "expected non-nil TOTP key")
		require.NotEmpty(t, totpKey.Secret, "expected non-empty TOTP secret")
		require.NotEmpty(t, totpKey.URL, "expected non-empty TOTP URL")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get repository error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		service := userservice.NewTOTPService(mockLogger, conf, nil)

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return((*string)(nil), errors.New("repository error"))

		totpKey, err := service.Get(ctx, mockUow, userID, email)
		require.Error(t, err, "expected an error due to repository failure")
		require.Nil(t, totpKey, "expected nil TOTP key on repository failure")
		require.Equal(t, "repository error", err.Error(), "expected specific repository error message")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get TOTP not enrolled", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		service := userservice.NewTOTPService(mockLogger, conf, nil)

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return((*string)(nil), nil)

		totpKey, err := service.Get(ctx, mockUow, userID, email)
		require.Error(t, err, "expected an error due to TOTP not being enrolled")
		require.Nil(t, totpKey, "expected nil TOTP key when not enrolled")
		require.IsType(t, &serviceerror.ServiceError{}, err, "expected ServiceError type")
		require.Equal(t, serviceerror.TOTPNotEnrolled, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get base32 decode error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		service := userservice.NewTOTPService(mockLogger, conf, nil)

		invalidBase32Secret := "invalid-base32-secret"
		key := helper.ToAESKey(conf.Auth.EncryptionKey)
		encryptedSecret, _ := helper.EncryptSecret([]byte(invalidBase32Secret), key)

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)

		totpKey, err := service.Get(ctx, mockUow, userID, email)
		require.Error(t, err, "expected an error due to base32 decoding failure")
		require.Nil(t, totpKey, "expected nil TOTP key on base32 decoding failure")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get TOTP key generation error", func(t *testing.T) {
		t.Parallel()

		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)

		key := helper.ToAESKey(conf.Auth.EncryptionKey)
		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), key)

		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("GetTOTPSecret", userID).Return(&encryptedSecret, nil)

		invalidConf := conf
		invalidConf.App.Name = ""

		service := userservice.NewTOTPService(mockLogger, invalidConf, nil)
		totpKey, err := service.Get(ctx, mockUow, userID, email)
		require.Error(t, err, "expected an error due to TOTP key generation failure")
		require.Nil(t, totpKey, "expected nil TOTP key on generation failure")

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
	})
}
