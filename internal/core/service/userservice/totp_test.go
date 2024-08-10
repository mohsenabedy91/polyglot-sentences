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
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Set", ctx, email, mock.AnythingOfType("string")).Return(nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		totpKey, err := service.Enroll(ctx, email)
		require.NoError(t, err)
		require.NotNil(t, totpKey)
		require.NotEmpty(t, totpKey.Secret)
		require.NotEmpty(t, totpKey.URL)

		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enroll TOTP generation error", func(t *testing.T) {
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockLogger.On("Error", logger.TOTP, logger.EnrollTOTP, mock.Anything, mock.Anything).Return()

		invalidConf := conf
		invalidConf.App.Name = ""
		invalidService := userservice.NewTOTPService(mockLogger, invalidConf, mockTOTPCache)

		totpKey, err := invalidService.Enroll(ctx, email)
		require.Error(t, err)
		require.Nil(t, totpKey)
	})

	t.Run("Enroll cache error", func(t *testing.T) {
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Set", ctx, email, mock.AnythingOfType("string")).Return(errors.New("cache error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		totpKey, err := service.Enroll(ctx, email)
		require.Error(t, err)
		require.Nil(t, totpKey)

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
	otpKey, generateErr := totp.Generate(totp.GenerateOpts{
		Issuer:      conf.App.Name,
		AccountName: email,
		Digits:      otp.DigitsSix,
	})
	require.NoError(t, generateErr)

	validCode, generateCodeErr := totp.GenerateCode(otpKey.Secret(), time.Now().UTC())
	require.NoError(t, generateCodeErr)

	t.Run("Enable success", func(t *testing.T) {
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("UpdateTOTPSecret", userID, &encryptedSecret).Return(nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		require.NoError(t, service.Enable(ctx, mockUow, userID, email, validCode))

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enable cache get error", func(t *testing.T) {
		mockUow := new(repository.MockUnitOfWork)
		mockTOTPCache := new(userrepository.MockTOTPCache)
		mockTOTPCache.On("Get", ctx, email).Return("", errors.New("cache error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Enable(ctx, mockUow, userID, email, validCode)
		require.Error(t, err)

		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enable verification error", func(t *testing.T) {
		mockUow := new(repository.MockUnitOfWork)
		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		invalidCode := "654321"
		mockLogger.On("Warn", logger.TOTP, logger.EnableTOTP, fmt.Sprintf("The code «%s» is not valid ", invalidCode), mock.Anything).Return()

		err := service.Enable(ctx, mockUow, userID, email, invalidCode)
		require.Error(t, err)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.InvalidTOTPCode, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockTOTPCache.AssertExpectations(t)
	})

	t.Run("Enable repository error", func(t *testing.T) {
		mockUow := new(repository.MockUnitOfWork)
		mockRepo := new(repository.MockUserRepository)
		mockTOTPCache := new(userrepository.MockTOTPCache)

		encryptedSecret, _ := helper.EncryptSecret([]byte(otpKey.Secret()), helper.ToAESKey(conf.Auth.EncryptionKey))

		mockTOTPCache.On("Get", ctx, email).Return(encryptedSecret, nil)
		mockUow.On("UserRepository").Return(mockRepo)
		mockRepo.On("UpdateTOTPSecret", userID, &encryptedSecret).Return(errors.New("repository error"))

		service := userservice.NewTOTPService(mockLogger, conf, mockTOTPCache)
		err := service.Enable(ctx, mockUow, userID, email, validCode)
		require.Error(t, err)
		require.Equal(t, "repository error", err.Error())

		mockUow.AssertExpectations(t)
		mockRepo.AssertExpectations(t)
		mockTOTPCache.AssertExpectations(t)
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
	otpKey, generateErr := totp.Generate(totp.GenerateOpts{
		Issuer:      conf.App.Name,
		AccountName: email,
		Digits:      otp.DigitsSix,
	})
	require.NoError(t, generateErr)

	t.Run("Verify success", func(t *testing.T) {
		validCode, err := totp.GenerateCode(otpKey.Secret(), time.Now().UTC())
		require.NoError(t, err)

		valid, err := service.Verify(validCode, otpKey.Secret())
		require.NoError(t, err)
		require.True(t, valid)
	})

	t.Run("Verify invalid code", func(t *testing.T) {
		invalidCode := "123456"

		mockLogger.On("Warn", logger.TOTP, logger.EnableTOTP, fmt.Sprintf("The code «%s» is not valid ", invalidCode), mock.Anything).Return()

		valid, err := service.Verify(invalidCode, otpKey.Secret())
		require.Error(t, err)
		require.False(t, valid)
		require.IsType(t, &serviceerror.ServiceError{}, err)
		require.Equal(t, serviceerror.InvalidTOTPCode, err.(*serviceerror.ServiceError).GetErrorMessage())
	})
}
