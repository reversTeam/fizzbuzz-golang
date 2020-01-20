# Fizzbuzz Golang Architecture

![Architecture](https://raw.github.com/reversTeam/fizzbuzz-golang/master/docs/assets/fizzbuzz-architecture.jpg)

In this architecture description we do not take into account information that is linked to deployment such as the notion of Loadbalancing service in Kubernetes. If you want more deployment information, please see [this section](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/deployment.md).


### Gateway & Grpc Server

The program currently works thanks to two main components:
 - [Gateway](https://github.com/reversTeam/fizzbuzz-golang/tree/master/main.go) : Which allows to expose GRPC endpoints on HTTP protocol
 - [GRPC Server](https://github.com/reversTeam/fizzbuzz-golang/tree/master/src/endpoin/main.go) : Who embed logic and grpc services


It is the gateway which connects to the GRPC service, as the following code shows us:
```golang
gw := common.NewGateway(ctx, *httpHost, *httpPort, *grpcHost, *grpcPort, opts)
fizzbuzzService := fizzbuzz.NewService()

gw.AddService(fizzbuzzService)

gw.Start()
```
In a configuration with several micro services we could imagine going through a configuration file or a command argument, in order to activate or deactivate certain modules in order to scale only the routes that need them.

The GRPC server is very similar, it will load the same service and instantiate it the elements necessary for its operation.
```golang
grpcServer := common.NewGrpcServer(ctx, *grpcHost, *grpcPort)
fizzbuzzService := fizzbuzz.NewService()

grpcServer.AddService(fizzbuzzService)

grpcServer.Start()
```

### Protobuf

The definition of the service is done using the [proto file](https://github.com/reversTeam/fizzbuzz-golang/tree/master/src/endpoin/fizzbuzz/protobuff/fizzbuzz.proto).

It will then be possible to determine the endpoints that will be exposed, with their paths and their HTTP verb thanks to the annotations:
```protobuf
import "google/api/annotations.proto";
```
This line allows you to include advanced logic such as the HTTP gateway, but note that it is possible on these same principles to make sure to connect an exposure layer in QraphQL. See this [project](https://github.com/google/rejoiner) for more details.
```
rpc Get(FizzBuzzGetRequest) returns (FizzBuzzGetResponse) {
	option (google.api.http) = { // use the annotations
		post: "/fizzbuzz"        // httpVerbe: "httpPath"
		body: "*"
	};
}
```
This simple declaration makes it possible to define the following routes:
 - http `POST` `/fizzbuzz`
 - grpc `Get` Please note this is not the verb http but the name of the GRPC method to implement

It also defines that the routes will take as parameter a `FizzBuzzGetRequest` structure, the definition of which is as follows:
```protobuf
message FizzBuzzGetRequest {
	uint64 Int1 = 1 [json_name="int1"];
	uint64 Int2 = 2 [json_name="int2"];
	uint64 Limit = 3 [json_name="limit"];
	string Str1 = 4 [json_name="str1"];
	string Str2 = 5 [json_name="str2"];
}
```

The routes wait for the parameters in the following format:
```protobuf
message FizzBuzzGetResponse {
	repeated string Items = 6 [json_name="items"]; // repeated type == []type
}
```

The `.proto` files are not directly used, the following files must be generated:
  - `pb.go` : for the grpc server
  - `pb.gw.go` : for the gateway

You can generate them with the following command, provided that you have followed the [setup project](https://github.com/reversTeam/fizzbuzz-golang/tree/master/docs/setup.md).
```bash
make protogen
```

## Service

The service is the component that is really carrying the logic to perform, unlike the monolytic code this is a micro service and it will be deployed with a very limited number of dependencies.

It is used as a [ServiceInterface](https://github.com/reversTeam/fizzbuzz-golang/tree/master/src/common/service.go) in the Gateway and the GRPC Server, it will therefore be necessary to implement the following functions in order to respect the interface.
 - RegisterGateway: which allows the service to register on the gateway
 - RegisterGrpc: which allows the service to declare itself to the Grpc server

```golang
type ServiceInterface interface {
	RegisterGateway(*Gateway) error
	RegisterGrpc(*GrpcServer)
}

// used by ./main.go
func (o *FizzBuzz) RegisterGateway(gw *common.Gateway) error {
	uri := fmt.Sprintf("%s:%d", gw.GrpcHost, gw.GrpcPort)
	return pb.RegisterFizzBuzzHandlerFromEndpoint(gw.Ctx, gw.Mux, uri, gw.GrpcOpts)
}

// used by ./src/endpoint/main.go
func (o *FizzBuzz) RegisterGrpc(gs *common.GrpcServer) {
	pb.RegisterFizzBuzzServer(gs.Server, o)
}
```

It only remains for us to add the function that takes care of the endpoint that we declared earlier:
```golang
func (o *FizzBuzz) Get(ctx context.Context, in *pb.FizzBuzzGetRequest) (*pb.FizzBuzzGetResponse, error) {
	results := []string{}
	limit := uint64(in.Limit)
	int1 := uint64(in.Int1)
	int2 := uint64(in.Int2)

	if int1 * int2 == 0 {
		return nil, errors.New("int1 and int2 parameters need to be more than 0")
	}
	if in.Str1 == "" || in.Str2 == "" {
		return nil, errors.New("str1 and str2 parameters cannot be empty")
	}

	fizzbuzz := in.Str1+in.Str2
	for i := uint64(1); i <= limit; i++ {
		if i%(int1*int2) == 0 {
			results = append(results, fizzbuzz)
		} else if i%int1 == 0 {
			results = append(results, in.Str1)
		} else if i%int2 == 0 {
			results = append(results, in.Str2)
		} else {
			results = append(results, strconv.FormatUint(i, 10))
		}
	}
	
	return &pb.FizzBuzzGetResponse{Items: results}, nil
}
```


## Exporter

We have an [exporter](https://github.com/reversTeam/fizzbuzz-golang/tree/master/src/common/exporter.go) which is automatically plugged into the gateway and which will provide us with endpoint statistics.

To function they act as middleware by encapsulating the request which is send to the GRPC server, then it increments the `GaugeVec` corresponding to the return code http from the GRPC service, as well as the method and the path to use.

```golang
func (o *Exporter) HandleHttpHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		rwh := NewResponseWriterHandler(w)
		h.ServeHTTP(rwh, r)
		o.IncrRequests(rwh.StatusCode, method, path)
	})
}
```

There is a part to rework so that the metrics can be declared directly from the service, currently it is hard in construction to export it:
```golang
func NewExporter(host string, port int, interval int) *Exporter {
	exp := &Exporter{
		host: host,
		port: port,
		interval: interval,
		requests : promauto.NewGaugeVec(             // Todo : load from service
			prometheus.GaugeOpts{
				Name: "fizzbuzz_request_sec",
				Help: "Number of requests",
			}, []string{"code", "method", "path"}),
	}

	return exp
}
```