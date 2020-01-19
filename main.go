package main
import (
	"log"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb_fizzbuzz "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"flag"
	"fmt"
)

const (
	// Default listen http server configuration
	HTTP_DEFAULT_HOST = "127.0.0.1"
	HTTP_DEFAULT_PORT = 8080

	GRPC_DEFAULT_HOST = "127.0.0.1"
	GRPC_DEFAULT_PORT = 42001

	EXPORTER_DEFAULT_HOST = "127.0.0.1"
	EXPORTER_DEFAULT_PORT = 4242
	EXPORTER_DEFAULT_INTERVAL = 1
)

var (
	httpHost = flag.String("http-host", HTTP_DEFAULT_HOST, "http gateway host")
	httpPort = flag.Int("http-port", HTTP_DEFAULT_PORT, "http gateway port")

	clientGrpcHost = flag.String("grpc-host", GRPC_DEFAULT_HOST, "grpc server host")
	clientGrpcPort = flag.Int("grpc-port", GRPC_DEFAULT_PORT, "grpc server port")

	exporterHost = flag.String("exporter-host", EXPORTER_DEFAULT_HOST, "exporter host")
	exporterPort = flag.Int("exporter-port", EXPORTER_DEFAULT_PORT, "exporter port")
	exporterInterval = flag.Int("exporter-interval", EXPORTER_DEFAULT_INTERVAL, "exporter interval")
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	flag.Parse()
	httpServer := common.NewHttpServer(*httpHost, *httpPort)
	httpServer.Init(*exporterHost, *exporterPort, *exporterInterval)
	done := common.GracefulSignals(httpServer)
	defer cancel()

	grpcGateway := common.NewGrpcGateway(ctx)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb_fizzbuzz.RegisterFizzBuzzHandlerFromEndpoint(ctx, grpcGateway.Mux, fmt.Sprintf("%s:%d", *clientGrpcHost, *clientGrpcPort), opts)
	if err != nil {
		log.Fatal("Error register proto: ", err)
	}

	httpServer.Handle("/", grpcGateway.Mux)

	uri := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Printf("[HTTP] Server listen on %s\n", uri)

	if err := httpServer.Listen(); err != nil {
		log.Fatal(err)
	}
	<-done
}