package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
)

func ErrorHandler(ctx *gin.Context, err any) {
	serviceErr := serviceerror.NewServiceError(serviceerror.ServerError)
	presenter.NewResponse(ctx, nil).Error(serviceErr).Echo()
	return
}
