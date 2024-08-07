package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/adapter/http/presenter"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/translation"
	"net/http"
)

type HealthHandler struct {
	trans translation.Translator
}

func NewHealthHandler(trans translation.Translator) *HealthHandler {
	return &HealthHandler{
		trans: trans,
	}
}

// Check godoc
// @Summary Health check
// @Description Check if the service is up and running
// @Tags Health
// @Accept json
// @Produce json
// @Param language path string true "language 2 abbreviations" default(en)
// @Success 200 {object} presenter.Response{message=string} "Successful response"
// @Router /{language}/v1/health/check [get]
func (r HealthHandler) Check(ctx *gin.Context) {
	presenter.NewResponse(ctx, r.trans).Message("iAmWorking").Echo(http.StatusOK)
}
