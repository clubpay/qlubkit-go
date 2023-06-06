package qmetrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

func Meter(instrument string) metric.Meter {
	return otel.GetMeterProvider().Meter(instrument)
}
