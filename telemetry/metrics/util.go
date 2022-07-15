package qmetrics

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

func Meter(instrument string) metric.Meter {
	return global.MeterProvider().Meter(instrument)
}
