package common

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"log"
	"fmt"
)

// Define the GRPC Server Struct
type GrpcServer struct {
	Ctx context.Context
	host string
	port int
	Server *grpc.Server
	State FizzbuzzServerState
	listener net.Listener
	services []ServiceInterface
}

// Create a grpc server
func NewGrpcServer(ctx context.Context, host string, port int) *GrpcServer {
	return &GrpcServer{
		Ctx: ctx,
		host: host,
		port: port,
		Server: grpc.NewServer(grpc.WriteBufferSize(2147483647), grpc.ReadBufferSize(2147483647)),
		State: Init,
		listener: nil,
		services: make([]ServiceInterface, 0),
	};
}

// Start listen tcp socket for handle grpc service
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

// Add service to the local service array, need for register later
func (o *GrpcServer) AddService(service ServiceInterface) {
	o.services = append(o.services, service)
}

// Register services to the grpc server
func (o *GrpcServer) startServices() {
	for _, service := range o.services {
		service.RegisterGrpc(o)
	}
}

// Start a grpc server ready for handle connexion
func (o *GrpcServer) Start() error {
	err := o.Listen()
	if err != nil {
		return err
	}
	o.startServices()
	go func(grpcServer *grpc.Server) {
		err := grpcServer.Serve(o.listener)
		// we can't catch this error, also we log it
		if err != nil {
			log.Fatal(err)
		}
	}(o.Server)
	o.State = Ready
	return nil
}

// Graceful stop, when SIG_TERM is send
func (o *GrpcServer) GracefulStop() error {
	if o.isGracefulStopable() { 
		o.Server.GracefulStop()
	}
	return nil
}

// Centralize GracefulStopable state
func (o *GrpcServer) isGracefulStopable() bool {
	switch (o.State) {
	case
		Ready,
		Listen:
		return true
	}
	return false
}