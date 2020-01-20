package common

import (
	"os/signal"
	"syscall"
	"log"
	"os"
)


type ServerGracefulStopableInterface interface{
	GracefulStop()
}

func GracefulStopSignals(server ServerGracefulStopableInterface) (done chan bool) {
	done = make(chan bool, 1)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer log.Println("[MAIN]: System is ready for catch exit's signals, To exit press CTRL+C")

	go func() {
		sig := <-sigs
		log.Println("[SYSTEM]: Signal catch:", sig)
		server.GracefulStop()
		done <- true
	}()
	return
}