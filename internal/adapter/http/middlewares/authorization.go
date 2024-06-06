package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/handler"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/domain"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/port"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/claim"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

func ACL(
	service port.ACLService,
	permissions ...domain.PermissionKeyType,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if len(ctx.Keys) == 0 {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.New(serviceerror.Unauthorized),
			).Echo()
			return
		}

		userUUID := claim.GetUserUUIDFromGinContext(ctx)

		isAllowed, err := service.CheckAccess(ctx.Request.Context(), userUUID, permissions...)
		if err != nil {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(err).Echo()
			return
		}

		if !isAllowed {
			presenter.NewResponse(ctx, nil, handler.StatusCodeMapping).Error(
				serviceerror.New(serviceerror.PermissionDenied),
			).Echo()
			return
		}

		ctx.Next()
	}
}
