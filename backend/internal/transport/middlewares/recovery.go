package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
	"github.com/ritchieridanko/stockflow/backend/internal/utils"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/ce"
)

func Recovery(l *logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				l.Error(
					ctx.Request.Context(),
					"PANIC RECOVERED",
					logger.NewField("method", ctx.Request.Method),
					logger.NewField("path", ctx.Request.URL.Path),
					logger.NewField("panic", fmt.Sprintf("%v", r)),
					logger.NewField("stack_trace", string(debug.Stack())),
				)

				ctx.AbortWithStatusJSON(
					http.StatusInternalServerError,
					gin.H{
						"status": http.StatusInternalServerError,
						"error": gin.H{
							"code":    ce.CodeInternal,
							"message": ce.MsgInternalServer,
						},
						"metadata": gin.H{
							"request_id": utils.CtxRequestID(ctx.Request.Context()),
							"timestamp":  time.Now().UTC(),
						},
					},
				)
			}
		}()

		ctx.Next()
	}
}
