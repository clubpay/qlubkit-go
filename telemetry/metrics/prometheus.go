package qmetrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func servePrometheus(registry *prometheus.Registry, port int) {
	http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
