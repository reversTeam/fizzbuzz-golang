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