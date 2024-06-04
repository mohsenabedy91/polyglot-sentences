package claim

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
)

func GetUserUUIDFromGinContext(ctx *gin.Context) uuid.UUID {
	userUUIDStr := ctx.GetString(config.AuthTokenUserUUID)
	return uuid.MustParse(userUUIDStr)
}

func GetJTIFromGinContext(ctx *gin.Context) string {
	return ctx.GetString(config.AuthTokenJTI)
}

func GetExpFromGinContext(ctx *gin.Context) int64 {
	return int64(ctx.GetFloat64(config.AuthTokenExpirationTime))
}
