package common

import (
	"os/signal"
	"syscall"
	"log"
	"os"
)

// Defition of ServerGracefulStopableInterface for http & grpc server graceful stop
type ServerGracefulStopableInterface interface{
	GracefulStop() error
}

// Catch SIG_TERM and exit propely
func GracefulStopSignals(server ServerGracefulStopableInterface) (done chan bool) {
	done = make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer log.Println("[MAIN]: System is ready for catch exit's signals, To exit press CTRL+C")

	go func() {
		sig := <-sigs
		log.Println("[SYSTEM]: Signal catch:", sig)
		err := server.GracefulStop()
		if err != nil {
			log.Println("Server can't GracefulStop", err)
		}
		done <- true
	}()
	return
}
