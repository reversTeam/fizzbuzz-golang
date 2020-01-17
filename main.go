package main
import (
	"log"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/reversTeam/fizzbuzz-golang/src/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb_fizzbuzz "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
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

	httpServer *http.Server
)

func NewServer(host *string, port *int, mux *http.ServeMux) *http.Server {
	uri := fmt.Sprintf("%s:%d", *host, *port)
	return &http.Server{
		Addr:           uri,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func configureSignals() (done chan bool) {
	done = make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer log.Println("[MAIN]: System is ready for catch exit's signals, To exit press CTRL+C")

	go func() {
		sig := <-sigs
		log.Println("[SYSTEM]: Signal catch:", sig)
		if httpServer != nil {
			httpServer.Shutdown(context.Background())
		}
		done <- true
	}()
	return
}

func prettier(h http.Handler, exporter *common.Exporter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		path := r.URL.Path
		rwh := common.NewResponseWriterHandler(w)
		exporter.IncrConcurrency(rwh.StatusCode, method, path)
		h.ServeHTTP(rwh, r)
		exporter.DecrConcurrency(rwh.StatusCode, method, path)
	})
}

func run(exporter *common.Exporter) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := http.NewServeMux()

	gwmux, err := newGateway(ctx)
	if err != nil {
		panic(err)
	}
	mux.Handle("/", prettier(gwmux, exporter))

	uri := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Printf("[HTTP] Server listen on %s\n", uri)
	httpServer = NewServer(httpHost, httpPort, mux)
	return httpServer.ListenAndServe()
}

func newGateway(ctx context.Context) (http.Handler, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	gwmux := runtime.NewServeMux()
	if err := pb_fizzbuzz.RegisterFizzBuzzHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%s:%d", *clientGrpcHost, *clientGrpcPort), opts); err != nil {
		return nil, err
	}

	return gwmux, nil
}

func main() {
	flag.Parse()
	done := configureSignals()

	exporter := common.NewExporter(*exporterHost, *exporterPort, *exporterInterval)
	exporter.Records()
	exporter.Serve()

	if err := run(exporter); err != nil {
		log.Fatal(err)
	}
	<-done
}