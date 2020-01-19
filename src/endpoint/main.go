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

	GRPC_DEFAULT_HOST = "127.0.0.1"
	GRPC_DEFAULT_PORT = 42001
)

var (
	redisHost = flag.String("redis-host", REDIS_DEFAULT_HOST, "Redis host")
	redisPassword = flag.String("redis-password", REDIS_DEFAULT_PASSWORD, "Redis password")
	redisPort = flag.Int("redis-port", REDIS_DEFAULT_PORT, "Redis port")

	grpcHost = flag.String("grpc-host", GRPC_DEFAULT_HOST, "Grpc listening host")
	grpcPort = flag.Int("grpc-port", GRPC_DEFAULT_PORT, "Grpc listening port")
)

func main() {
	grpcServer := common.NewGrpcServer(*grpcHost, *grpcPort)
	done := common.GracefulSignals(grpcServer)

	redis := common.NewRedisClient(*redisHost, *redisPort, *redisPassword)
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