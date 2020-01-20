package common

type ServiceInterface interface {
	RegisterGateway(*Gateway) error
	RegisterGrpc(*GrpcServer)
}