package authservice_test

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/authservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

var wg sync.WaitGroup

func TestJWTService_GenerateToken(t *testing.T) {
	mockLogger := new(logger.MockLogger)

	conf := config.Jwt{
		AccessTokenSecret:    "secret",
		AccessTokenExpireDay: 1,
	}

	userUUID := uuid.New().String()
	expectedJTI := uuid.New().String()
	jtiGenerator := func() string {
		return expectedJTI
	}

	t.Run("GenerateToken success", func(t *testing.T) {
		mockCache := new(authrepository.MockAuthCache)

		wg.Add(1)

		key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, expectedJTI)
		mockCache.On("SetTokenState", mock.Anything, key, "", 24*time.Hour).
			Run(func(args mock.Arguments) {
				defer wg.Done()
			}).
			Return(nil)

		service := authservice.New(mockLogger, conf, mockCache, jtiGenerator, nil)
		token, err := service.GenerateToken(userUUID)

		wg.Wait()

		require.NoError(t, err)
		require.NotNil(t, token)

		parsedToken, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
			return []byte(conf.AccessTokenSecret), nil
		})
		require.NoError(t, err)
		require.True(t, parsedToken.Valid)

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		require.True(t, ok)
		require.Equal(t, userUUID, claims[config.AuthTokenUserUUID])
		require.Equal(t, expectedJTI, claims[config.AuthTokenJTI])

		mockCache.AssertExpectations(t)
	})

	t.Run("GenerateToken JWT signing error", func(t *testing.T) {
		mockCache := new(authrepository.MockAuthCache)

		mockSignJWT := func(token *jwt.Token, secret []byte) (string, error) {
			return "", fmt.Errorf("signing error")
		}

		mockLogger.On("Error", logger.JWT, logger.JWTGenerate, mock.AnythingOfType("string"), mock.Anything).Return()

		service := authservice.New(mockLogger, conf, mockCache, nil, mockSignJWT)
		token, err := service.GenerateToken(userUUID)

		require.Error(t, err)
		require.Nil(t, token)

		mockLogger.AssertExpectations(t)
	})
}

func TestJWTService_LogoutToken(t *testing.T) {
	mockLogger := new(logger.MockLogger)
	conf := config.Jwt{
		AccessTokenSecret: "secret",
	}

	expectedJTI := uuid.New().String()
	jtiGenerator := func() string {
		return expectedJTI
	}

	key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, expectedJTI)

	exp := time.Now().Add(24 * time.Hour)
	ctx := context.TODO()

	t.Run("LogoutToken success", func(t *testing.T) {
		mockCache := new(authrepository.MockAuthCache)

		mockCache.On("SetTokenState", mock.Anything, key, constant.LogoutRedisValue, mock.Anything).Return(nil)

		service := authservice.New(mockLogger, conf, mockCache, jtiGenerator, nil)
		err := service.LogoutToken(ctx, expectedJTI, exp.Unix())

		require.NoError(t, err)

		mockCache.AssertExpectations(t)
	})

	t.Run("LogoutToken failure cache error", func(t *testing.T) {
		mockCache := new(authrepository.MockAuthCache)

		mockCache.On("SetTokenState", mock.Anything, key, constant.LogoutRedisValue, mock.Anything).Return(serviceerror.NewServerError())

		service := authservice.New(mockLogger, conf, mockCache, jtiGenerator, nil)
		err := service.LogoutToken(ctx, expectedJTI, exp.Unix())

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockCache.AssertExpectations(t)
	})
}
