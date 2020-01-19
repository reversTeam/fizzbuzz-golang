package common

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"fmt"
)

// Default listen http server configuration
const (
	
)

// Define the GRPC Server Struct
type GrpcServer struct {
	host string
	port int
	Server *grpc.Server
	State FizzbuzzServerState
	listener net.Listener
}

func NewGrpcServer(host string, port int) *GrpcServer {
	return &GrpcServer{
		host: host,
		port: port,
		Server: grpc.NewServer(),
		State: Init,
		listener: nil,
	};
}

func (o *GrpcServer) Listen() (err error) {
	uri := fmt.Sprintf("%s:%d", o.host, o.port)
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

func (o *GrpcServer) isGracefulable() bool {
	switch (o.State) {
	case
		Ready,
		Listen:
		return true
	}
	return false
}

func (o *GrpcServer) Graceful() {
	if o.isGracefulable() { 
		o.Server.GracefulStop()
	}
}