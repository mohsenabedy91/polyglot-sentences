package authservice

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

type JWTService struct {
	log          logger.Logger
	conf         config.Jwt
	cache        port.AuthCache
	jtiGenerator func() string
	signJWT      func(token *jwt.Token, secret []byte) (string, error)
}

func New(
	log logger.Logger,
	conf config.Jwt,
	cache port.AuthCache,
	jtiGenerator func() string,
	signJWT func(token *jwt.Token, secret []byte) (string, error),
) *JWTService {
	if jtiGenerator == nil {
		jtiGenerator = func() string {
			return uuid.New().String()
		}
	}
	if signJWT == nil {
		signJWT = func(token *jwt.Token, secret []byte) (string, error) {
			return token.SignedString(secret)
		}
	}
	return &JWTService{
		log:          log,
		conf:         conf,
		cache:        cache,
		jtiGenerator: jtiGenerator,
		signJWT:      signJWT,
	}
}

func (r JWTService) GenerateToken(userUUIDStr string) (*string, error) {
	mapClaims := jwt.MapClaims{}
	now := time.Now()
	mapClaims[config.AuthTokenUserUUID] = userUUIDStr

	jti := r.jtiGenerator()
	mapClaims[config.AuthTokenJTI] = jti

	mapClaims[config.AuthTokenIssuedAt] = now.Unix()

	accessTokenExpirationHour := r.conf.AccessTokenExpireDay * (24 * time.Hour)
	mapClaims[config.AuthTokenExpirationTime] = int(now.Add(accessTokenExpirationHour).Unix())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	jwtString, err := r.signJWT(token, []byte(r.conf.AccessTokenSecret))
	if err != nil {
		r.log.Error(logger.JWT, logger.JWTGenerate, err.Error(), nil)
		return nil, serviceerror.NewServerError()
	}

	go func() {
		ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, jti)
		_ = r.cache.SetTokenState(ctxWithTimeout, key, "", accessTokenExpirationHour)
	}()

	return &jwtString, nil
}

func (r JWTService) LogoutToken(ctx context.Context, jti string, exp int64) error {
	expTime := time.Unix(exp, 0)

	key := fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, jti)

	if err := r.cache.SetTokenState(ctx, key, constant.LogoutRedisValue, time.Until(expTime)); err != nil {
		return err
	}

	return nil
}
