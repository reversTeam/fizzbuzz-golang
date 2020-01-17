package main

import (
	"flag"
	"log"
	pbFizzbuzz "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz"
)

// Default listen http server configuration
const (
	REDIS_DEFAULT_HOST = "127.0.0.1"
	REDIS_DEFAULT_PASSWORD = ""
	REDIS_DEFAULT_PORT = 6379
)

func main() {
	grpcServer := common.NewGrpcServer()
	done := common.GrpcGracefullSignals(grpcServer.Server)

	redis := common.NewRedisClient(
		flag.String("redis-host", REDIS_DEFAULT_HOST, "Redis host"),
		flag.Int("redis-port", REDIS_DEFAULT_PORT, "Redis port"),
		flag.String("redis-password", REDIS_DEFAULT_PASSWORD, "Redis password"),
	)
	flag.Parse()
	fizzbuzzService := fizzbuzz.NewService(redis)

	if redis.Connect() != nil {
		log.Fatal("Cannot enabled redis connexion")
	}
	err := fizzbuzzService.Init()
	if err != nil {
		log.Fatal(err)
	}

	grpcServer.Listen()
	pbFizzbuzz.RegisterFizzBuzzServer(grpcServer.Server, fizzbuzzService)
	grpcServer.Ready()
	<-done
}