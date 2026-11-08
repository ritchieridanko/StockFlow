package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/handlers"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/middlewares"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type Router struct {
	router *gin.Engine
}

func Init(
	appName,
	clientAddr string,
	l *logger.Logger,
	oh *handlers.OrderHandler,
	ph *handlers.ProductHandler,
) *Router {
	r := gin.New()
	r.ContextWithFallback = true

	r.Use(middlewares.CORS(clientAddr))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "OK",
			},
		)
	})

	api := r.Group(
		"/api",
		middlewares.Request(l),
		middlewares.Recovery(l),
		otelgin.Middleware(appName),
		middlewares.Logging(l),
	)

	product := api.Group("/products")
	{
		product.POST("", ph.CreateProduct)

		id := product.Group("/:id")
		{
			id.PUT("", ph.UpdateProduct)
			id.DELETE("", ph.DeleteProduct)

			id.POST("/stock", ph.UpdateStock)
			id.GET("/movements", ph.GetAllProductStockMovements)
		}
	}

	order := api.Group("/orders")
	{
		order.POST("", oh.Create)

		id := order.Group("/:id")
		{
			id.GET("", oh.GetByID)
		}
	}

	return &Router{router: r}
}

func (r *Router) Router() *gin.Engine {
	return r.router
}
