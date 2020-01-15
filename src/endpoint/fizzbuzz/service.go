package fizzbuzz

import (
	"context"
	"strconv"
	"strings"
	"fmt"
	"log"
	"errors"
	"github.com/go-redis/redis/v7"
	pb "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
)

type FizzBuzz struct{
	redis *common.RedisClient
}

func NewService() *FizzBuzz {
	return &FizzBuzz{
		redis: common.NewRedisClient(),
	}
}

func (o *FizzBuzz) Init() (err error) {
	err = o.redis.Connect()
	if err != nil {
		return err
	}
	return nil
}

func initRedisIndex(rc *redis.Client, key string) (err error) {
	_, err = rc.Get(key).Result()
	if err == redis.Nil {
		_, err = rc.Set(key, 1, 0).Result()
		_, err = rc.ZAdd("counter", &redis.Z{
			Score:  0,
			Member: key,
		}).Result()	
	} else if err != nil {
		log.Println("Cannot be created", err)
		return err

	}

	return nil
}

func (o *FizzBuzz) Get(ctx context.Context, in *pb.FizzBuzzGetRequest) (*pb.FizzBuzzGetResponse, error) {
	results := []string{}
	limit := int(in.Limit)
	fizz := int(in.Int1)
	buzz := int(in.Int2)
	fizzbuzz := in.Str1+in.Str2
	key := fmt.Sprintf("%d:%d:%d:%s:%s", fizz, buzz, limit, in.Str1, in.Str2);
	err := initRedisIndex(o.redis.Client, key)
	if err != nil {
		// Don't kill the request but loggin it
		log.Printf("Error %s cannot init the redis structure %s\n", key, err)
	}
	for i := 1; i <= limit; i++ {
		if i%(fizz*buzz) == 0 {
			results = append(results, fizzbuzz)
		} else if i%fizz == 0 {
			results = append(results, in.Str1)
		} else if i%buzz == 0 {
			results = append(results, in.Str2)
		} else {
			results = append(results, strconv.Itoa(i))
		}
	}
	_, err = o.redis.Client.ZIncrBy("counter", 1, key).Result()
	if err != nil {
		log.Println("Cannot be incremented", err)
	}
	
	return &pb.FizzBuzzGetResponse{Items: results}, nil
}

func (o *FizzBuzz) Stats(ctx context.Context, in *pb.FizzBuzzStatsRequest) (*pb.FizzBuzzStatsResponse, error) {
	items, err := o.redis.Client.Keys("counter").Result()
	if err == redis.Nil || len(items) == 0 {
		return nil, errors.New("No data found")
	} else if err != nil {
		log.Fatal("Error get Sort: ", err)
	}

	res, err := o.redis.Client.ZRangeWithScores("counter", -1, -1).Result()

	params := strings.Split(res[0].Member.(string), ":");
	int1, err := strconv.Atoi(params[0])
	if err != nil {
		log.Fatal("Atoi failed int1: ", err)
	}
	int2, err := strconv.Atoi(params[1])
	if err != nil {
		log.Fatal("Atoi failed int2: ", err)
	}
	limit, err := strconv.Atoi(params[2])
	if err != nil {
		log.Fatal("Atoi failed limit: ", err)
	}
	str1 := params[3]
	str2 := params[4]
	request := uint64(res[0].Score)
	

	return &pb.FizzBuzzStatsResponse{
		Int1: int32(int1),
		Int2: int32(int2),
		Limit: int32(limit),
		Str1: str1,
		Str2: str2,
		Requests: uint64(request),
	}, nil
}