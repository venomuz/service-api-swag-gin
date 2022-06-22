package v1

import (
	jwt "github.com/venomuz/service-api-swag-gin/ApiGateway/api/token"
	"github.com/venomuz/service-api-swag-gin/ApiGateway/config"
	"github.com/venomuz/service-api-swag-gin/ApiGateway/pkg/logger"
	"github.com/venomuz/service-api-swag-gin/ApiGateway/services"
	"github.com/venomuz/service-api-swag-gin/ApiGateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.RepositoryStorage
	jwtHandler     jwt.JwtHendler
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RepositoryStorage
	jwtHandler     jwt.JwtHendler
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:   c.Redis,
		jwtHandler:     c.jwtHandler,
	}
}
