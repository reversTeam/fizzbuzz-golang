package common

import (
	"log"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)


type Exporter struct {
	host string
	port int
	interval int
	requests *prometheus.GaugeVec
}

func NewExporter(host string, port int, interval int) *Exporter {
	exp := &Exporter{
		host: host,
		port: port,
		interval: interval,
		requests : promauto.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "fizzbuzz_request_sec",
				Help: "Number of requests",
			}, []string{"code", "method", "path"}),
	}

	return exp
}

func (o *Exporter) Records() {
	go func() {
		for {
			
			time.Sleep(time.Duration(o.interval) * time.Second)
		}
	}()
}

func (o *Exporter) Serve() {
	uri := fmt.Sprintf("%s:%d", o.host, o.port)
	http.Handle("/metrics", promhttp.Handler())
	go func (uri string) {
		err := http.ListenAndServe(uri, nil)
		if err != nil {
			log.Fatal("Error listen metrics: ", err)
		}
	} (uri)
}

func (o *Exporter) IncrRequests(code int, method string, path string) {
	str_code := strconv.Itoa(code)
	o.requests.WithLabelValues(str_code, method, path).Inc()
}