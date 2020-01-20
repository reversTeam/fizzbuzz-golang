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
	GRPC_DEFAULT_HOST = "127.0.0.1"
	GRPC_DEFAULT_PORT = 42001

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
	grpcHost = flag.String("grpc-host", GRPC_DEFAULT_HOST, "grpc server host")
	grpcPort = flag.Int("grpc-port", GRPC_DEFAULT_PORT, "grpc server port")

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

	// Create a gateway
	gw := common.NewGateway(ctx, *httpHost, *httpPort, *grpcHost, *grpcPort, []grpc.DialOption{grpc.WithInsecure()})
	// Add exporter for grafana export metrics
	gw.Http.InitExporter(*exporterHost, *exporterPort, *exporterInterval)
	// Catch ctrl+c and graceful stop 
	done := common.GracefulStopSignals(gw.Http)

	fizzbuzzService := fizzbuzz.NewService()
	gw.AddService(fizzbuzzService)

	if err := gw.Start(); err != nil {
		log.Fatal(err)
	}
	<-done
}