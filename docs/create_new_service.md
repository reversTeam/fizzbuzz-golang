# How to create my Service

To create a service you only need a few files:

```bash
mkdir -p src/endpoint/example/protobuf
touch src/endpoint/example/service.go src/endpoint/example/protobuf/hello.proto
```

Now you have all the files and nomenclature needed to run your service.

## Definition protobuf

Now let's go to the endpoint declaration in your protobuf file

```golang
syntax = "proto3";

package go.micro.service.example;

import "google/api/annotations.proto";
import "google/protobuf/empty.proto";

service Hello {
	rpc Hello(google.protobuf.Empty) returns (HelloResponse) {
		option (google.api.http) = {
			get: "/hello"
		};
	}
}

message HelloResponse {
	string Message = 6 [json_name="message"];
}
```

With this description of protobuf, you will have an endpoint `http` available on` GET` `/ hello`, or in grpc via the` Hello` method.

You can now launch the command to generate the files associated with your prototype.
```bash
make protogen
```


## Service

We must now take care of the service in order to be able to expose it.

```golang
package hello

import (
	"fmt"
	pb "github.com/reversTeam/fizzbuzz-golang/src/endpoint/hello/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context"
)

// Define the service structure
type Hello struct {

}

// Instanciate the service without dependency because it's role of ServiceFactory
// And because Gateway no need redis connexion for work
func NewService() *Hello {
	return &Hello{}
}


// Interface Service method for register protos on Gateway
func (o *Hello) RegisterGateway(gw *common.Gateway) error {
	uri := fmt.Sprintf("%s:%d", gw.GrpcHost, gw.GrpcPort)
	return pb.RegisterHelloHandlerFromEndpoint(gw.Ctx, gw.Mux, uri, gw.GrpcOpts)
}

// Interface Service method for register on GRPC server
func (o *Hello) RegisterGrpc(gs *common.GrpcServer) {
	pb.RegisterHelloServer(gs.Server, o)
}


// Endpoint :
//  - grpc : hello
//  - http : GET /hello
func (o *Hello) Hello(ctx context.Context, in *empty.Empty) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Message: "Hello world !",
	}, nil
}
```

## Expose service

Now we just have to expose the service through the different servers.

grpc: `./cmd/grpc/fizzbuzz_grpc.go`
```golang
package main

import (
	// [...]
	"github.com/reversTeam/fizzbuzz-golang/src/hello/fizzbuzz"
)

func main() {
	// [...]

	// Create your service
	helloService := hello.NewService()
	// Add your service
	grpcServer.AddService(helloService)

	// [...]
}
```

http: `./cmd/http/fizzbuzz_http.go`
```golang
package main

import (
	// [...]
	"github.com/reversTeam/fizzbuzz-golang/src/hello/fizzbuzz"
)

func main() {
	// [...]

	// Create your service
	helloService := hello.NewService()
	// Add your service
	gw.AddService(helloService)

	// [...]
}
```

You can now check if your service is exposed
```
go run cmd/grpc/fizzbuzz_grpc.go
go run cmd/http/fizzbuzz_http.go
curl -X GET http://127.0.0.1:8080/hello
```
