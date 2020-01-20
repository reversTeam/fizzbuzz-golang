package fizzbuzz

import (
	"strconv"
	"strings"
	"fmt"
	"log"
	"errors"
	"github.com/go-redis/redis/v7"
	pb "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
	"sync"
	"math"
)

// Define the service structure
type FizzBuzz struct {
	redis *common.RedisClient
}

// Instanciate the service without dependency because it's role of ServiceFactory
// And because Gateway no need redis connexion for work
func NewService() *FizzBuzz {
	return &FizzBuzz{
		redis: nil,
	}
}

// Register external dependency after construct
func (o *FizzBuzz) SetRedis(redis *common.RedisClient) (err error) {
	o.redis = redis
	_, err = o.redis.Client.Ping().Result()

	return err
}

// Interface Service method for register protos on Gateway 
func (o *FizzBuzz) RegisterGateway(gw *common.Gateway) error {
	uri := fmt.Sprintf("%s:%d", gw.GrpcHost, gw.GrpcPort)
	return pb.RegisterFizzBuzzHandlerFromEndpoint(gw.Ctx, gw.Mux, uri, gw.GrpcOpts)
}

// Interface Service method for register on GRPC server
func (o *FizzBuzz) RegisterGrpc(gs *common.GrpcServer) {
	pb.RegisterFizzBuzzServer(gs.Server, o)
}

// Endpoint :
//  - grpc : Get
//  - http : POST /fizzbuzz
func (o *FizzBuzz) Get(ctx context.Context, in *pb.FizzBuzzGetRequest) (*pb.FizzBuzzGetResponse, error) {
	var wg sync.WaitGroup
	limit := uint64(in.Limit)
	results := make([]string, limit)
	int1 := uint64(in.Int1)
	int2 := uint64(in.Int2)

	if int1 * int2 == 0 {
		return nil, errors.New("int1 and int2 parameters need to be more than 0")
	}
	if int1 == int2 {
		return nil, errors.New("We can't choice between str1 or str2, because int1 and int2 has a same value")	
	}
	if in.Str1 == "" || in.Str2 == "" {
		return nil, errors.New("str1 and str2 parameters cannot be empty")
	}

	fizzbuzz := in.Str1+in.Str2
	key := fmt.Sprintf("%d:%d:%d:%s:%s", int1, int2, limit, in.Str1, in.Str2);
	if o.redis.CreateSortableIndex("counter", key) != nil {
		// we can accept to continue but we lost the bonus
		return nil, errors.New("Internal error cannot init the counter")
	}
	nbWait := int(math.Ceil(float64(limit) / float64(10000)))
	wg.Add(nbWait)
	for i := uint64(0); i < limit; {
		go func(waitGroup *sync.WaitGroup, results []string, k uint64, l uint64, int1 uint64, int2 uint64, str1 string, str2 string, str1str2 string) {
			for  j := uint64(1); j <= 10000 && k < l; j++ {
					if (k+1)%(int1*int2) == 0 {
						results[k] = fizzbuzz
					} else if (k+1)%int1 == 0 {
						results[k] = str1
					} else if (k+1)%int2 == 0 {
						results[k] = str2
					} else {
						results[k] = strconv.FormatUint((k+1), 10)
					}
				k++
			}
			waitGroup.Done()
		}(&wg, results, i, limit, int1, int2, in.Str1, in.Str2, fizzbuzz)
		i += 10000
	}
	wg.Wait()
	if o.redis.IncrIndex("counter", key) != nil {
		// we can accept to continue but we lost the bonus
		return nil, errors.New("Internal error the index cannot be increase")
	}
	
	return &pb.FizzBuzzGetResponse{Items: results}, nil
}

// Endpoint :
//  - grpc : Stats
//  - http : GET /fizzbuzz
func (o *FizzBuzz) Stats(ctx context.Context, in *empty.Empty) (*pb.FizzBuzzStatsResponse, error) {
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
	int1, err := strconv.ParseUint(params[0], 10, 64)
	if err != nil {
		return nil, errors.New("Cannot read int1")
	}
	int2, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		return nil, errors.New("Cannot read int2")

	}
	limit, err := strconv.ParseUint(params[2], 10, 64)
	if err != nil {
		return nil, errors.New("Cannot read limit")
	}
	str1 := params[3]
	str2 := params[4]
	

	return &pb.FizzBuzzStatsResponse{
		Int1: uint64(int1),
		Int2: uint64(int2),
		Limit: uint64(limit),
		Str1: str1,
		Str2: str2,
		Requests: score,
	}, nil
}