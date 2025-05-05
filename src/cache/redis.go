package cache

import (
	"context"
	"fmt"
	"twitter/src/configs"
	"twitter/src/logger"

	"github.com/redis/go-redis/v9"
)

var log = logger.NewLogger()
var RedisClient *redis.Client

func InitRedis(cfg *configs.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		WriteTimeout: cfg.Redis.WriteTimeOut,
		ReadTimeout: cfg.Redis.ReadTimeOut,
		DialTimeout: cfg.Redis.DialTimeOut,
		Addr: fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB: cfg.Redis.DB,
		PoolSize: cfg.Redis.Poolsize,
		PoolTimeout: cfg.Redis.PoolTimeOut,
	})

	err := RedisClient.Ping(context.Background()).Err()
	if err != nil {
		return err
	}
	
	log.Info(logger.Redis, logger.Start, "started successfuly", nil)
	return nil
}

func GetRedis() *redis.Client {
	return RedisClient
}

func CloseRedis() {
	RedisClient.Close()
}