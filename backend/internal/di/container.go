package di

import (
	"github.com/ritchieridanko/stockflow/backend/configs"
	"github.com/ritchieridanko/stockflow/backend/internal/infra"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/database"
	"github.com/ritchieridanko/stockflow/backend/internal/infra/logger"
	"github.com/ritchieridanko/stockflow/backend/internal/repositories"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/handlers"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/router"
	"github.com/ritchieridanko/stockflow/backend/internal/transport/server"
	"github.com/ritchieridanko/stockflow/backend/internal/usecases"
	"github.com/ritchieridanko/stockflow/backend/internal/utils/validator"
	"go.opentelemetry.io/otel"
)

type Container struct {
	config   *configs.Config
	database *database.Database
	logger   *logger.Logger

	pr repositories.ProductRepository

	validator *validator.Validator

	pu usecases.ProductUsecase

	ph *handlers.ProductHandler

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Repositories
	pr := repositories.NewProductRepository(db)

	// Utils
	v := validator.Init()

	// Usecases
	ut := otel.Tracer("usecase")
	pu := usecases.NewProductUsecase(pr, v, ut)

	// Handlers
	ph := handlers.NewProductHandler(pu)

	// Router
	r := router.Init(cfg.App.Name, cfg.Client.Addr, l, ph)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:    cfg,
		database:  db,
		logger:    l,
		pr:        pr,
		validator: v,
		pu:        pu,
		ph:        ph,
		router:    r,
		server:    srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
