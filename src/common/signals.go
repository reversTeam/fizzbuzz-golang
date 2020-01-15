package common

import (
	"os/signal"
	"syscall"
	"log"
	"os"
	"google.golang.org/grpc"
)

func GrpcGracefullSignals(grpcServer *grpc.Server) (done chan bool) {
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

// func configureSignals() (done chan bool) {
// 	done = make(chan bool, 1)
// 	sigs := make(chan os.Signal, 1)
// 	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
// 	defer log.Println("[MAIN]: System is ready for catch exit's signals, To exit press CTRL+C")

// 	go func() {
// 		sig := <-sigs
// 		log.Println("[SYSTEM]: Signal catch:", sig)
// 		if httpServer != nil {
// 			httpServer.Shutdown(context.Background())
// 		}
// 		done <- true
// 	}()
// 	return
// }