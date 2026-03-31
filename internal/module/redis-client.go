package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/jsonschema-go/jsonschema"
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
		Type:  "redis.client",
		Title: "Redis Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Type: "string",
				},
				"port": {
					Type:    "integer",
					Minimum: jsonschema.Ptr[float64](1),
					Maximum: jsonschema.Ptr[float64](65535),
				},
			},
			Required:             []string{"host", "port"},
			AdditionalProperties: &jsonschema.Schema{Not: &jsonschema.Schema{}},
		},
		New: func(config config.ModuleConfig) (common.Module, error) {
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

func (rc *RedisClient) Printf(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	rc.logger.Debug(msg)
}

func (rc *RedisClient) Start(ctx context.Context) error {
	redis.SetLogger(rc)
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

func (rc *RedisClient) Stop() {
	rc.cancel()
}

func (rc *RedisClient) Get(key string) (any, error) {
	if rc.client != nil {
		val, err := rc.client.Get(rc.ctx, key).Result()
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, errors.New("redis.client not setup")
}

func (rc *RedisClient) Set(key string, value any) error {
	if rc.client != nil {
		status := rc.client.Set(rc.ctx, key, value, 0)
		return status.Err()
	}
	return errors.New("redis.client not setup")
}
