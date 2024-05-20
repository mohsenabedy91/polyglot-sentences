package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"strings"
)

func Authentication(cfg config.Jwt) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeaderToken := ctx.Request.Header.Get(config.AuthorizationHeaderKey)
		if authHeaderToken == "" {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.NewServiceError(serviceerror.Unauthorized),
			).Echo()
			return
		}

		token := authHeaderToken[len("Bearer"):]
		token = strings.TrimSpace(token)

		validatedToken, err := validationToken(cfg.AccessTokenSecret, token)
		if err != nil {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(err).Echo()
			return
		}

		if !validatedToken.Valid {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.NewServiceError(serviceerror.InvalidToken),
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
				serviceerror.NewServiceError(serviceerror.ServerError),
			).Echo()
			return
		}

		ctx.Set(config.AuthUserUUIDKey, claims[config.AuthUserUUIDKey])

		ctx.Next()
	}
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
		return nil, serviceerror.NewServiceError(serviceerror.InvalidToken)
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

	return nil, serviceerror.NewServiceError(serviceerror.TokenExpired)
}
