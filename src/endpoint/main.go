package main

import (
	"flag"
	"log"
	pb "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
)

func main() {
	grpcServer := common.NewGrpcServer()
	done := common.GrpcGracefullSignals(grpcServer.Server)
	fizzbuzzService := fizzbuzz.NewService()
	flag.Parse()
	err := fizzbuzzService.Init()
	if err != nil {
		log.Fatal(err)
	}

	grpcServer.Listen()
	pb.RegisterFizzBuzzServer(grpcServer.Server, fizzbuzzService)
	grpcServer.Ready()
	<-done
}