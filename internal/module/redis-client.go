package module

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/jwetzell/showbridge-go/internal/common"
	"github.com/jwetzell/showbridge-go/internal/config"
	"github.com/redis/go-redis/v9"
)

func init() {
	RegisterModule(ModuleRegistration{
		Type:  "redis.client",
		Title: "Redis Client",
		ParamsSchema: &jsonschema.Schema{
			Type: "object",
			Properties: map[string]*jsonschema.Schema{
				"host": {
					Title:       "Host",
					Description: "the hostname or IP address of the Redis server to connect to",
					Type:        "string",
				},
				"port": {
					Title:       "Port",
					Description: "the port to use when connecting to the Redis server",
					Type:        "integer",
					Minimum:     jsonschema.Ptr[float64](1),
					Maximum:     jsonschema.Ptr[float64](65535),
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

type RedisClient struct {
	config       config.ModuleConfig
	ctx          context.Context
	inputHandler common.InputHandler
	Host         string
	Port         uint16
	client       *redis.Client
	logger       *slog.Logger
	cancel       context.CancelFunc
	clientMu     sync.Mutex
}

func (rc *RedisClient) Id() string {
	return rc.config.Id
}

func (rc *RedisClient) Type() string {
	return rc.config.Type
}

func (rc *RedisClient) Printf(ctx context.Context, format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	rc.logger.Debug(msg)
}

func (rc *RedisClient) Start(ctx context.Context, inputHandler common.InputHandler) error {
	redis.SetLogger(rc)
	rc.logger.Debug("running")
	rc.inputHandler = inputHandler
	moduleContext, cancel := context.WithCancel(ctx)
	rc.ctx = moduleContext
	rc.cancel = cancel

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rc.Host, rc.Port),
		Password: "",
		DB:       0,
	})

	rc.clientMu.Lock()
	rc.client = client
	rc.clientMu.Unlock()

	<-rc.ctx.Done()
	rc.logger.Debug("done")
	return nil
}

func (rc *RedisClient) Stop() {
	if rc.cancel != nil {
		defer rc.cancel()
	}
	rc.clientMu.Lock()
	defer rc.clientMu.Unlock()
	if rc.client != nil {
		rc.client.Close()
	}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (any, error) {
	if rc.client != nil {
		val, err := rc.client.Get(ctx, key).Result()
		if err != nil {
			return nil, err
		}
		return val, nil
	}
	return nil, errors.New("redis.client not setup")
}

func (rc *RedisClient) Set(ctx context.Context, key string, value any) error {
	if rc.client != nil {
		status := rc.client.Set(ctx, key, value, 0)
		return status.Err()
	}
	return errors.New("redis.client not setup")
}
