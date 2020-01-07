package main

import (
	"google.golang.org/grpc"
	pb "github.com/reversTeam/fizzbuzz-golang/src/client/protobuf"
	"github.com/reversTeam/fizzbuzz-golang/src/client/services"
	"log"
	"net"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"fmt"
)

const (
	// Default listen http server configuration
	DEFAULT_HOST = "127.0.0.1"
	DEFAULT_PORT = 42001
)

var (
	grpcServer = grpc.NewServer()
	host, port = getFlags()
)


func getFlags() (host *string, port *int) {
	host = flag.String("host", DEFAULT_HOST, "Default listening host")
	port = flag.Int("port", DEFAULT_PORT, "Default listening port")

	flag.Parse()
	return
}

func configureSignals() (done chan bool) {
	done = make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer log.Println("[MAIN]: System is ready for catch exit's signals, To exit press CTRL+C")

	go func() {
		sig := <-sigs
		log.Println("[SYSTEM]: Signal catch:", sig)
		grpcServer.GracefulStop()
		done <- true
	}()
	return
}

func main() {
	done := configureSignals()
	uri := fmt.Sprintf("%s:%d", *host, *port)
	lis, err := net.Listen("tcp", uri)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
  
	pb.RegisterClientServer(grpcServer, services.NewServer())
	log.Printf("[GRPC] services started, listen on %s\n", uri)
	grpcServer.Serve(lis)
	<-done
}