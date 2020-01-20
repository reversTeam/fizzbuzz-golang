package main

import (
	"flag"
	"log"
	"golang.org/x/net/context"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz"
)

const (
	// Default flag values for redis connexion
	REDIS_DEFAULT_HOST = "127.0.0.1"
	REDIS_DEFAULT_PASSWORD = ""
	REDIS_DEFAULT_PORT = 6379

	// Default flag values for GRPC server
	GRPC_DEFAULT_HOST = "127.0.0.1"
	GRPC_DEFAULT_PORT = 42001
)

var (
	// flags for redis configuration
	redisHost = flag.String("redis-host", REDIS_DEFAULT_HOST, "Redis host")
	redisPassword = flag.String("redis-password", REDIS_DEFAULT_PASSWORD, "Redis password")
	redisPort = flag.Int("redis-port", REDIS_DEFAULT_PORT, "Redis port")

	// flags for Grpc server
	grpcHost = flag.String("grpc-host", GRPC_DEFAULT_HOST, "Grpc listening host")
	grpcPort = flag.Int("grpc-port", GRPC_DEFAULT_PORT, "Grpc listening port")
)

func main() {
	// Instantiate context in background
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	flag.Parse()

	// Start a GRPC Server
	grpcServer := common.NewGrpcServer(ctx, *grpcHost, *grpcPort)
	done := common.GracefulStopSignals(grpcServer)

	// Instantiate a redis connexion
	redis := common.NewRedisClient(*redisHost, *redisPort, *redisPassword)
	if redis.IsConnected() != nil {
		log.Fatal("Cannot enabled redis connexion")
	}

	// Create your service
	fizzbuzzService := fizzbuzz.NewService()
	fizzbuzzService.SetRedis(redis)

	// Add your service
	grpcServer.AddService(fizzbuzzService)

	// Start Grpc Server
	grpcServer.Start()
	<-done
}