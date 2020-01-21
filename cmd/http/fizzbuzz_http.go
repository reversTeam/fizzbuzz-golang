package main
import (
	"log"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz"
	"flag"
)

const (
	// Default flag values for http server
	HTTP_DEFAULT_HOST = "127.0.0.1"
	HTTP_DEFAULT_PORT = 8080

	// Default flag values for grpc connexion
	SERVER_GRPC_DEFAULT_HOST = "127.0.0.1"
	SERVER_GRPC_DEFAULT_PORT = 42001

	// Default flag values for exporter expose
	EXPORTER_DEFAULT_HOST = "127.0.0.1"
	EXPORTER_DEFAULT_PORT = 4242
	EXPORTER_DEFAULT_INTERVAL = 1
)

var (
	// flags for http server
	httpHost = flag.String("http-host", HTTP_DEFAULT_HOST, "http gateway host")
	httpPort = flag.Int("http-port", HTTP_DEFAULT_PORT, "http gateway port")

	// flags for gprc connexion
	serverGrpcHost = flag.String("server-grpc-host", SERVER_GRPC_DEFAULT_HOST, "grpc server host")
	serverGrpcPort = flag.Int("server-grpc-port", SERVER_GRPC_DEFAULT_PORT, "grpc server port")

	// flags for exporter expose and interval conf
	exporterHost = flag.String("exporter-host", EXPORTER_DEFAULT_HOST, "exporter host")
	exporterPort = flag.Int("exporter-port", EXPORTER_DEFAULT_PORT, "exporter port")
	// Todo: Add watching values
	exporterInterval = flag.Int("exporter-interval", EXPORTER_DEFAULT_INTERVAL, "exporter interval")
)

func main() {
	// Instantiate context in background
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Parse flags
	flag.Parse()

	// Create a gateway configuration
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1099511627776), grpc.MaxCallSendMsgSize(1099511627776)),
	}
	// Create a gateway
	gw := common.NewGateway(ctx, *httpHost, *httpPort, *serverGrpcHost, *serverGrpcPort, opts)
	// Add exporter for grafana export metrics
	gw.Http.InitExporter(*exporterHost, *exporterPort, *exporterInterval)
	// Catch ctrl+c and graceful stop 
	done := common.GracefulStopSignals(gw.Http)

	// Create a fizzbuzz service
	fizzbuzzService := fizzbuzz.NewService()
	// Handle fizzbuzz service in gateway server
	gw.AddService(fizzbuzzService)

	// Start a gateway server handle http connexions
	if err := gw.Start(); err != nil {
		log.Fatal(err)
	}
	<-done
}