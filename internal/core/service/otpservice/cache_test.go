package otpservice_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis/authrepository"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/service/otpservice"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func generateOTP(digits int) string {
	return strings.Repeat("0", digits)
}

func TestOTPCacheService_Set(t *testing.T) {
	conf := config.OTP{
		ExpireSecond:               60,
		ForgetPasswordExpireSecond: 120,
		Digits:                     6,
	}
	ctx := context.TODO()
	email := faker.Email()
	otpValue := generateOTP(int(conf.Digits))

	key := fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(email))

	t.Run("Set success first time", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Value == otpValue
		}), conf.ExpireSecond).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Set(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Set success update existing OTP", func(t *testing.T) {
		existingOTPState := &domain.OTP{
			Value:        "123456",
			Used:         false,
			RequestCount: 1,
			CreatedAt:    time.Now().Unix() - 1_000,
			LastRequest:  time.Now().Unix() - 1_000,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(existingOTPState, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Value == otpValue && state.RequestCount == 2
		}), conf.ExpireSecond).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Set(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Set failure cache set error", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, nil)
		mockOTPCache.On("Set", ctx, key, mock.Anything, conf.ExpireSecond).Return(serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Set(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}

func TestOTPCacheService_Validate(t *testing.T) {
	conf := config.OTP{
		ExpireSecond:               60,
		ForgetPasswordExpireSecond: 120,
		Digits:                     6,
	}
	ctx := context.TODO()
	email := faker.Email()
	otpValue := generateOTP(int(conf.Digits))

	key := fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(email))

	t.Run("Validate success", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: otpValue,
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Validate(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Validate failure invalid OTP", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Validate(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidOTP, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Validate failure OTP not found", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, serviceerror.New(serviceerror.InvalidOTP))

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Validate(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidOTP, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}

func TestOTPCacheService_Used(t *testing.T) {
	conf := config.OTP{
		ExpireSecond: 60,
	}
	ctx := context.TODO()
	email := faker.Email()

	key := fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(email))

	t.Run("Used success", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Used == true
		}), conf.ExpireSecond).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Used(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Used failure OTP already used", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  true,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Used(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Used failure OTP not found", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Used(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Used failure cache get error", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Used(ctx, email)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("Used failure cache set error", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)
		mockOTPCache.On("Set", ctx, key, mock.Anything, conf.ExpireSecond).Return(serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.Used(ctx, email)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}

func TestOTPCacheService_SetForgetPassword(t *testing.T) {
	conf := config.OTP{
		ExpireSecond:               60,
		ForgetPasswordExpireSecond: 120,
		Digits:                     6,
	}
	ctx := context.TODO()
	email := faker.Email()
	otpValue := generateOTP(int(conf.Digits))

	key := fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(email))

	t.Run("SetForgetPassword success first time", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Value == otpValue
		}), conf.ForgetPasswordExpireSecond).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.SetForgetPassword(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("SetForgetPassword success update existing OTP", func(t *testing.T) {
		existingOTPState := &domain.OTP{
			Value:        "123456",
			Used:         false,
			RequestCount: 1,
			CreatedAt:    time.Now().Unix() - 1_000,
			LastRequest:  time.Now().Unix() - 1_000,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(existingOTPState, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Value == otpValue && state.RequestCount == 2
		}), conf.ForgetPasswordExpireSecond).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.SetForgetPassword(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("SetForgetPassword failure cache set error", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, nil)
		mockOTPCache.On("Set", ctx, key, mock.Anything, conf.ForgetPasswordExpireSecond).Return(serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.SetForgetPassword(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}

func TestOTPCacheService_ValidateForgetPassword(t *testing.T) {
	conf := config.OTP{
		ExpireSecond:               60,
		ForgetPasswordExpireSecond: 120,
		Digits:                     6,
	}
	ctx := context.TODO()
	email := faker.Email()
	otpValue := generateOTP(int(conf.Digits))

	key := fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(email))

	t.Run("ValidateForgetPassword success", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: otpValue,
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.ValidateForgetPassword(ctx, email, otpValue)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("ValidateForgetPassword failure invalid OTP", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.ValidateForgetPassword(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidOTP, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("ValidateForgetPassword failure OTP not found", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, serviceerror.New(serviceerror.InvalidOTP))

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.ValidateForgetPassword(ctx, email, otpValue)

		require.Error(t, err)
		require.Equal(t, serviceerror.InvalidOTP, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}

func TestOTPCacheService_UsedForgetPassword(t *testing.T) {
	conf := config.OTP{
		ExpireSecond:               60,
		ForgetPasswordExpireSecond: 120,
		Digits:                     6,
	}
	ctx := context.TODO()
	email := faker.Email()

	key := fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(email))

	t.Run("UsedForgetPassword success", func(t *testing.T) {
		otpState := &domain.OTP{
			Value:       "123456",
			Used:        false,
			LastRequest: time.Now().Add(-1 * time.Hour).Unix(),
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)
		mockOTPCache.On("Set", ctx, key, mock.MatchedBy(func(state *domain.OTP) bool {
			return state.Used == true
		}), mock.Anything).Return(nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.UsedForgetPassword(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("UsedForgetPassword failure OTP already used", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "123456",
			Used:  true,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.UsedForgetPassword(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("UsedForgetPassword failure OTP not found", func(t *testing.T) {
		otpState := &domain.OTP{
			Value: "",
			Used:  false,
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.UsedForgetPassword(ctx, email)

		require.NoError(t, err)

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("UsedForgetPassword failure cache get error", func(t *testing.T) {
		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(&domain.OTP{}, serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.UsedForgetPassword(ctx, email)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})

	t.Run("UsedForgetPassword failure cache set error", func(t *testing.T) {
		otpState := &domain.OTP{
			Value:       "123456",
			Used:        false,
			LastRequest: time.Now().Add(-1 * time.Hour).Unix(),
		}

		mockOTPCache := new(authrepository.MockOTPCache)
		mockOTPCache.On("Get", ctx, key).Return(otpState, nil)
		mockOTPCache.On("Set", ctx, key, mock.Anything, mock.Anything).Return(serviceerror.NewServerError())

		service := otpservice.NewOTPCache(conf, mockOTPCache)
		err := service.UsedForgetPassword(ctx, email)

		require.Error(t, err)
		require.Equal(t, serviceerror.ServerError, err.(*serviceerror.ServiceError).GetErrorMessage())

		mockOTPCache.AssertExpectations(t)
	})
}
