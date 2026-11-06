package middlewares

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ritchieridanko/stockflow/backend/internal/constants"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
)

func Request(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestID := strings.TrimSpace(ctx.GetHeader("X-Request-ID"))
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Writer.Header().Set("X-Request-ID", requestID)
		ctx.Request = ctx.Request.WithContext(
			context.WithValue(
				ctx.Request.Context(),
				constants.CtxKeyRequestID,
				requestID,
			),
		)
		ctx.Next()
	}
}
