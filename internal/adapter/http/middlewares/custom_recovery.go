package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

func ErrorHandler(trans translation.Translator) func(ctx *gin.Context, err interface{}) {
	return func(ctx *gin.Context, err interface{}) {
		serviceErr := serviceerror.NewServerError()
		presenter.NewResponse(ctx, trans).Error(serviceErr).Echo()
	}
}
