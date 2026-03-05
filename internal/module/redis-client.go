package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	config config.ModuleConfig
	ctx    context.Context
	router common.RouteIO
	Host   string
	Port   uint16
	client *redis.Client
	logger *slog.Logger
	cancel context.CancelFunc
}

func init() {
	RegisterModule(ModuleRegistration{
		Type: "redis.client",
		New: func(config config.ModuleConfig) (Module, error) {
			params := config.Params
			hostString, err := params.GetString("host")
			if err != nil {
				return nil, errors.New("redis.client host error: " + err.Error())
			}

			portInt, err := params.GetInt("port")

			if err != nil {
				return nil, errors.New("redis.client port error: " + err.Error())
			}

			return &RedisClient{config: config, Host: hostString, Port: uint16(portInt), logger: CreateLogger(config)}, nil
		},
	})
}

func (rc *RedisClient) Id() string {
	return rc.config.Id
}

func (rc *RedisClient) Type() string {
	return rc.config.Type
}

func (rc *RedisClient) Start(ctx context.Context) error {
	rc.logger.Debug("running")
	router, ok := ctx.Value(common.RouterContextKey).(common.RouteIO)

	if !ok {
		return errors.New("redis.client unable to get router from context")
	}

	rc.router = router
	moduleContext, cancel := context.WithCancel(ctx)
	rc.ctx = moduleContext
	rc.cancel = cancel

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Password: "",
		DB:       0,
	})

	rc.client = client

	defer client.Close()

	<-rc.ctx.Done()
	rc.logger.Debug("done")
	return nil
}

func (rc *RedisClient) Output(ctx context.Context, payload any) error {

	return errors.ErrUnsupported
}

func (rc *RedisClient) Stop() {
	rc.cancel()
}

func (rc *RedisClient) Get(key string) (any, error) {

	switch key {
	case "host":
		return rc.Host, nil
	case "port":
		return rc.Port, nil
	default:
		if rc.client != nil {
			val, err := rc.client.Get(rc.ctx, key).Result()
			if err != nil {
				return nil, err
			}
			return val, nil
		}
		return nil, errors.New("redis.client key not found")
	}
}
