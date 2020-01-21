package common

import (
	"github.com/go-redis/redis/v7"
	"fmt"
)

// Definition of RedisClient
type RedisClient struct {
	Client *redis.Client
}


// Init a RedisClient
func NewRedisClient(host string, port int, password string) *RedisClient {
	return &RedisClient{
		Client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d", host, port),
			Password: password,
			DB: 0,
		}),
	}
}

// Check if the connexion is etablish with success
func (o *RedisClient) IsConnected() (err error) {
	_, err = o.Client.Ping().Result()
	return err
}

// Create an index with sortable capacity, and lock the key for unicity
func (o *RedisClient) CreateSortableIndex(index string, key string) (err error) {
	_, err = o.Client.Get(key).Result()
	if err == redis.Nil {
		_, err = o.Client.Set(key, 1, 0).Result()
		if err != nil {
			return err
		}
		_, err = o.Client.ZAdd(index, &redis.Z{
			Score:  0,
			Member: key,
		}).Result()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil	
}

// Get the key has better score (than more called)
func (o *RedisClient) GetHighterScore(index string) (params string, score uint64, err error) {
	res, err := o.Client.ZRangeWithScores(index, -1, -1).Result()

	if err != nil {
		return "", 0, err
	}

	return res[0].Member.(string), uint64(res[0].Score), nil
}

// Atomic increment the score of the key
func (o *RedisClient) IncrIndex(index string, key string) (err error) {
	_, err = o.Client.ZIncrBy("counter", 1, key).Result()
	return err
}