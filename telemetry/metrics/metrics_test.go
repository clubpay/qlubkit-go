package qmetrics_test

import (
	"context"
	"testing"
	"time"

	qmetrics "github.com/clubpay/qlubkit-go/telemetry/metrics"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMetrics(t *testing.T) {
	Convey("Metrics", t, func(c C) {
		Convey("Prometheus", func(c C) {
			met, err := qmetrics.New(
				qmetrics.WithPrometheus(2374),
			)
			c.So(err, ShouldBeNil)
			c.So(met, ShouldNotBeNil)
			m := qmetrics.Meter("x")

			c.So(m, ShouldNotBeNil)
			g, err := m.Int64UpDownCounter("g1")
			c.So(err, ShouldBeNil)
			c.So(g, ShouldNotBeNil)
			for i := 0; i < 10; i++ {
				g.Add(context.TODO(), 1)
				time.Sleep(time.Second)
			}

			h, err := m.Int64Histogram("hist1")
			c.So(err, ShouldBeNil)
			h.Record(context.TODO(), 10)

			time.Sleep(time.Second * 10)
		})

	})
}
