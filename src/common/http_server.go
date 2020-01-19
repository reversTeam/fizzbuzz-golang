package common

import (
	"net/http"
	"time"
	"fmt"
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
	return &HttpServer{
		Host: host,
		Port: port,
		Server: nil,
		State: Init,
		mux: http.NewServeMux(),
		exporter: nil,
	}
}

func (o *HttpServer) Init(exporterHost string, exporterPort int, exporterInterval int) {
	uri := fmt.Sprintf("%s:%d", o.Host, o.Port)
	o.Server = &http.Server{
		Addr:           uri,
		Handler:        o.mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	o.exporter = NewExporter(exporterHost, exporterPort, exporterInterval)
	// Todo: Add watched metrics
	// exporter.WatchedMetrics()
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
	return o.Server.ListenAndServe()
}

func (o *HttpServer) Graceful() {
	o.Server.Shutdown(context.Background())
}