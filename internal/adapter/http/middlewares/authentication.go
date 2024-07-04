package middlewares

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	cache "github.com/mohsenabedy91/polyglot-sentences/internal/adapter/storage/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/constant"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
)

func Authentication(conf config.Jwt, cacheDriver cache.Interface[any]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeaderToken := ctx.Request.Header.Get(config.AuthorizationHeaderKey)
		if authHeaderToken == "" || len(authHeaderToken) < len("Bearer") {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.New(serviceerror.Unauthorized),
			).Echo()
			return
		}

		token := authHeaderToken[len("Bearer"):]
		token = strings.TrimSpace(token)

		validatedToken, err := validationToken(conf.AccessTokenSecret, token)
		if err != nil {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(err).Echo()
			return
		}

		if !validatedToken.Valid {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.New(serviceerror.InvalidToken),
			).Echo()
			return
		}

		claims, err := getClaims(validatedToken)
		if err != nil {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(err).Echo()
			return
		}

		if claims == nil {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.NewServerError(),
			).Echo()
			return
		}

		if jti, ok := claims[config.AuthTokenJTI].(string); ok {
			if err = checkLogout(ctx.Request.Context(), cacheDriver, jti); err != nil {
				presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(err).Echo()
				return
			}
		}

		ctx.Set(config.AuthTokenJTI, claims[config.AuthTokenJTI])
		ctx.Set(config.AuthTokenExpirationTime, claims[config.AuthTokenExpirationTime])
		ctx.Set(config.AuthTokenUserUUID, claims[config.AuthTokenUserUUID])
		ctx.Next()
	}
}

func checkLogout(ctx context.Context, cacheDriver cache.Interface[any], jti string) error {

	result, err := cacheDriver.Get(ctx, fmt.Sprintf("%s:%s", constant.RedisAuthTokenPrefix, jti))
	if err != nil || result == nil {
		return serviceerror.NewServerError()

	} else if *result == constant.LogoutRedisValue {
		return serviceerror.New(serviceerror.UserLogout)
	}

	return nil
}

func validationToken(accessTokenSecret string, token string) (*jwt.Token, error) {
	validToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("it's not a valid token")
		}
		return []byte(accessTokenSecret), nil
	})

	if err != nil {
		return nil, serviceerror.New(serviceerror.InvalidToken)
	}

	return validToken, nil
}

func getClaims(validToken *jwt.Token) (map[string]interface{}, error) {
	claimsMaps := map[string]interface{}{}

	if claims, ok := validToken.Claims.(jwt.MapClaims); ok {
		for key, value := range claims {
			claimsMaps[key] = value
		}

		return claimsMaps, nil
	}

	return nil, serviceerror.New(serviceerror.TokenExpired)
}
