package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/serviceerror"
	"time"
)

func New(log logger.Logger, conf config.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", conf.Redis.Host, conf.Redis.Port),
		Password:           conf.Redis.Password,
		DB:                 conf.Redis.DB,
		DialTimeout:        conf.Redis.DialTimeout * time.Second,
		ReadTimeout:        conf.Redis.ReadTimeout * time.Second,
		WriteTimeout:       conf.Redis.WriteTimeout * time.Second,
		PoolSize:           conf.Redis.PoolSize,
		PoolTimeout:        conf.Redis.PoolTimeout * time.Second,
		IdleTimeout:        conf.Redis.IdleTimeout * time.Millisecond,
		IdleCheckFrequency: conf.Redis.IdleCheckFrequency * time.Millisecond,
	})

	if val, err := client.Ping().Result(); err != nil {
		log.Error(
			logger.Cache,
			logger.RedisPing,
			fmt.Sprintf("Redis client doesn't connected val: %s, error: %v", val, err.Error()),
			nil,
		)
		return nil, serviceerror.New(serviceerror.ServiceUnavailable)
	}

	log.Info(logger.Cache, logger.Startup, "Redis client initialized", nil)

	return client, nil
}
