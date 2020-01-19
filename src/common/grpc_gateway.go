package common

import (
	// "net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	// "google.golang.org/grpc"
	"golang.org/x/net/context"
)

type GrpcGateway struct {
	ctx context.Context
	Mux *runtime.ServeMux
}

func NewGrpcGateway(ctx context.Context) *GrpcGateway {
	return &GrpcGateway{
		ctx: ctx,
		Mux: runtime.NewServeMux(),
	}
}