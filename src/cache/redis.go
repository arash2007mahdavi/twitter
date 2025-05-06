package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"twitter/src/configs"
	"twitter/src/logger"

	"github.com/redis/go-redis/v9"
)

var log = logger.NewLogger()
var RedisClient *redis.Client

func InitRedis(cfg *configs.Config) error {
	RedisClient = redis.NewClient(&redis.Options{
		WriteTimeout: cfg.Redis.WriteTimeOut,
		ReadTimeout:  cfg.Redis.ReadTimeOut,
		DialTimeout:  cfg.Redis.DialTimeOut,
		Addr:         fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.Poolsize,
		PoolTimeout:  cfg.Redis.PoolTimeOut,
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

type RedisValue struct {
	Value interface{} `json:"value"`
	Valid bool        `json:"valid"`
}

func Get(redisClient *redis.Client, key string) (*RedisValue, error) {
	res, err := redisClient.Get(context.Background(), key).Result()
	if err != nil {
		return nil, fmt.Errorf("doesnt exists")
	}
	res2 := RedisValue{}
	err = json.Unmarshal([]byte(res), &res2)
	if err != nil {
		return nil, fmt.Errorf("problem in unmarshaling in redis get")
	}
	new_res := res2
	new_res.Valid = false
	redisClient.Set(context.Background(), key, new_res, 2*time.Minute)
	return &res2, nil
}

func Set(redisClient *redis.Client, key string, value *RedisValue, expiretime time.Duration) error {
	_, err := Get(redisClient, key)
	if err == nil {
		fmt.Errorf("this key exists")
	}
	err = redisClient.Set(context.Background(), key, value, expiretime).Err()
	if err != nil {
		return fmt.Errorf("error in set new key")
	}
	return nil
}
