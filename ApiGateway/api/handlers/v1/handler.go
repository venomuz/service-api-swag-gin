package v1

import (
	"github.com/venomuz/service_api_swag_gin/ApiGateway/config"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/pkg/logger"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/services"
	"github.com/venomuz/service_api_swag_gin/ApiGateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.RepositoryStorage
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RepositoryStorage
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:   c.Redis,
	}
}
