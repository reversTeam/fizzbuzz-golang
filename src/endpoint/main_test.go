package main

import (
	"testing"
	"flag"
	"log"
	"os"
	pbFizzbuzz "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

func TestGrpcGetRequest(t *testing.T) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial("127.0.0.1:42001", grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1099511627776)))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()
	c := pbFizzbuzz.NewFizzBuzzClient(conn)
	// defer os.Exit(0)
	t.Run("Fizzbuzz GetRequest", func (t2 *testing.T) {
		response, err := c.Get(context.Background(), &pbFizzbuzz.FizzBuzzGetRequest{
			Int1: 3,
			Int2: 5,
			Limit: 10,
			Str1: "fizz",
			Str2: "buzz",
		})
		if err != nil {
			log.Fatalf("Error when calling FizzBuzz Get Request: %s", err)
		}
		log.Printf("Response from server: %s", response.Items)
		
	})
}

func TestMain(m *testing.M) {
	flag.Parse()
    os.Exit(m.Run())
}