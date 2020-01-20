package common

import (
	"net/http"
	"time"
	"fmt"
	"log"
	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type HttpServer struct {
	Host string
	Port int
	Server *http.Server
	State FizzbuzzServerState
	mux *http.ServeMux
	exporter *Exporter
}

func NewHttpServer(host string, port int) *HttpServer {
	uri := fmt.Sprintf("%s:%d", host, port)
	mux := http.NewServeMux()
	return &HttpServer{
		Host: host,
		Port: port,
		Server: &http.Server{
			Addr:           uri,
			Handler:        mux,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
		State: Init,
		mux: mux,
		exporter: nil,
	}
}

func (o *HttpServer) InitExporter(exporterHost string, exporterPort int, exporterInterval int) {
	o.exporter = NewExporter(exporterHost, exporterPort, exporterInterval)
	// Todo: Add watched metrics
	// exporter.WatchedMetrics()

	// Serve start a goroutine
	o.exporter.Serve()
}

func (o *HttpServer) Handle(path string, mux *runtime.ServeMux) {
	if o.exporter != nil {
		o.mux.Handle(path, o.exporter.HandleHttpHandler(mux))
	} else {
		o.mux.Handle(path, mux)
	}
}

func (o *HttpServer) Listen() error {
	uri := fmt.Sprintf("%s:%d", o.Host, o.Port)
	log.Printf("[HTTP] Server listen on %s\n", uri)
	return o.Server.ListenAndServe()
}

func (o *HttpServer) GracefulStop() {
	o.Server.Shutdown(context.Background())
}