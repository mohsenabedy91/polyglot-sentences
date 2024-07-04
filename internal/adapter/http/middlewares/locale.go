package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
)

type LanguageUri struct {
	Language string `uri:"language" binding:"required"`
}

func LocaleMiddleware(trans translation.Translator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var langUri LanguageUri
		if err := ctx.ShouldBindUri(&langUri); err != nil {
			serviceErr := serviceerror.NewServerError()
			presenter.NewResponse(ctx, trans).Error(serviceErr).Echo()
		}
		_ = trans.GetLocalizer(langUri.Language)

		ctx.Next()
	}
}
