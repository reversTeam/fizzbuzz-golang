package common

import (
	"net/http"
	"time"
	"fmt"
	"log"
	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// Definition of HttpServer struct
type HttpServer struct {
	Host string
	Port int
	Server *http.Server
	State FizzbuzzServerState
	mux *http.ServeMux
	exporter *Exporter
}

// Init HttpServer
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

// Todo: move this code in Gateway
// Init and attach exporter to the HttpServer
func (o *HttpServer) InitExporter(exporterHost string, exporterPort int, exporterInterval int) {
	o.exporter = NewExporter(exporterHost, exporterPort, exporterInterval)
	// Todo: Add watched metrics
	// exporter.WatchedMetrics()

	// Serve start a goroutine
	o.exporter.Serve()
}

// If the exporter is setup, add http handler for catch metrics
func (o *HttpServer) Handle(path string, mux *runtime.ServeMux) {
	if o.exporter != nil {
		o.mux.Handle(path, o.exporter.HandleHttpHandler(mux))
	} else {
		o.mux.Handle(path, mux)
	}
}

// Start the http server, ready for handle connexion
func (o *HttpServer) Start() error {
	uri := fmt.Sprintf("%s:%d", o.Host, o.Port)
	log.Printf("[HTTP] Server listen on %s\n", uri)
	return o.Server.ListenAndServe()
}

// Catch the SIG_TERM and exit cleanly
func (o *HttpServer) GracefulStop() error {
	return o.Server.Shutdown(context.Background())
}