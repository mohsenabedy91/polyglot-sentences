package userservice

import (
	"context"
	"encoding/base32"
	"fmt"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"time"
)

type TOTPService struct {
	log       logger.Logger
	conf      config.Config
	totpCache port.TOTPCache
}

func NewTOTPService(log logger.Logger, conf config.Config, totpCache port.TOTPCache) *TOTPService {
	return &TOTPService{
		log:       log,
		conf:      conf,
		totpCache: totpCache,
	}
}

func (r TOTPService) Enroll(ctx context.Context, email string) (*domain.TOTPKey, error) {
	otpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      r.conf.App.Name,
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		r.log.Error(logger.TOTP, logger.EnrollTOTP, err.Error(), nil)
		return nil, err
	}

	key := helper.ToAESKey(r.conf.Auth.EncryptionKey)
	encryptedSecret, encryptErr := helper.EncryptSecret([]byte(otpKey.Secret()), key)
	if encryptErr != nil {
		r.log.Error(logger.TOTP, logger.EnrollTOTP, encryptErr.Error(), nil)
		return nil, encryptErr
	}

	if cacheErr := r.totpCache.Set(ctx, email, encryptedSecret); cacheErr != nil {
		return nil, cacheErr
	}

	return &domain.TOTPKey{
		Secret: otpKey.Secret(),
		URL:    otpKey.URL(),
	}, nil
}

func (r TOTPService) Enable(ctx context.Context, uow port.UserUnitOfWork, userID uint64, email string, code string) error {
	encryptedSecret, err := r.totpCache.Get(ctx, email)
	if err != nil {
		return err
	}

	key := helper.ToAESKey(r.conf.Auth.EncryptionKey)
	decryptedSecret, decryptErr := helper.DecryptSecret(encryptedSecret, key)
	if decryptErr != nil {
		r.log.Error(logger.TOTP, logger.EnableTOTP, decryptErr.Error(), nil)
		return decryptErr
	}

	if valid, verifyErr := r.Verify(code, string(decryptedSecret)); verifyErr != nil && !valid {
		return verifyErr
	}

	if updateErr := uow.UserRepository().UpdateTOTPSecret(userID, &encryptedSecret); updateErr != nil {
		return updateErr
	}

	return nil
}

func (r TOTPService) Disable(ctx context.Context, uow port.UserUnitOfWork, userID uint64, code string) error {
	encryptedSecret, err := uow.UserRepository().GetTOTPSecret(userID)
	if err != nil {
		return err
	}
	if encryptedSecret == nil {
		return serviceerror.New(serviceerror.TOTPNotEnrolled)
	}

	key := helper.ToAESKey(r.conf.Auth.EncryptionKey)
	decryptedSecret, decryptErr := helper.DecryptSecret(*encryptedSecret, key)
	if decryptErr != nil {
		r.log.Error(logger.TOTP, logger.DisableTOTP, decryptErr.Error(), nil)
		return decryptErr
	}

	if valid, verifyErr := r.Verify(code, string(decryptedSecret)); verifyErr != nil && !valid {
		return verifyErr
	}

	if updateErr := uow.UserRepository().UpdateTOTPSecret(userID, nil); updateErr != nil {
		return updateErr
	}

	return nil
}

func (r TOTPService) Get(ctx context.Context, uow port.UserUnitOfWork, userID uint64, email string) (*domain.TOTPKey, error) {
	encryptedSecret, err := uow.UserRepository().GetTOTPSecret(userID)
	if err != nil {
		return nil, err
	}
	if encryptedSecret == nil {
		return nil, serviceerror.New(serviceerror.TOTPNotEnrolled)
	}

	key := helper.ToAESKey(r.conf.Auth.EncryptionKey)
	decryptedSecret, decryptErr := helper.DecryptSecret(*encryptedSecret, key)
	if decryptErr != nil {
		r.log.Error(logger.TOTP, logger.GetTOTP, decryptErr.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	decodeString, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(string(decryptedSecret))
	if err != nil {
		r.log.Error(logger.TOTP, logger.GetTOTP, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}
	otpKey, err := totp.Generate(totp.GenerateOpts{
		Issuer:      r.conf.App.Name,
		AccountName: email,
		Digits:      otp.DigitsSix,
		Secret:      decodeString,
	})
	if err != nil {
		r.log.Error(logger.TOTP, logger.GetTOTP, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	return &domain.TOTPKey{
		Secret: otpKey.Secret(),
		URL:    otpKey.URL(),
	}, nil
}

func (r TOTPService) Verify(code string, secret string) (bool, error) {
	if valid, err := totp.ValidateCustom(
		code,
		secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    30,
			Skew:      0,
			Digits:    otp.DigitsSix,
			Algorithm: otp.AlgorithmSHA1,
		},
	); err != nil || !valid {
		if err != nil {
			r.log.Error(logger.TOTP, logger.VerifyTOTP, err.Error(), map[logger.ExtraKey]interface{}{
				"Secret": secret,
			})
			return false, serviceerror.NewServerError()
		}

		r.log.Warn(logger.TOTP, logger.VerifyTOTP, fmt.Sprintf("The code «%s» is not valid ", code), nil)
		return false, serviceerror.New(serviceerror.InvalidTOTPCode)
	}

	return true, nil
}
