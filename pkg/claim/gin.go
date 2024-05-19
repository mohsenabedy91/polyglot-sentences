package claim

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
)

func GetUserUUIDFromGinContext(ctx *gin.Context) uuid.UUID {
	userUUIDStr := ctx.GetString(config.AuthUserUUIDKey)
	return uuid.MustParse(userUUIDStr)
}
