package authservice

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/helper"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type JWTService struct {
	log      logger.Logger
	conf     config.Jwt
	jwtCache *redis.CacheDriver[string]
}

func New(log logger.Logger, conf config.Jwt, cache *redis.CacheDriver[any]) *JWTService {
	return &JWTService{
		log:      log,
		conf:     conf,
		jwtCache: (*redis.CacheDriver[string])(cache),
	}
}

func (r JWTService) GenerateToken(userUUIDStr string) (*string, error) {
	mapClaims := jwt.MapClaims{}
	now := time.Now()
	mapClaims[config.AuthTokenUserUUID] = userUUIDStr

	jti := uuid.New().String()
	mapClaims[config.AuthTokenJTI] = jti

	mapClaims[config.AuthTokenIssuedAt] = now.Unix()

	accessTokenExpirationHour := r.conf.AccessTokenExpireDay * (24 * time.Hour)
	mapClaims[config.AuthTokenExpirationTime] = int(now.Add(accessTokenExpirationHour).Unix())

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, jti)
		err := r.jwtCache.Set(ctxWithTimeout, key, nil, accessTokenExpirationHour)
		if err != nil {
			r.log.Error(logger.JWT, logger.RedisSet, err.Error(), map[logger.ExtraKey]interface{}{
				logger.CacheKey: jti,
			})
		}
	}()

	jwtString, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		mapClaims,
	).SignedString([]byte(r.conf.AccessTokenSecret))

	if err != nil {
		r.log.Error(logger.JWT, logger.JWTGenerate, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	return &jwtString, nil
}

func (r JWTService) LogoutToken(ctx context.Context, jti string, exp int64) error {
	expTime := time.Unix(exp, 0)

	key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, jti)

	if err := r.jwtCache.Set(ctx, key, helper.StringPtr(constant.LogoutRedisValue), expTime.Sub(time.Now())); err != nil {
		return err
	}

	return nil
}
