package main
import (
	"log"
	"net/http"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	gw "github.com/reversTeam/fizzbuzz-golang/src/endpoint/fizzbuzz/protobuf"
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

	CLIENT_GRPC_DEFAULT_HOST = "127.0.0.1"
	CLIENT_GRPC_DEFAULT_PORT = 42001
)

var (
	httpHost, httpPort, clientGrpcHost, clientGrpcPort = getFlags()
	httpServer *http.Server
)


func getFlags() (httpHost *string, httpPort *int, clientGrpcHost *string, clientGrpcPort *int) {
	httpHost = flag.String("http-host", HTTP_DEFAULT_HOST, "Default listening host")
	httpPort = flag.Int("http-port", HTTP_DEFAULT_PORT, "Default listening port")

	clientGrpcHost = flag.String("client-grpc-host", CLIENT_GRPC_DEFAULT_HOST, "Default client host connexion")
	clientGrpcPort = flag.Int("client-grpc-port", CLIENT_GRPC_DEFAULT_PORT, "Default client port connexion")

	flag.Parse()
	return
}

func NewServer(host *string, port *int, mux *http.ServeMux) *http.Server {
	uri := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
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


func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mux := http.NewServeMux()

	gwmux, err := newGateway(ctx)
	if err != nil {
		panic(err)
	}
	mux.Handle("/", gwmux)
	
	uri := fmt.Sprintf("%s:%d", *httpHost, *httpPort)
	log.Printf("[HTTP] Server listen on %s\n", uri)
	httpServer = NewServer(httpHost, httpPort, mux)
	return httpServer.ListenAndServe()
}

func newGateway(ctx context.Context) (http.Handler, error) {
	opts := []grpc.DialOption{grpc.WithInsecure()}

	gwmux := runtime.NewServeMux()
	if err := gw.RegisterFizzBuzzHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%s:%d", *clientGrpcHost, *clientGrpcPort), opts); err != nil {
		return nil, err
	}

	return gwmux, nil
}

func main() {
	done := configureSignals()

	if err := run(); err != nil {
		log.Fatal(err)
	}
	<-done
}