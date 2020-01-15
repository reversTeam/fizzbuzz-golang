package common

import (
	"google.golang.org/grpc"
	"flag"
	"net"
	"log"
	"fmt"
)

// Default listen http server configuration
const (
	DEFAULT_HOST = "127.0.0.1"
	DEFAULT_PORT = 42001
)

// Define the GRPC server state in ENUM
type GrpcServerState int
const (
	Init GrpcServerState = 0
	Boot GrpcServerState = 1
	Listen GrpcServerState = 3
	Ready GrpcServerState = 4
	Error GrpcServerState = 5
	Gracefull GrpcServerState = 6
	Stop GrpcServerState = 7
)

// Define the GRPC Server Struct
type GrpcServer struct {
	host *string
	port *int
	Server *grpc.Server
	State GrpcServerState
	listener net.Listener
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{
		host: flag.String("grpc-host", DEFAULT_HOST, "Grpc listening host"),
		port: flag.Int("grpc-port", DEFAULT_PORT, "Grpc listening port"),
		Server: grpc.NewServer(),
		State: Init,
		listener: nil,
	};
}

func (o *GrpcServer) Listen() (err error) {
	uri := fmt.Sprintf("%s:%d", *o.host, *o.port)
	o.listener, err = net.Listen("tcp", uri)
	if err != nil {
		o.State = Error
		log.Fatalf("failed to listen: %v", err)
	}
	o.State = Listen
	log.Printf("[GRPC] services started, listen on %s\n", uri)
	return err
}

func (o *GrpcServer) Ready() {
	go o.Server.Serve(o.listener)
	o.State = Ready
}