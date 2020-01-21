package common

// Definition of service interface for register http & grpc server
type ServiceInterface interface {
	RegisterGateway(*Gateway) error
	RegisterGrpc(*GrpcServer)
}
