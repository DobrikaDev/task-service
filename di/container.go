package di

import (
	"DobrikaDev/task-service/internal/delivery"
	"DobrikaDev/task-service/internal/service/task"
	"DobrikaDev/task-service/internal/storage/sql"
	"DobrikaDev/task-service/internal/storage/sqlxtrm"
	"DobrikaDev/task-service/utils/config"
	"context"
	"net"
	"net/http"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Container struct {
	ctx                context.Context
	cfg                *config.Config
	logger             *zap.Logger
	taskService        *task.TaskService
	httpClient         *http.Client
	server             *delivery.Server
	transactionFactory *sqlxtrm.SqlxTransactionFactory
	transactionManager *sqlxtrm.SqlxTransactionManager
	db                 *sqlx.DB
	storage            *sql.SqlStorage
	netListener        *net.Listener
	grpcServer         *grpc.Server
}

func NewContainer(ctx context.Context, cfg *config.Config, logger *zap.Logger) *Container {
	return &Container{ctx: ctx, cfg: cfg, logger: logger}
}

func (c *Container) GetTaskService() *task.TaskService {
	return get(&c.taskService, func() *task.TaskService {
		return task.NewTaskService(c.GetStorage(), c.cfg, c.logger)
	})
}

func (c *Container) GetTransactionFactory() *sqlxtrm.SqlxTransactionFactory {
	return get(&c.transactionFactory, func() *sqlxtrm.SqlxTransactionFactory {
		return sqlxtrm.NewSqlxTransactionFactory(c.GetDB())
	})
}

func (c *Container) GetTransactionManager() *sqlxtrm.SqlxTransactionManager {
	return get(&c.transactionManager, func() *sqlxtrm.SqlxTransactionManager {
		trm, err := sqlxtrm.NewSqlxTransactionManager(c.GetDB())
		if err != nil {
			panic(err)
		}

		return trm
	})
}

func (c *Container) GetDB() *sqlx.DB {
	return get(&c.db, func() *sqlx.DB {
		return sql.MustCreateDB(c.cfg)
	})
}

func (c *Container) GetStorage() *sql.SqlStorage {
	return get(&c.storage, func() *sql.SqlStorage {
		return sql.NewStorage(c.GetTransactionFactory(), c.GetTransactionManager(), c.logger)
	})
}

func (c *Container) GetHTTPClient() *http.Client {
	return get(&c.httpClient, func() *http.Client {
		return http.DefaultClient
	})
}

func (c *Container) GetNetListener() *net.Listener {
	return get(&c.netListener, func() *net.Listener {
		listener, err := net.Listen("tcp", ":"+c.cfg.Port)
		if err != nil {
			panic(err)
		}
		return &listener
	})
}

func (c *Container) GetGRPCServer() *grpc.Server {
	return get(&c.grpcServer, func() *grpc.Server {
		grpcServer := grpc.NewServer()

		reflection.Register(grpcServer)
		return grpcServer
	})
}
func (c *Container) GetRpcServer() *delivery.Server {
	return get(&c.server, func() *delivery.Server {
		return delivery.NewServer(c.ctx, c.GetTaskService(), c.cfg, c.logger)
	})
}

func get[T comparable](obj *T, builder func() T) T {
	if *obj != *new(T) {
		return *obj
	}

	*obj = builder()
	return *obj
}
