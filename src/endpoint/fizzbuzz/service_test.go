package fizzbuzz_test

import (
	"testing"
	"flag"
	"log"
	"os"
	pbFizzbuzz "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
)

func TestGrpcRequest(t *testing.T) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("127.0.0.1:42001", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1099511627776)))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pbFizzbuzz.NewFizzBuzzClient(conn)

	fbStatsRequestError(t, c, "Fizzbuzz StatsRequest Check stats without data", 3, 5, 15, "fizz", "buzz", "rpc error: code = Unknown desc = No data found")

	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest Sample", 3, 5, 15, "fizz", "buzz", []string{"1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"})
	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest change buzz str", 3, 5, 15, "ok", "buzz", []string{"1","2","ok","4","buzz","ok","7","8","ok","buzz","11","ok","13","14","okbuzz"})
	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest change buzz str", 3, 5, 15, "fizz", "lol", []string{"1","2","fizz","4","lol","fizz","7","8","fizz","lol","11","fizz","13","14","fizzlol"})
	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest change int1", 2, 5, 15, "fizz", "buzz", []string{"1","fizz","3","fizz","buzz","fizz","7","fizz","9","fizzbuzz","11","fizz","13","fizz","buzz"})
	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest change int2", 3, 4, 15, "fizz", "buzz", []string{"1","2","fizz","buzz","5","fizz","7","buzz","fizz","10","11","fizzbuzz","13","14","fizz"})
	fbGetRequestSuccess(t, c, "Fizzbuzz GetRequest change limit", 3, 5, 25, "fizz", "buzz", []string{"1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz","16","17","fizz","19","buzz","fizz","22","23","fizz","buzz"})

	fbGetRequestError(t, c, "Fizzbuzz GetRequest 0 int1", 0, 5, 25, "fizz", "buzz", "rpc error: code = Unknown desc = int1 and int2 parameters need to be more than 0")
	fbGetRequestError(t, c, "Fizzbuzz GetRequest 0 int2", 3, 0, 25, "fizz", "buzz", "rpc error: code = Unknown desc = int1 and int2 parameters need to be more than 0")
	fbGetRequestError(t, c, "Fizzbuzz GetRequest int1 == int2", 3, 3, 25, "fizz", "buzz", "rpc error: code = Unknown desc = We can't choice between str1 or str2, because int1 and int2 has a same value")

	fbGetRequestSuccessLoop(t, c, "Fizzbuzz GetRequest Sample loop 100", 3, 5, 15, "fizz", "buzz", []string{"1","2","fizz","4","buzz","fizz","7","8","fizz","buzz","11","fizz","13","14","fizzbuzz"}, 10)

	fbStatsRequestSuccess(t, c, "Fizzbuzz StatsRequest Sample 101", 3, 5, 15, "fizz", "buzz", 11)
	fbGetRequestSuccessLoop(t, c, "Fizzbuzz GetRequest Sample loop 100", 3, 5, 10, "fizz", "buzz", []string{"1","2","fizz","4","buzz","fizz","7","8","fizz","buzz"}, 20)
	fbStatsRequestSuccess(t, c, "Fizzbuzz StatsRequest Sample 101", 3, 5, 10, "fizz", "buzz", 20)
}

func fbStatsRequestSuccess(t *testing.T, c pbFizzbuzz.FizzBuzzClient, name string, int1 uint64, int2 uint64, limit uint64, str1 string, str2 string, value uint64)  {
	t.Run(name, func (t2 *testing.T) {
		have, err := c.Stats(context.Background(), &empty.Empty{})
		if err != nil {
			log.Fatalf("Error when calling FizzBuzz Stats Request: %s", err)
		}

		want := &pbFizzbuzz.FizzBuzzStatsResponse{
			Int1: int1,
			Int2: int2,
			Limit: limit,
			Str1: str1,
			Str2: str2,
			Requests: value,
		}

		if !proto.Equal(have, want) {
			t.Errorf("[%s] Proto.Equal returned false:\nhave -> %v\nwant -> %v\n", name, have, want)
		}
	})
}

func fbStatsRequestError(t *testing.T, c pbFizzbuzz.FizzBuzzClient, name string, int1 uint64, int2 uint64, limit uint64, str1 string, str2 string, error string)  {
	t.Run(name, func (t2 *testing.T) {
		_, err := c.Stats(context.Background(), &empty.Empty{})
		if err == nil {
			log.Fatalf("[%s] The request success or want failed: %s", name, err)
		} else if err.Error() != error {
			t.Errorf("[%s] The error message is not name \nhave -> %s\nwant -> %s\n", name, err.Error(), error)
		}
	})
}

func fbGetRequestSuccess(t *testing.T, c pbFizzbuzz.FizzBuzzClient, name string, int1 uint64, int2 uint64, limit uint64, str1 string, str2 string, result []string)  {
	t.Run(name, func (t2 *testing.T) {
		have, err := c.Get(context.Background(), &pbFizzbuzz.FizzBuzzGetRequest{
			Int1: int1,
			Int2: int2,
			Limit: limit,
			Str1: str1,
			Str2: str2,
		})
		if err != nil {
			log.Fatalf("[%s] Error when calling FizzBuzz Get Request: %s", name, err)
		}

		want := &pbFizzbuzz.FizzBuzzGetResponse{
			Items: result,
		}

		if !proto.Equal(have, want) {
			t.Errorf("[%s] Proto.Equal returned false:\nhave -> %v\nwant -> %v\n", name, have, want)
		}
	})
}

func fbGetRequestSuccessLoop(t *testing.T, c pbFizzbuzz.FizzBuzzClient, name string, int1 uint64, int2 uint64, limit uint64, str1 string, str2 string, result []string, n int)  {
	t.Run(name, func (t2 *testing.T) {
		for i := 0; i < n; i++ {
			have, err := c.Get(context.Background(), &pbFizzbuzz.FizzBuzzGetRequest{
				Int1: int1,
				Int2: int2,
				Limit: limit,
				Str1: str1,
				Str2: str2,
			})
			if err != nil {
				log.Fatalf("Error when calling FizzBuzz Get Request: %s", err)
			}

			want := &pbFizzbuzz.FizzBuzzGetResponse{
				Items: result,
			}

			if !proto.Equal(have, want) {
				t.Errorf("[%s] Proto.Equal returned false:\nhave -> %v\nwant -> %v\n", name, have, want)
			}
		}
	})
}

func fbGetRequestError(t *testing.T, c pbFizzbuzz.FizzBuzzClient, name string, int1 uint64, int2 uint64, limit uint64, str1 string, str2 string, result string)  {
	t.Run(name, func (t2 *testing.T) {
		_, err := c.Get(context.Background(), &pbFizzbuzz.FizzBuzzGetRequest{
			Int1: int1,
			Int2: int2,
			Limit: limit,
			Str1: str1,
			Str2: str2,
		})
		if err == nil {
			t.Errorf("[%s] The request success or want failed", name)
		} else if err.Error() != result {
			t.Errorf("[%s] The error message is not name \nhave -> %s\nwant -> %s\n", name, err.Error(), result)
		}

	})
}


func TestMain(m *testing.M) {
	flag.Parse()
    os.Exit(m.Run())
}
