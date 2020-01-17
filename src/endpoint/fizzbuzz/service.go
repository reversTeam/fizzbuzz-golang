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

func NewService(redis *common.RedisClient) *FizzBuzz {
	return &FizzBuzz{
		redis: redis,
	}
}

func (o *FizzBuzz) Init() (err error) {
	return nil
}

func (o *FizzBuzz) Get(ctx context.Context, in *pb.FizzBuzzGetRequest) (*pb.FizzBuzzGetResponse, error) {
	results := []string{}
	limit := int(in.Limit)
	int1 := int(in.Int1)
	int2 := int(in.Int2)

	if int1 * int2 == 0 {
		return nil, errors.New("The int1 or int2 parameters cannot be equals to 0")
	}

	fizzbuzz := in.Str1+in.Str2
	key := fmt.Sprintf("%d:%d:%d:%s:%s", int1, int2, limit, in.Str1, in.Str2);
	if o.redis.CreateSortableIndex("counter", key) != nil {
		// we can accept to continue but we lost the bonus
		return nil, errors.New("Internal error cannot init the counter")
	}
	for i := 1; i <= limit; i++ {
		if i%(int1*int2) == 0 {
			results = append(results, fizzbuzz)
		} else if i%int1 == 0 {
			results = append(results, in.Str1)
		} else if i%int1 == 0 {
			results = append(results, in.Str2)
		} else {
			results = append(results, strconv.Itoa(i))
		}
	}
	if o.redis.IncrIndex("counter", key) != nil {
		// we can accept to continue but we lost the bonus
		return nil, errors.New("Internal error the index cannot be increase")
	}
	
	return &pb.FizzBuzzGetResponse{Items: results}, nil
}

func (o *FizzBuzz) Stats(ctx context.Context, in *pb.FizzBuzzStatsRequest) (*pb.FizzBuzzStatsResponse, error) {
	items, err := o.redis.Client.Keys("counter").Result()
	if err == redis.Nil || len(items) == 0 {
		return nil, errors.New("No data found")
	} else if err != nil {
		log.Fatal("No keys in counter: ", err)
	}

	key, score, err := o.redis.GetHighterScore("counter")
	if err != nil {
		return nil, err
	}

	params := strings.Split(key, ":");
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
	

	return &pb.FizzBuzzStatsResponse{
		Int1: int32(int1),
		Int2: int32(int2),
		Limit: int32(limit),
		Str1: str1,
		Str2: str2,
		Requests: score,
	}, nil
}