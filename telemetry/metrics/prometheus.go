package qmetrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func servePrometheus(c prometheus.Collector, port int) {
	registry := prometheus.NewRegistry()
	registry.MustRegister(
		c,
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
