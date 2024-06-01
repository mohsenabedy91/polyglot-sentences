package otpservice

import (
	"context"
	"crypto/subtle"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
	"time"
)

type OtpService struct {
	log       logger.Logger
	otpConfig config.OTP
	otpCache  *redis.CacheDriver[TransformOTPState]
}

func New(log logger.Logger, otpConfig config.OTP, cache *redis.CacheDriver[any]) *OtpService {
	return &OtpService{
		log:       log,
		otpConfig: otpConfig,
		otpCache:  (*redis.CacheDriver[TransformOTPState])(cache),
	}
}

func (r OtpService) Set(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err == nil && !otpState.Used && otpState.Value != "" {

		otpState.RequestCount++
		otpState.LastRequest = time.Now().Unix()
	} else {

		otpState = TransformOTPState{
			Value:        otp,
			Used:         false,
			RequestCount: 1,
			CreatedAt:    time.Now().Unix(),
			LastRequest:  time.Now().Unix(),
		}
	}

	if err = r.otpCache.Set(ctx, key, otpState, r.otpConfig.ExpireSecond); err != nil {
		r.log.Error(logger.Redis, logger.RedisSet, err.Error(), map[logger.ExtraKey]interface{}{
			logger.CacheKey:    key,
			logger.CacheSetArg: otpState,
		})
		return serviceerror.NewServerError()
	}

	return nil
}

func (r OtpService) Validate(ctx context.Context, key string, otp string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState.Value == "" {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	if subtle.ConstantTimeCompare([]byte(otpState.Value), []byte(otp)) != 1 {
		return serviceerror.New(serviceerror.InvalidOTP)
	}

	return nil
}

func (r OtpService) Used(ctx context.Context, key string) error {
	key = fmt.Sprintf("%s:%s", constant.RedisOTPPrefix, strings.ToLower(key))

	otpState, err := r.otpCache.Get(ctx, key)
	if err != nil || otpState.Used || otpState.Value == "" {
		if err != nil {
			r.log.Error(logger.Redis, logger.RedisSet, err.Error(), map[logger.ExtraKey]interface{}{
				logger.CacheKey: key,
			})
			return serviceerror.NewServerError()
		}

		return nil
	}

	otpState.Used = true
	if err = r.otpCache.Set(ctx, key, otpState, r.otpConfig.ExpireSecond); err != nil {
		r.log.Error(logger.Redis, logger.RedisSet, err.Error(), map[logger.ExtraKey]interface{}{
			logger.CacheKey:    key,
			logger.CacheSetArg: otpState,
		})
		return serviceerror.NewServerError()
	}

	return nil
}

// ----------------------
//        DTO
// ----------------------

type TransformOTPState struct {
	Value        string
	Used         bool
	RequestCount int8
	CreatedAt    int64
	LastRequest  int64
}
