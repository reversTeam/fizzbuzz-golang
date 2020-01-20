package common

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
)

type Gateway struct {
	Ctx context.Context
	Mux *runtime.ServeMux
	Http *HttpServer
	GrpcHost string
	GrpcPort int
	GrpcOpts []grpc.DialOption
	services []ServiceInterface
}

func NewGateway(
	ctx context.Context,
	httpHost string,
	httpPort int,
	grpcHost string,
	grpcPort int,
	grpcOpts []grpc.DialOption,
) *Gateway {
	return &Gateway{
		Ctx: ctx,
		Mux: runtime.NewServeMux(),
		Http: NewHttpServer(httpHost, httpPort),
		GrpcHost: grpcHost,
		GrpcPort: grpcPort,
		GrpcOpts: grpcOpts,
		services: make([]ServiceInterface, 0),
	}
}

func (o *Gateway) AddService(service ServiceInterface) {
	o.services = append(o.services, service)
}

func (o *Gateway) startServices() error {
	for _, service := range o.services {
		err := service.RegisterGateway(o)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *Gateway) Start() error {
	err :=  o.startServices()
	if err != nil {
		return err
	}
	o.Http.Handle("/", o.Mux)

	return o.Http.Listen();
}
