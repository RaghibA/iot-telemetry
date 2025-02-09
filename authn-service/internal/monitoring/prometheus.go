package monitoring

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/**
 * Prometheus metric monitoring
 *
 * HttpRequestDuration: time taken to process a request
 * HttpRequestStatus: running total of status codes for all requests
 */
var (
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_dur_sec",
			Help:    "Duration of http requests measure in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "route"},
	)

	HttpRequestStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_status_ct",
			Help: "Total requests",
		},
		[]string{"method", "route", "status_code"},
	)
)

func InitPrometheus() {
	prometheus.MustRegister(HttpRequestDuration)
	log.Println("Prometheus Collector Registered")
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
