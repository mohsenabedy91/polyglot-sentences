package otpservice

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
	"time"
)

type OTPCacheService struct {
	otpConfig config.OTP
	otpCache  port.OTPCache
}

func NewOTPCache(otpConfig config.OTP, cache port.OTPCache) OTPCacheService {
	return OTPCacheService{
		otpConfig: otpConfig,
		otpCache:  cache,
	}
}

func (r OTPCacheService) Set(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err == nil && otpState != nil && !otpState.Used && otpState.Value != "" {

		otpState.Value = otp
		otpState.RequestCount++
		otpState.LastRequest = time.Now().Unix()
	} else {
		otpState = &domain.OTP{
			Value:        otp,
			Used:         false,
			RequestCount: 1,
			CreatedAt:    time.Now().Unix(),
			LastRequest:  time.Now().Unix(),
		}
	}

	if err = r.otpCache.Set(ctx, key, otpState, r.otpConfig.ExpireSecond); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}

func (r OTPCacheService) Validate(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState == nil || otpState.Value == "" {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	if subtle.ConstantTimeCompare([]byte(otpState.Value), []byte(otp)) != 1 {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	return nil
}

func (r OTPCacheService) Used(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState.Used || otpState.Value == "" {
		if err != nil {
			return serviceerror.NewServerError()
		}

		return nil
	}

	otpState.Used = true
	if err = r.otpCache.Set(ctx, key, otpState, r.otpConfig.ExpireSecond); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}

func (r OTPCacheService) SetForgetPassword(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err == nil && otpState != nil && !otpState.Used && otpState.Value != "" {

		otpState.Value = otp
		otpState.RequestCount++
		otpState.LastRequest = time.Now().Unix()
	} else {

		otpState = &domain.OTP{
			Value:        otp,
			Used:         false,
			RequestCount: 1,
			CreatedAt:    time.Now().Unix(),
			LastRequest:  time.Now().Unix(),
		}
	}

	if err = r.otpCache.Set(ctx, key, otpState, r.otpConfig.ForgetPasswordExpireSecond); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}

func (r OTPCacheService) ValidateForgetPassword(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState.Value == "" {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	if subtle.ConstantTimeCompare([]byte(otpState.Value), []byte(otp)) != 1 {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	return nil
}

func (r OTPCacheService) UsedForgetPassword(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisForgetPasswordPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState.Used || otpState.Value == "" {
		if err != nil {
			return serviceerror.NewServerError()
		}

		return nil
	}

	requestTime := time.Unix(otpState.LastRequest, 0)
	expiryTime := requestTime.Add(r.otpConfig.ForgetPasswordExpireSecond)
	otpState.Used = true
	if err = r.otpCache.Set(ctx, key, otpState, expiryTime.Sub(time.Now())); err != nil {
		return serviceerror.NewServerError()
	}

	return nil
}
