package qmetrics

import (
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
)

func init() {
	model.NameValidationScheme = model.LegacyValidation
}

func servePrometheus(registry *prometheus.Registry, port int) {
	http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("failed to start prometheus server", err)
	}
}
