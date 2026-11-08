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
	config     *configs.Config
	database   *database.Database
	transactor *database.Transactor
	logger     *logger.Logger

	or  repositories.OrderRepository
	oir repositories.OrderItemRepository
	pr  repositories.ProductRepository
	smr repositories.StockMovementRepository

	validator *validator.Validator

	ou usecases.OrderUsecase
	pu usecases.ProductUsecase

	oh *handlers.OrderHandler
	ph *handlers.ProductHandler

	router *router.Router
	server *server.Server
}

func Init(cfg *configs.Config, inf *infra.Infra) *Container {
	// Infra
	db := database.NewDatabase(inf.Database())
	tx := database.NewTransactor(inf.Database())
	l := logger.NewLogger(inf.Logger())

	// Repositories
	or := repositories.NewOrderRepository(db)
	oir := repositories.NewOrderItemRepository(db)
	pr := repositories.NewProductRepository(db)
	smr := repositories.NewStockMovementRepository(db)

	// Utils
	v := validator.Init()

	// Usecases
	ut := otel.Tracer("usecase")
	ou := usecases.NewOrderUsecase(or, oir, pr, smr, tx, v, ut)
	pu := usecases.NewProductUsecase(pr, smr, tx, v, ut)

	// Handlers
	oh := handlers.NewOrderHandler(ou)
	ph := handlers.NewProductHandler(pu)

	// Router
	r := router.Init(cfg.App.Name, cfg.Client.Addr, l, oh, ph)

	// Server
	srv := server.Init(&cfg.Server, r, l)

	return &Container{
		config:     cfg,
		database:   db,
		transactor: tx,
		logger:     l,
		or:         or,
		oir:        oir,
		pr:         pr,
		smr:        smr,
		validator:  v,
		ou:         ou,
		pu:         pu,
		oh:         oh,
		ph:         ph,
		router:     r,
		server:     srv,
	}
}

func (c *Container) Server() *server.Server {
	return c.server
}
