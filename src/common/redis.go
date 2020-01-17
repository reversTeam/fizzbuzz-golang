package common

import (
	"flag"
	"github.com/go-redis/redis/v7"
	"fmt"
)

// Default listen http server configuration
const (
	REDIS_DEFAULT_HOST = "127.0.0.1"
	REDIS_DEFAULT_PASSWORD = ""
	REDIS_DEFAULT_PORT = 6379
)

type RedisClient struct {
	host *string
	port *int
	password *string
	Client *redis.Client
}


func NewRedisClient() *RedisClient {
	return &RedisClient{
		host: flag.String("redis-host", REDIS_DEFAULT_HOST, "Redis host"),
		port: flag.Int("redis-port", REDIS_DEFAULT_PORT, "Redis port"),
		password: flag.String("redis-password", REDIS_DEFAULT_PASSWORD, "Redis password"),
		Client: nil,
	}
}

func (o *RedisClient) Connect() (err error) {
	o.Client = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", *o.host, *o.port),
		Password: *o.password,
		DB: 0,
	})
	// _ muted pong return
	_, err = o.Client.Ping().Result()
	return err
}

func (o *RedisClient) CreateSortableIndex(index string, key string) (err error) {
	_, err = o.Client.Get(key).Result()
	if err == redis.Nil {
		_, err = o.Client.Set(key, 1, 0).Result()
		_, err = o.Client.ZAdd(index, &redis.Z{
			Score:  0,
			Member: key,
		}).Result()	
	} else if err != nil {
		return err
	}
	return nil	
}

func (o *RedisClient) GetHighterScore(index string) (params string, score uint64, err error) {
	res, err := o.Client.ZRangeWithScores(index, -1, -1).Result()

	if err != nil {
		return "", 0, err
	}

	return res[0].Member.(string), uint64(res[0].Score), nil
}

func (o *RedisClient) IncrIndex(index string, key string) (err error) {
	_, err = o.Client.ZIncrBy("counter", 1, key).Result()
	return err
}